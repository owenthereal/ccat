package sourcegraph

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-vcs/vcs"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/vcsstore/vcsclient"
)

func TestRepoTreeService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := &TreeEntry{
		TreeEntry: &vcsclient.TreeEntry{
			Name:     "p",
			Type:     vcsclient.FileEntry,
			Size:     123,
			Contents: []byte("hello"),
		},
	}
	want.ModTime = want.ModTime.In(time.UTC)

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoTreeEntry, map[string]string{"RepoSpec": "r.com/x", "Rev": "v", "Path": "p"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"Formatted":          "true",
			"ExpandContextLines": "2",
			"StartByte":          "123",
			"EndByte":            "456",
		})

		writeJSON(w, want)
	})

	opt := &RepoTreeGetOptions{
		Formatted: true,
		GetFileOptions: vcsclient.GetFileOptions{
			FileRange:          vcsclient.FileRange{StartByte: 123, EndByte: 456},
			ExpandContextLines: 2,
		},
	}
	data, _, err := client.RepoTree.Get(TreeEntrySpec{
		RepoRev: RepoRevSpec{RepoSpec: RepoSpec{URI: "r.com/x"}, Rev: "v"},
		Path:    "p",
	}, opt)
	if err != nil {
		t.Errorf("RepoTree.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(data, want) {
		t.Errorf("RepoTree.Get returned %+v, want %+v", data, want)
	}
}

func TestRepoTreeService_Search(t *testing.T) {
	setup()
	defer teardown()

	want := []*vcs.SearchResult{
		{
			File:      "f",
			Match:     []byte("abc"),
			StartLine: 1,
			EndLine:   2,
		},
	}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoTreeSearch, map[string]string{"RepoSpec": "r.com/x", "Rev": "v", "CommitID": "c"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"Query":        "q",
			"QueryType":    "t",
			"ContextLines": "1",
			"N":            "3",
			"Offset":       "0",
			"Formatted":    "true",
		})

		writeJSON(w, want)
	})

	opt := RepoTreeSearchOptions{
		Formatted:     true,
		SearchOptions: vcs.SearchOptions{Query: "q", QueryType: "t", N: 3, ContextLines: 1},
	}
	data, _, err := client.RepoTree.Search(RepoRevSpec{RepoSpec: RepoSpec{URI: "r.com/x"}, Rev: "v"}, &opt)
	if err != nil {
		t.Errorf("RepoTree.Search returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(data, want) {
		t.Errorf("RepoTree.Search returned %+v, want %+v", data, want)
	}
}
