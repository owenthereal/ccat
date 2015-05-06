package src

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"strings"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
)

func init() {
	remoteGroup, err := CLI.AddCommand("remote",
		"remote operations",
		"The remote command displays information about the remote repository corresponding to the local repository. Its subcommands perform operations on the remote repository.",
		&remoteCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
	remoteGroup.SubcommandsOptional = true
	setDefaultRepoURIOpt(remoteGroup)

	initRemoteRepoCmds(remoteGroup)
	initRemoteImportBuildCmd(remoteGroup)
}

type RemoteCmd struct {
	RepoURI string `short:"r" long:"repo" description:"repository URI (defaults to VCS 'srclib' or 'origin' remote URL)"`
}

var remoteCmd RemoteCmd

func (c *RemoteCmd) Execute(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("no such subcommand: src remote %v; see src remote --help", strings.Join(args, " "))
	}

	cl := NewAPIClientWithAuthIfPresent()
	rrepo, err := c.getRemoteRepo(cl)
	if err != nil {
		return err
	}
	printRemoteRepo(rrepo)

	log.Println()
	log.Println("# Run 'src remote --help' to see other remote operations you can perform.")

	return nil
}

func (c *RemoteCmd) getRemoteRepo(cl *sourcegraph.Client) (*sourcegraph.Repo, error) {
	if c.RepoURI == "" {
		lrepo, err := openLocalRepo()
		var errMsg string
		if lrepo == nil {
			errMsg = "no git/hg repository found in or above the current dir"
		} else if err == errNoVCSCloneURL {
			errMsg = err.Error() + "\n\n"
		} else {
			errMsg = err.Error()
		}
		return nil, errors.New(errMsg + "; to specify which remote repository to act upon instead of attempting automatic detection, use --repo (e.g., '--repo github.com/owner/repo')")
	}

	rrepo, _, err := cl.Repos.Get(sourcegraph.RepoSpec{URI: c.RepoURI}, nil)
	if sourcegraph.IsHTTPErrorCode(err, http.StatusNotFound) {
		return nil, fmt.Errorf("No repository exists on the remote with the URI %q. To add it, use 'src remote add'. The underlying error was: %s", c.RepoURI, err)
	}
	return rrepo, err
}
