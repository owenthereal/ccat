package sourcegraph

import "github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/rwvfs"

type MockBuildDataService struct {
	FileSystem_ func(repo RepoRevSpec) (rwvfs.FileSystem, error)
}

func (s MockBuildDataService) FileSystem(repo RepoRevSpec) (rwvfs.FileSystem, error) {
	return s.FileSystem_(repo)
}
