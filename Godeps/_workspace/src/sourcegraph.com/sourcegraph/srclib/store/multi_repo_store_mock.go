package store

import (
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/graph"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/unit"
)

type MockMultiRepoStore struct {
	Repos_    func(...RepoFilter) ([]string, error)
	Versions_ func(...VersionFilter) ([]*Version, error)
	Units_    func(...UnitFilter) ([]*unit.SourceUnit, error)
	Defs_     func(...DefFilter) ([]*graph.Def, error)
	Refs_     func(...RefFilter) ([]*graph.Ref, error)
}

func (m MockMultiRepoStore) Repos(f ...RepoFilter) ([]string, error) {
	return m.Repos_(f...)
}

func (m MockMultiRepoStore) Versions(f ...VersionFilter) ([]*Version, error) {
	return m.Versions_(f...)
}

func (m MockMultiRepoStore) Units(f ...UnitFilter) ([]*unit.SourceUnit, error) {
	return m.Units_(f...)
}

func (m MockMultiRepoStore) Defs(f ...DefFilter) ([]*graph.Def, error) {
	return m.Defs_(f...)
}

func (m MockMultiRepoStore) Refs(f ...RefFilter) ([]*graph.Ref, error) {
	return m.Refs_(f...)
}

var _ MultiRepoStore = MockMultiRepoStore{}
