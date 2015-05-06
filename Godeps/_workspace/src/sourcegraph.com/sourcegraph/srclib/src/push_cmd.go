package src

import (
	"log"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
)

func init() {
	_, err := CLI.AddCommand("push",
		"upload and import the current commit (to a remote)",
		"The push command uploads and imports the current repository commit's build data to a remote. It is a wrapper around `src build-data upload` and `src remote import-build`.",
		&pushCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type PushCmd struct {
}

var pushCmd PushCmd

func (c *PushCmd) Execute(args []string) error {
	cl := NewAPIClientWithAuthIfPresent()
	rrepo, err := getRemoteRepo(cl)
	if err != nil {
		return err
	}

	repoSpec := sourcegraph.RepoSpec{URI: rrepo.URI}
	repoRevSpec := sourcegraph.RepoRevSpec{RepoSpec: repoSpec, Rev: localRepo.CommitID}

	if _, err := getCommitWithRefreshAndRetry(cl, repoRevSpec); err != nil {
		return err
	}

	if err := buildDataUploadCmd.Execute(nil); err != nil {
		return err
	}
	if err := remoteImportBuildCmd.Execute(nil); err != nil {
		return err
	}
	return nil
}
