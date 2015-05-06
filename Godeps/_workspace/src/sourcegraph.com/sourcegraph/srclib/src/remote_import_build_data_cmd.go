package src

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
	"sourcegraph.com/sourcegraph/go-flags"
)

func initRemoteImportBuildCmd(remoteGroup *flags.Command) {
	importBuildCmd, err := remoteGroup.AddCommand("import-build",
		"tell a remote to import a build for a repository at a specific commit",
		"The import-build command tells the remote to import build data for a repository at a specific commit. To import build data that was produced locally, first run `src build-data upload` (or run `src push`, which performs both steps).",
		&remoteImportBuildCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
	if lrepo, err := openLocalRepo(); err == nil {
		SetOptionDefaultValue(importBuildCmd.Group, "commit", lrepo.CommitID)
	}
}

type RemoteImportBuildCmd struct {
	CommitID string `short:"c" long:"commit" description:"commit ID of data to import" required:"yes"`
}

var remoteImportBuildCmd RemoteImportBuildCmd

func (c *RemoteImportBuildCmd) Execute(args []string) error {
	cl := NewAPIClientWithAuthIfPresent()

	if GlobalOpt.Verbose {
		log.Printf("Creating a new import-only build for repo %q commit %q", remoteCmd.RepoURI, c.CommitID)
	}

	repo, _, err := cl.Repos.Get(sourcegraph.RepoSpec{URI: remoteCmd.RepoURI}, nil)
	if err != nil {
		return err
	}

	repoSpec := sourcegraph.RepoSpec{URI: remoteCmd.RepoURI}
	repoRevSpec := sourcegraph.RepoRevSpec{RepoSpec: repoSpec, Rev: c.CommitID}

	// Resolve to the full commit ID, and ensure that the remote
	// server knows about the commit.
	commit, err := getCommitWithRefreshAndRetry(cl, repoRevSpec)
	if err != nil {
		return err
	}
	repoRevSpec.CommitID = string(commit.ID)

	build, _, err := cl.Builds.Create(repoRevSpec, &sourcegraph.BuildCreateOptions{
		BuildConfig: sourcegraph.BuildConfig{
			Import: true,
			Queue:  false,
		},
		Force: true,
	})
	if err != nil {
		return err
	}
	if GlobalOpt.Verbose {
		log.Printf("Created build #%d", build.BID)
	}

	now := time.Now()
	host := fmt.Sprintf("local (USER=%s)", os.Getenv("USER"))
	buildUpdate := sourcegraph.BuildUpdate{StartedAt: &now, Host: &host}
	if _, _, err := cl.Builds.Update(build.Spec(), buildUpdate); err != nil {
		return err
	}

	importTask := &sourcegraph.BuildTask{
		BID:   build.BID,
		Op:    sourcegraph.ImportTaskOp,
		Queue: true,
	}
	tasks, _, err := cl.Builds.CreateTasks(build.Spec(), []*sourcegraph.BuildTask{importTask})
	if err != nil {
		return err
	}
	importTask = tasks[0]
	if GlobalOpt.Verbose {
		log.Printf("Created import task #%d", importTask.TaskID)
	}

	// Stream logs.
	done := make(chan struct{})
	go func() {
		var logOpt sourcegraph.BuildGetLogOptions
		loopsSinceLastLog := 0
		for {
			select {
			case <-done:
				return
			case <-time.After(time.Duration(loopsSinceLastLog+1) * 500 * time.Millisecond):
				logs, _, err := cl.Builds.GetTaskLog(importTask.Spec(), &logOpt)
				if err != nil {
					log.Printf("Warning: failed to get build logs: %s.", err)
					return
				}
				if len(logs.Entries) == 0 {
					loopsSinceLastLog++
					continue
				}
				logOpt.MinID = logs.MaxID
				for _, e := range logs.Entries {
					fmt.Println(e)
				}
				loopsSinceLastLog = 0
			}
		}
	}()

	defer func() {
		done <- struct{}{}
	}()
	taskID := importTask.TaskID
	started := false
	log.Printf("# Import queued. Waiting for task #%d in build #%d to start...", importTask.TaskID, build.BID)
	for i, start := 0, time.Now(); ; i++ {
		if time.Since(start) > 45*time.Minute {
			return fmt.Errorf("import timed out after %s", time.Since(start))
		}

		tasks, _, err := cl.Builds.ListBuildTasks(build.Spec(), nil)
		if err != nil {
			return err
		}
		importTask = nil
		for _, task := range tasks {
			if task.TaskID == taskID {
				importTask = task
				break
			}
		}
		if importTask == nil {
			return fmt.Errorf("task #%d not found in task list for build #%d", taskID, build.BID)
		}

		if !started && importTask.StartedAt.Valid {
			log.Printf("# Import started.")
			started = true
		}

		if importTask.EndedAt.Valid {
			if importTask.Success {
				log.Printf("# Import succeeded!")
			} else if importTask.Failure {
				log.Printf("# Import failed!")
				return fmt.Errorf("import failed")
			}
			break
		}

		time.Sleep(time.Duration(i) * 200 * time.Millisecond)
	}

	log.Printf("# View the repository at:")
	log.Printf("# %s://%s/%s@%s", cl.BaseURL.Scheme, cl.BaseURL.Host, repo.URI, repoRevSpec.Rev)

	return nil
}
