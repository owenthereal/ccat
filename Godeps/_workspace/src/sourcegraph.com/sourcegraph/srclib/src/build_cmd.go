package src

import (
	"log"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
)

func init() {
	buildCmd, err := CLI.AddCommand("build",
		"trigger a remote build",
		"The build command triggers a remote build of the repository.",
		&buildCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
	setDefaultCommitIDOpt(buildCmd)
}

type BuildCmd struct {
	CommitID string `short:"c" long:"commit" description:"commit ID to build" required:"yes"`
	Priority int    `short:"p" long:"priority" description:"build priority" default:"2"`
}

var buildCmd BuildCmd

func (c *BuildCmd) Execute(args []string) error {
	cl := NewAPIClientWithAuthIfPresent()

	rrepo, err := getRemoteRepo(cl)
	if err != nil {
		return err
	}

	repoRev := sourcegraph.RepoRevSpec{RepoSpec: rrepo.RepoSpec(), Rev: c.CommitID, CommitID: c.CommitID}
	build, _, err := cl.Builds.Create(repoRev, &sourcegraph.BuildCreateOptions{
		BuildConfig: sourcegraph.BuildConfig{
			Import:   true,
			Queue:    true,
			Priority: c.Priority,
		},
		Force: true,
	})
	if err != nil {
		return err
	}
	log.Printf("# Created build #%d", build.BID)

	return nil
}
