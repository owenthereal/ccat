package sourcegraph

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/sourcegraph/go-github/github"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
)

func TestReposService_GetCombinedStatus(t *testing.T) {
	setup()
	defer teardown()

	want := &CombinedStatus{CombinedStatus: github.CombinedStatus{State: github.String("s")}}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoCombinedStatus, map[string]string{"RepoSpec": "r.com/x", "Rev": "r"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	cs, _, err := client.Repos.GetCombinedStatus(RepoRevSpec{RepoSpec: RepoSpec{URI: "r.com/x"}, Rev: "r"})
	if err != nil {
		t.Errorf("Repos.GetCombinedStatus returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(cs, want) {
		t.Errorf("Repos.GetCombinedStatus returned %+v, want %+v", cs, want)
	}
}

func TestReposService_CreateStatus(t *testing.T) {
	setup()
	defer teardown()

	want := RepoStatus{RepoStatus: github.RepoStatus{State: github.String("s")}}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoStatusCreate, map[string]string{"RepoSpec": "r.com/x", "Rev": "r"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")

		var st RepoStatus
		if err := json.NewDecoder(r.Body).Decode(&st); err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(st, want) {
			t.Errorf("got status %+v, want %+v", st, want)
		}

		writeJSON(w, want)
	})

	s, _, err := client.Repos.CreateStatus(RepoRevSpec{RepoSpec: RepoSpec{URI: "r.com/x"}, Rev: "r"}, want)
	if err != nil {
		t.Errorf("Repos.CreateStatus returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(s, &want) {
		t.Errorf("Repos.CreateStatus returned %+v, want %+v", s, &want)
	}
}
