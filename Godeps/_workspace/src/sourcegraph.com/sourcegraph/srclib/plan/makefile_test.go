package plan_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/config"
	_ "github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/config"
	_ "github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/dep"
	_ "github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/grapher"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/plan"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/unit"
	"sourcegraph.com/sourcegraph/makex"
)

func TestCreateMakefile(t *testing.T) {
	buildDataDir := "testdata"
	c := &config.Tree{
		SourceUnits: []*unit.SourceUnit{
			{
				Name:  "n",
				Type:  "t",
				Files: []string{"f"},
				Ops: map[string]*srclib.ToolRef{
					"graph":      {Toolchain: "tc", Subcmd: "t"},
					"depresolve": {Toolchain: "tc", Subcmd: "t"},
				},
			},
		},
	}

	mf, err := plan.CreateMakefile(buildDataDir, nil, "", c, plan.Options{NoCache: true})
	if err != nil {
		t.Fatal(err)
	}

	want := `
all: testdata/n/t.graph.json testdata/n/t.depresolve.json

testdata/n/t.graph.json: testdata/n/t.unit.json f
	src tool  "tc" "t" < $< | src internal normalize-graph-data --unit-type "t" --dir . 1> $@

testdata/n/t.depresolve.json: testdata/n/t.unit.json
	src tool  "tc" "t" < $^ 1> $@

.DELETE_ON_ERROR:
`

	gotBytes, err := makex.Marshal(mf)
	if err != nil {
		t.Fatal(err)
	}

	want = strings.TrimSpace(want)
	got := string(bytes.TrimSpace(gotBytes))

	if got != want {
		t.Errorf("got makefile:\n==========\n%s\n==========\n\nwant makefile:\n==========\n%s\n==========", got, want)
	}
}
