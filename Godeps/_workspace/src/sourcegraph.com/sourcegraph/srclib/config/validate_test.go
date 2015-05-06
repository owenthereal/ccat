package config

import (
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/unit"
)

func TestTree_validate(t *testing.T) {
	tests := map[string]*Tree{
		"absolute path":                &Tree{SourceUnits: []*unit.SourceUnit{{Files: []string{"/foo"}}}},
		"relative path above root":     &Tree{SourceUnits: []*unit.SourceUnit{{Files: []string{"../foo"}}}},
		"bad path after being cleaned": &Tree{SourceUnits: []*unit.SourceUnit{{Files: []string{"foo/bar/../../../../baz"}}}},
	}

	for label, tree := range tests {
		if err := tree.validate(); err != ErrInvalidFilePath {
			t.Errorf("%s: got err %v, want ErrInvalidFilePath", label, err)
		}
	}
}
