package src

import (
	"log"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
	"sourcegraph.com/sourcegraph/go-flags"
)

func initRemoteRepoCmds(remoteGroup *flags.Command) {
	c, err := remoteGroup.AddCommand("add",
		"add the local repository to the remote",
		"The add command creates a remote repository corresponding to the current local repository.",
		&remoteAddCmd,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Set defaults.
	openLocalRepo()
	if localRepo != nil {
		if localRepo.CloneURL != "" {
			//SetOptionDefaultValue(c.Group, "uri", localRepo.URI())
			SetOptionDefaultValue(c.Group, "clone-url", localRepo.CloneURL)
		}
		SetOptionDefaultValue(c.Group, "vcs", localRepo.VCSType)
	}
}

type RemoteAddCmd struct {
	VCSType  string `long:"vcs" description:"VCS type" required:"yes"`
	CloneURL string `long:"clone-url" description:"clone URL" required:"yes"`
}

var remoteAddCmd RemoteAddCmd

func (c *RemoteAddCmd) Execute(args []string) error {
	cl := NewAPIClientWithAuthIfPresent()

	if lrepo, _ := openLocalRepo(); lrepo != nil {
		if c.CloneURL != lrepo.CloneURL {
			log.Printf("# Warning: you are creating a remote repository with a clone URL (%q) that doesn't match that of the current dir's repository (%q).", c.CloneURL, lrepo.CloneURL)
		}
		if c.VCSType != lrepo.VCSType {
			log.Printf("# Warning: you are creating a remote repository with a VCS type (%q) that doesn't match that of the current dir's repository (%q).", c.VCSType, lrepo.VCSType)
		}
	}

	newRepo := sourcegraph.NewRepoSpec{
		Type:        c.VCSType,
		CloneURLStr: c.CloneURL,
	}
	rrepo, _, err := cl.Repos.Create(newRepo)
	if err != nil {
		return err
	}
	printRemoteRepo(rrepo)
	return nil
}
