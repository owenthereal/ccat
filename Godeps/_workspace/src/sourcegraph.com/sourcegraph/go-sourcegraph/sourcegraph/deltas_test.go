package sourcegraph

import (
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-diff/diff"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/graph"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/unit"
	"github.com/kr/pretty"
)

var (
	baseRev = RepoRevSpec{RepoSpec: RepoSpec{URI: "base.com/repo"}, Rev: "baserev", CommitID: "basecommit"}
	headRev = RepoRevSpec{RepoSpec: RepoSpec{URI: "head.com/repo"}, Rev: "headrev", CommitID: "headcommit"}
)

func TestDeltas(t *testing.T) {
	tests := []struct {
		spec          DeltaSpec
		wantRouteVars map[string]string
	}{
		{
			spec: DeltaSpec{
				Base: RepoRevSpec{RepoSpec: RepoSpec{URI: "samerepo"}, Rev: "baserev", CommitID: "basecommit"},
				Head: RepoRevSpec{RepoSpec: RepoSpec{URI: "samerepo"}, Rev: "headrev", CommitID: "headcommit"},
			},
			wantRouteVars: map[string]string{
				"RepoSpec":     "samerepo",
				"Rev":          "baserev===basecommit",
				"DeltaHeadRev": "headrev===headcommit",
			},
		},
		{
			spec: DeltaSpec{
				Base: baseRev,
				Head: headRev,
			},
			wantRouteVars: map[string]string{
				"RepoSpec":     "base.com/repo",
				"Rev":          "baserev===basecommit",
				"DeltaHeadRev": encodeCrossRepoRevSpecForDeltaHeadRev(headRev),
			},
		},
	}
	for _, test := range tests {
		vars := test.spec.RouteVars()
		if !reflect.DeepEqual(vars, test.wantRouteVars) {
			t.Errorf("got route vars != want\n\n%s", strings.Join(pretty.Diff(vars, test.wantRouteVars), "\n"))
		}

		spec, err := UnmarshalDeltaSpec(vars)
		if err != nil {
			t.Errorf("UnmarshalDeltaSpec(%+v): %s", err)
			continue
		}
		if !reflect.DeepEqual(spec, test.spec) {
			t.Errorf("got spec != original spec\n\n%s", strings.Join(pretty.Diff(spec, test.spec), "\n"))
		}
	}
}

func TestDeltasService_Get(t *testing.T) {
	setup()
	defer teardown()

	ds := DeltaSpec{
		Base: baseRev,
		Head: headRev,
	}
	want := &Delta{
		Base: ds.Base,
		Head: ds.Head,
	}

	var called bool
	mux.HandleFunc(urlPath(t, router.Delta, ds.RouteVars()), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	delta, _, err := client.Deltas.Get(ds, nil)
	if err != nil {
		t.Errorf("Deltas.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(delta, want) {
		t.Errorf("Deltas.Get returned %+v, want %+v", delta, want)
	}
}
func TestDeltasService_ListUnits(t *testing.T) {
	setup()
	defer teardown()

	ds := DeltaSpec{
		Base: baseRev,
		Head: headRev,
	}
	want := []*UnitDelta{{Head: &unit.SourceUnit{Type: "t", Name: "u"}}}

	var called bool
	mux.HandleFunc(urlPath(t, router.DeltaUnits, ds.RouteVars()), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})

		writeJSON(w, want)
	})

	units, _, err := client.Deltas.ListUnits(ds, &DeltaListUnitsOptions{})
	if err != nil {
		t.Errorf("Deltas.ListUnits returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(units, want) {
		t.Errorf("Deltas.ListUnits returned %+v, want %+v", units, want)
	}
}

func TestDeltasService_ListDefs(t *testing.T) {
	setup()
	defer teardown()

	ds := DeltaSpec{
		Base: baseRev,
		Head: headRev,
	}
	want := &DeltaDefs{Defs: []*DefDelta{{Base: &Def{Def: graph.Def{Name: "x"}}}}}

	var called bool
	mux.HandleFunc(urlPath(t, router.DeltaDefs, ds.RouteVars()), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"UnitType": "t",
			"Unit":     "u",
		})

		writeJSON(w, want)
	})

	defs, _, err := client.Deltas.ListDefs(ds, &DeltaListDefsOptions{DeltaFilter: DeltaFilter{UnitType: "t", Unit: "u"}})
	if err != nil {
		t.Errorf("Deltas.ListDefs returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(defs, want) {
		t.Errorf("Deltas.ListDefs returned %+v, want %+v", defs, want)
	}
}

func TestDeltasService_ListDependencies(t *testing.T) {
	setup()
	defer teardown()

	ds := DeltaSpec{
		Base: baseRev,
		Head: headRev,
	}
	want := &DeltaDependencies{}

	var called bool
	mux.HandleFunc(urlPath(t, router.DeltaDependencies, ds.RouteVars()), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"UnitType": "t",
			"Unit":     "u",
		})

		writeJSON(w, want)
	})

	dependencies, _, err := client.Deltas.ListDependencies(ds, &DeltaListDependenciesOptions{DeltaFilter: DeltaFilter{UnitType: "t", Unit: "u"}})
	if err != nil {
		t.Errorf("Deltas.ListDependencies returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(dependencies, want) {
		t.Errorf("Deltas.ListDependencies returned %+v, want %+v", dependencies, want)
	}
}

func TestDeltasService_ListFiles(t *testing.T) {
	setup()
	defer teardown()

	ds := DeltaSpec{
		Base: baseRev,
		Head: headRev,
	}
	want := &DeltaFiles{FileDiffs: []*diff.FileDiff{{OrigName: "o"}}}

	var called bool
	mux.HandleFunc(urlPath(t, router.DeltaFiles, ds.RouteVars()), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"UnitType": "t",
			"Unit":     "u",
		})

		writeJSON(w, want)
	})

	files, _, err := client.Deltas.ListFiles(ds, &DeltaListFilesOptions{DeltaFilter: DeltaFilter{UnitType: "t", Unit: "u"}})
	if err != nil {
		t.Errorf("Deltas.ListFiles returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(files, want) {
		t.Errorf("Deltas.ListFiles returned %+v, want %+v", files, want)
	}
}

func TestDeltasService_ListAffectedAuthors(t *testing.T) {
	setup()
	defer teardown()

	ds := DeltaSpec{
		Base: baseRev,
		Head: headRev,
	}
	want := []*DeltaAffectedPerson{}

	var called bool
	mux.HandleFunc(urlPath(t, router.DeltaAffectedAuthors, ds.RouteVars()), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"UnitType": "t",
			"Unit":     "u",
		})

		writeJSON(w, want)
	})

	affectedAuthors, _, err := client.Deltas.ListAffectedAuthors(ds, &DeltaListAffectedAuthorsOptions{DeltaFilter: DeltaFilter{UnitType: "t", Unit: "u"}})
	if err != nil {
		t.Errorf("Deltas.ListAffectedAuthors returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(affectedAuthors, want) {
		t.Errorf("Deltas.ListAffectedAuthors returned %+v, want %+v", affectedAuthors, want)
	}
}

func TestDeltasService_ListAffectedClients(t *testing.T) {
	setup()
	defer teardown()

	ds := DeltaSpec{
		Base: baseRev,
		Head: headRev,
	}
	want := []*DeltaAffectedPerson{}

	var called bool
	mux.HandleFunc(urlPath(t, router.DeltaAffectedClients, ds.RouteVars()), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"UnitType": "t",
			"Unit":     "u",
		})

		writeJSON(w, want)
	})

	affectedClients, _, err := client.Deltas.ListAffectedClients(ds, &DeltaListAffectedClientsOptions{DeltaFilter: DeltaFilter{UnitType: "t", Unit: "u"}})
	if err != nil {
		t.Errorf("Deltas.ListAffectedClients returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(affectedClients, want) {
		t.Errorf("Deltas.ListAffectedClients returned %+v, want %+v", affectedClients, want)
	}
}

func TestDeltasService_ListAffectedDependents(t *testing.T) {
	setup()
	defer teardown()

	ds := DeltaSpec{
		Base: baseRev,
		Head: headRev,
	}
	want := []*DeltaAffectedRepo{}

	var called bool
	mux.HandleFunc(urlPath(t, router.DeltaAffectedDependents, ds.RouteVars()), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"UnitType": "t",
			"Unit":     "u",
		})

		writeJSON(w, want)
	})

	affectedDependents, _, err := client.Deltas.ListAffectedDependents(ds, &DeltaListAffectedDependentsOptions{DeltaFilter: DeltaFilter{UnitType: "t", Unit: "u"}})
	if err != nil {
		t.Errorf("Deltas.ListAffectedDependents returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(affectedDependents, want) {
		t.Errorf("Deltas.ListAffectedDependents returned %+v, want %+v", affectedDependents, want)
	}
}

func TestDeltasService_ListReviewers(t *testing.T) {
	setup()
	defer teardown()

	ds := DeltaSpec{
		Base: baseRev,
		Head: headRev,
	}
	want := []*DeltaReviewer{}

	var called bool
	mux.HandleFunc(urlPath(t, router.DeltaReviewers, ds.RouteVars()), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"UnitType": "t",
			"Unit":     "u",
		})

		writeJSON(w, want)
	})

	reviewers, _, err := client.Deltas.ListReviewers(ds, &DeltaListReviewersOptions{DeltaFilter: DeltaFilter{UnitType: "t", Unit: "u"}})
	if err != nil {
		t.Errorf("Deltas.ListReviewers returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(reviewers, want) {
		t.Errorf("Deltas.ListReviewers returned %+v, want %+v", reviewers, want)
	}
}

func TestDeltasService_ListIncoming(t *testing.T) {
	setup()
	defer teardown()

	rr := RepoRevSpec{RepoSpec: RepoSpec{URI: "x.com/r"}, Rev: "r"}
	want := []*Delta{}

	var called bool
	mux.HandleFunc(urlPath(t, router.DeltasIncoming, rr.RouteVars()), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"UnitType": "t",
			"Unit":     "u",
		})

		writeJSON(w, want)
	})

	incoming, _, err := client.Deltas.ListIncoming(rr, &DeltaListIncomingOptions{DeltaFilter: DeltaFilter{UnitType: "t", Unit: "u"}})
	if err != nil {
		t.Errorf("Deltas.ListIncoming returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(incoming, want) {
		t.Errorf("Deltas.ListIncoming returned %+v, want %+v", incoming, want)
	}
}
