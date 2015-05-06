package vcsclient

import (
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-vcs/vcs"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/vcsstore/git"
)

type VCSStore interface {
	RepositoryOpener
	git.GitTransporter
}

type MockVCSStore struct {
	Repository_   func(repoPath string) (vcs.Repository, error)
	GitTransport_ func(repoPath string) (git.GitTransport, error)
}

var _ VCSStore = (*MockVCSStore)(nil)

func (s *MockVCSStore) Repository(repoPath string) (vcs.Repository, error) {
	return s.Repository_(repoPath)
}
func (s *MockVCSStore) GitTransport(repoPath string) (git.GitTransport, error) {
	return s.GitTransport_(repoPath)
}
