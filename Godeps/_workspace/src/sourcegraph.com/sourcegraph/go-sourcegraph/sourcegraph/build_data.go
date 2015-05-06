package sourcegraph

import (
	"io"
	"os"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/rwvfs"
)

// BuildDataService communicates with the build data-related endpoints in the
// Sourcegraph API.
type BuildDataService interface {
	// FileSystem returns a virtual filesystem interface to the build
	// data for a repo at a specific commit.
	FileSystem(repo RepoRevSpec) (rwvfs.FileSystem, error)
}

type buildDataService struct {
	client *Client
}

var _ BuildDataService = &buildDataService{}

func (s *buildDataService) FileSystem(repo RepoRevSpec) (rwvfs.FileSystem, error) {
	v := repo.RouteVars()
	v["Path"] = "."
	baseURL, err := s.client.URL(router.RepoBuildDataEntry, v, nil)
	if err != nil {
		return nil, err
	}
	return rwvfs.HTTP(s.client.BaseURL.ResolveReference(baseURL), s.client.httpClient), nil
}

// BuildDataFileSpec specifies a new or existing build data file in a
// repository.
type BuildDataFileSpec struct {
	RepoRev RepoRevSpec
	Path    string
}

// RouteVars returns route variables used to construct URLs to a build
// data file.
func (s *BuildDataFileSpec) RouteVars() map[string]string {
	m := s.RepoRev.RouteVars()
	m["Path"] = s.Path
	return m
}

// GetBuildDataFile is a helper function that calls Stat and Open on
// the FileSystem returned for file's RepoRevSpec. Callers are
// responsible for closing the file (unless an error is returned).
func GetBuildDataFile(s BuildDataService, file BuildDataFileSpec) (io.ReadCloser, os.FileInfo, error) {
	fs, err := s.FileSystem(file.RepoRev)
	if err != nil {
		return nil, nil, err
	}
	fi, err := fs.Stat(file.Path)
	if err != nil {
		return nil, nil, err
	}
	f, err := fs.Open(file.Path)
	if err != nil {
		return nil, fi, err
	}
	return f, fi, err
}

var _ BuildDataService = &MockBuildDataService{}
