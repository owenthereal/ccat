package src

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
)

// getRemoteRepo gets the remote repository that corresponds to the
// local repository (from openLocalRepo). It does not respect any
// flags that override the repo URI to use. Commands that need to
// allow the user to override the repo URI should be under the
// "remote" subcommand and use "RemoteCmd.getRemoteRepo".
func getRemoteRepo(cl *sourcegraph.Client) (*sourcegraph.Repo, error) {
	lrepo, err := openLocalRepo()
	if err != nil {
		return nil, localRepoErr
	}
	if lrepo.CloneURL == "" {
		return nil, errNoVCSCloneURL
	}
	uri := lrepo.URI()
	if uri == "" {
		return nil, fmt.Errorf("getRemoteRepo: the local repo's URI is malformed: %s", lrepo.CloneURL)
	}

	rrepo, _, err := cl.Repos.Get(sourcegraph.RepoSpec{URI: uri}, nil)
	return rrepo, err
}

// getCommitWithRefreshAndRetry tries to get a repository commit. If
// it doesn't exist, it triggers a refresh of the repo's VCS data and
// then retries (until maxGetCommitVCSRefreshWait has elapsed).
func getCommitWithRefreshAndRetry(cl *sourcegraph.Client, repoRevSpec sourcegraph.RepoRevSpec) (*sourcegraph.Commit, error) {
	timeout := time.After(maxGetCommitVCSRefreshWait)
	done := make(chan struct{})
	var commit *sourcegraph.Commit
	var err error
	go func() {
		refreshTriggered := false
		for {
			commit, _, err = cl.Repos.GetCommit(repoRevSpec, nil)

			// Keep retrying if it's a 404, but stop trying if we succeeded, or if it's some other
			// error.
			if !sourcegraph.IsHTTPErrorCode(err, http.StatusNotFound) {
				break
			}

			if !refreshTriggered {
				_, err = cl.Repos.RefreshVCSData(repoRevSpec.RepoSpec)
				if err != nil {
					err = fmt.Errorf("failed to trigger VCS refresh for repo %s: %s", repoRevSpec.URI, err)
					break
				}
				log.Printf("Repository %s revision %s wasn't found on remote. Triggered refresh of VCS data; waiting %s.", repoRevSpec.URI, repoRevSpec.Rev, maxGetCommitVCSRefreshWait)
				refreshTriggered = true
			}
			time.Sleep(time.Second)
		}
		done <- struct{}{}
	}()
	select {
	case <-done:
		return commit, err
	case <-timeout:
		return nil, fmt.Errorf("repo %s revision %s not found on remote, even after triggering a VCS refresh and waiting %s (if you are sure that commit has been pushed, try again later)", repoRevSpec.URI, repoRevSpec.Rev, maxGetCommitVCSRefreshWait)
	}
}

const maxGetCommitVCSRefreshWait = time.Second * 10
