package sourcegraph

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/sourcegraph/go-github/github"
	"strings"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
	"github.com/kr/pretty"
)

func TestIssues(t *testing.T) {
	tests := []struct {
		spec          IssueSpec
		wantRouteVars map[string]string
	}{
		{
			spec:          IssueSpec{Repo: RepoSpec{URI: "foo.com/bar"}, Number: 1},
			wantRouteVars: map[string]string{"RepoSpec": "foo.com/bar", "Issue": "1"},
		},
	}

	for _, test := range tests {
		routeVars := test.spec.RouteVars()
		if !reflect.DeepEqual(routeVars, test.wantRouteVars) {
			t.Errorf("Got route vars %+v, but wanted %+v", routeVars, test.wantRouteVars)
		}

		spec, err := UnmarshalIssueSpec(test.wantRouteVars)
		if err != nil {
			t.Errorf("UnmarshalIssueSpec(%+v): %s", test.wantRouteVars, err)
		}
		if !reflect.DeepEqual(spec, test.spec) {
			t.Errorf("Got spec %+v, but wanted %+v", spec, test.spec)
		}
	}
}

func TestIssuesService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := &Issue{Issue: github.Issue{Number: github.Int(1)}}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoIssue, map[string]string{"RepoSpec": "r.com/x", "Issue": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	issue, _, err := client.Issues.Get(IssueSpec{Repo: RepoSpec{URI: "r.com/x"}, Number: 1}, nil)
	if err != nil {
		t.Errorf("Issues.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(issue, want) {
		t.Errorf("Issues.Get returned %+v, want %+v", issue, want)
	}
}

func TestIssuesService_ListByRepo(t *testing.T) {
	setup()
	defer teardown()

	want := []*Issue{&Issue{Issue: github.Issue{Number: github.Int(1)}}}
	repoSpec := RepoSpec{URI: "x.com/r"}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoIssues, repoSpec.RouteVars()), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"PerPage": "1",
			"Page":    "2",
		})

		writeJSON(w, want)
	})

	issues, _, err := client.Issues.ListByRepo(
		repoSpec,
		&IssueListOptions{
			ListOptions: ListOptions{PerPage: 1, Page: 2},
		},
	)
	if err != nil {
		t.Errorf("Issues.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(issues, want) {
		t.Errorf("Issues.List returned %+v, want %+v with diff: %s", issues, want, strings.Join(pretty.Diff(want, issues), "\n"))
	}
}

func TestIssuesService_ListComments(t *testing.T) {
	setup()
	defer teardown()

	want := []*IssueComment{&IssueComment{IssueComment: github.IssueComment{ID: github.Int(1)}}}
	issueSpec := IssueSpec{Repo: RepoSpec{URI: "r.com/x"}, Number: 1}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoIssueComments, issueSpec.RouteVars()), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"PerPage": "1",
			"Page":    "2",
		})

		writeJSON(w, want)
	})

	comments, _, err := client.Issues.ListComments(
		issueSpec,
		&IssueListCommentsOptions{
			ListOptions: ListOptions{PerPage: 1, Page: 2},
		},
	)
	if err != nil {
		t.Errorf("Issues.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(comments, want) {
		t.Errorf("Issues.List returned %+v, want %+v with diff: %s", comments, want, strings.Join(pretty.Diff(want, comments), "\n"))
	}
}
