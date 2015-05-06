package store

import "github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/unit"

type MockTreeStore struct {
	Units_ func(...UnitFilter) ([]*unit.SourceUnit, error)
	MockUnitStore
}

func (m MockTreeStore) Units(f ...UnitFilter) ([]*unit.SourceUnit, error) {
	return m.Units_(f...)
}

var _ TreeStore = MockTreeStore{}
