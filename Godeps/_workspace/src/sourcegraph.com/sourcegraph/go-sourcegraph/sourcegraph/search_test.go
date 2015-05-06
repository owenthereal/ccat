package sourcegraph

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
)

func TestSearchService_Search(t *testing.T) {
	setup()
	defer teardown()

	want := &SearchResults{
		ResolvedTokens: Tokens{},
	}

	var called bool
	mux.HandleFunc(urlPath(t, router.Search, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"q":       "q",
			"People":  "false",
			"Repos":   "false",
			"Defs":    "false",
			"Tree":    "false",
			"PerPage": "1",
			"Page":    "2",
		})

		writeJSON(w, want)
	})

	results, _, err := client.Search.Search(&SearchOptions{
		Query:       "q",
		ListOptions: ListOptions{PerPage: 1, Page: 2},
	})
	if err != nil {
		t.Errorf("Search.Search returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(results, want) {
		t.Errorf("Search.Search returned %+v, want %+v", results, want)
	}
}

func TestSearchService_Complete(t *testing.T) {
	setup()
	defer teardown()

	want := &Completions{
		TokenCompletions: Tokens{Term("x")},
		ResolvedTokens:   Tokens{},
	}

	var called bool
	mux.HandleFunc(urlPath(t, router.SearchComplete, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"String":         "abc",
			"InsertionPoint": "2",
		})

		writeJSON(w, want)
	})

	comps, _, err := client.Search.Complete(RawQuery{String: "abc", InsertionPoint: 2})
	if err != nil {
		t.Errorf("Search.Complete returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(comps, want) {
		t.Errorf("Search.Complete returned %+v, want %+v", comps, want)
	}
}

func TestSearchService_Suggest(t *testing.T) {
	setup()
	defer teardown()

	want := []*Suggestion{
		{
			Query:       Tokens{Term("x")},
			Description: "d",
		},
	}

	var called bool
	mux.HandleFunc(urlPath(t, router.SearchSuggestions, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"String":         "abc",
			"InsertionPoint": "2",
		})

		writeJSON(w, want)
	})

	suggs, _, err := client.Search.Suggest(RawQuery{String: "abc", InsertionPoint: 2})
	if err != nil {
		t.Errorf("Search.Suggest returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(suggs, want) {
		t.Errorf("Search.Suggest returned %+v, want %+v", suggs, want)
	}
}

func asJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
