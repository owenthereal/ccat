package sourcegraph

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/schema"
	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/sourcegraph/go-github/github"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
	"github.com/kr/pretty"
)

func TestPullRequestsService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := &PullRequest{PullRequest: github.PullRequest{Number: github.Int(1)}}
	opts := &PullRequestGetOptions{Checklist: true}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoPullRequest, map[string]string{"RepoSpec": "r.com/x", "Pull": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		var gotOpts PullRequestGetOptions
		schema.NewDecoder().Decode(&gotOpts, r.URL.Query())
		if !reflect.DeepEqual(&gotOpts, opts) {
			t.Errorf("got requested opts %+v, but got %+v", &gotOpts, opts)
		}

		writeJSON(w, want)
	})

	pull, _, err := client.PullRequests.Get(PullRequestSpec{Repo: RepoSpec{URI: "r.com/x"}, Number: 1}, opts)
	if err != nil {
		t.Errorf("PullRequests.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(pull, want) {
		t.Errorf("PullRequests.Get returned %+v, want %+v", pull, want)
	}
}

func TestPullRequestsService_ListByRepo(t *testing.T) {
	setup()
	defer teardown()

	want := []*PullRequest{&PullRequest{PullRequest: github.PullRequest{Number: github.Int(1)}}}
	repoSpec := RepoSpec{URI: "x.com/r"}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoPullRequests, repoSpec.RouteVars()), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"PerPage": "1",
			"Page":    "2",
		})

		writeJSON(w, want)
	})

	pulls, _, err := client.PullRequests.ListByRepo(
		repoSpec,
		&PullRequestListOptions{
			ListOptions: ListOptions{PerPage: 1, Page: 2},
		},
	)
	if err != nil {
		t.Errorf("PullRequests.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(pulls, want) {
		t.Errorf("PullRequests.List returned %+v, want %+v with diff: %s", pulls, want, strings.Join(pretty.Diff(want, pulls), "\n"))
	}
}

func TestPullRequestsService_ListComments(t *testing.T) {
	setup()
	defer teardown()

	want := []*PullRequestComment{&PullRequestComment{PullRequestComment: github.PullRequestComment{ID: github.Int(1)}}}
	pullSpec := PullRequestSpec{Repo: RepoSpec{URI: "r.com/x"}, Number: 1}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoPullRequestComments, pullSpec.RouteVars()), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"PerPage": "1",
			"Page":    "2",
		})

		writeJSON(w, want)
	})

	comments, _, err := client.PullRequests.ListComments(
		pullSpec,
		&PullRequestListCommentsOptions{
			ListOptions: ListOptions{PerPage: 1, Page: 2},
		},
	)
	if err != nil {
		t.Errorf("PullRequests.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(comments, want) {
		t.Errorf("PullRequests.List returned %+v, want %+v with diff: %s", comments, want, strings.Join(pretty.Diff(want, comments), "\n"))
	}
}

func TestPullRequestsService_CreateComment(t *testing.T) {
	setup()
	defer teardown()

	pullSpec := PullRequestSpec{Repo: RepoSpec{URI: "r.com/foo"}, Number: 22}
	comment := PullRequestComment{
		PullRequestComment: github.PullRequestComment{
			Body:      github.String("this is a comment"),
			Path:      github.String("/"),
			Position:  github.Int(2),
			CommitID:  github.String("54be46135e45be9bd3318b8fd39a456ff1e2895e"),
			User:      &github.User{},
			CreatedAt: timePtr(time.Unix(100, 100).UTC()),
			UpdatedAt: timePtr(time.Unix(200, 200).UTC()),
		},
	}
	wantComment := comment
	wantComment.ID = github.Int(1)

	called := false
	mux.HandleFunc(urlPath(t, router.RepoPullRequestCommentsCreate, pullSpec.RouteVars()), func(w http.ResponseWriter, req *http.Request) {
		called = true
		testMethod(t, req, "POST")

		var unmarshalled PullRequestComment
		err := json.NewDecoder(req.Body).Decode(&unmarshalled)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(unmarshalled, comment) {
			t.Errorf("Got unmarshalled comment %+v, want %+v", unmarshalled, comment)
		}

		writeJSON(w, wantComment)
	})

	gotComment, _, err := client.PullRequests.CreateComment(pullSpec, &comment)
	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Errorf("!called")
	}

	if !jsonEqual(t, gotComment, wantComment) {
		t.Errorf("Got %+v, want %+v", gotComment, wantComment)
	}
}

func TestPullRequestsService_EditComment(t *testing.T) {
	setup()
	defer teardown()

	pullSpec := PullRequestSpec{Repo: RepoSpec{URI: "r.com/foo"}, Number: 22}
	comment := PullRequestComment{
		PullRequestComment: github.PullRequestComment{
			ID:        github.Int(1),
			Body:      github.String("this is a comment"),
			Path:      github.String("/"),
			Position:  github.Int(2),
			CommitID:  github.String("54be46135e45be9bd3318b8fd39a456ff1e2895e"),
			User:      &github.User{},
			CreatedAt: timePtr(time.Unix(100, 100).UTC()),
			UpdatedAt: timePtr(time.Unix(200, 200).UTC()),
		},
	}

	called := false
	mux.HandleFunc(urlPath(t, router.RepoPullRequestCommentsEdit, PullRequestCommentSpec{Pull: pullSpec, Comment: *comment.ID}.RouteVars()), func(w http.ResponseWriter, req *http.Request) {
		called = true
		testMethod(t, req, "PATCH")

		var unmarshalled PullRequestComment
		err := json.NewDecoder(req.Body).Decode(&unmarshalled)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(unmarshalled, comment) {
			t.Errorf("Got unmarshalled comment %+v, want %+v", unmarshalled, comment)
		}

		writeJSON(w, comment)
	})

	updatedComment, _, err := client.PullRequests.EditComment(pullSpec, &comment)
	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Errorf("!called")
	}

	if !jsonEqual(t, updatedComment, comment) {
		t.Errorf("Got %+v, want %+v", updatedComment, comment)
	}
}

func TestPullRequestsService_EditComment_ErrorNoID(t *testing.T) {
	setup()
	defer teardown()

	pullSpec := PullRequestSpec{Repo: RepoSpec{URI: "r.com/foo"}, Number: 22}
	comment := PullRequestComment{
		PullRequestComment: github.PullRequestComment{
			Body:      github.String("this is a comment"),
			Path:      github.String("/"),
			Position:  github.Int(2),
			CommitID:  github.String("54be46135e45be9bd3318b8fd39a456ff1e2895e"),
			User:      &github.User{},
			CreatedAt: timePtr(time.Unix(100, 100).UTC()),
			UpdatedAt: timePtr(time.Unix(200, 200).UTC()),
		},
	}

	_, _, err := client.PullRequests.EditComment(pullSpec, &comment)
	if err.Error() != "comment ID not specified" {
		t.Errorf(`expected error "comment ID not specified", but got none`)
	}
}

func TestPullRequestsService_DeleteComment(t *testing.T) {
	setup()
	defer teardown()

	pullSpec := PullRequestSpec{Repo: RepoSpec{URI: "r.com/foo"}, Number: 22}
	commentID := 1
	commentSpec := PullRequestCommentSpec{Pull: pullSpec, Comment: commentID}

	called := false
	mux.HandleFunc(urlPath(t, router.RepoPullRequestCommentsEdit, commentSpec.RouteVars()), func(w http.ResponseWriter, req *http.Request) {
		called = true
		testMethod(t, req, "DELETE")
	})

	_, err := client.PullRequests.DeleteComment(pullSpec, commentID)
	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Errorf("!called")
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func jsonEqual(t *testing.T, u, v interface{}) bool {
	uj, err := json.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}
	vj, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.Equal(uj, vj)
}

func TestPullRequestsService_Merge(t *testing.T) {
	setup()
	defer teardown()

	pullSpec := PullRequestSpec{Repo: RepoSpec{URI: "r.com/foo"}, Number: 22}
	mergeRequest := &PullRequestMergeRequest{CommitMessage: "message"}
	wantMergeResult := &PullRequestMergeResult{github.PullRequestMergeResult{
		SHA:     github.String("42c066aecc289359dcaaff61e9a85396ed7c376f"),
		Merged:  github.Bool(true),
		Message: github.String("Pull Request successfully merged"),
	}}

	called := false
	mux.HandleFunc(urlPath(t, router.RepoPullRequestMerge, pullSpec.RouteVars()), func(w http.ResponseWriter, req *http.Request) {
		called = true
		testMethod(t, req, "PUT")

		var unmarshalled PullRequestMergeRequest
		err := json.NewDecoder(req.Body).Decode(&unmarshalled)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(&unmarshalled, mergeRequest) {
			t.Errorf("Got unmarshalled merge request %+v, want %+v", unmarshalled, mergeRequest)
		}

		writeJSON(w, wantMergeResult)
	})

	mergeResult, _, err := client.PullRequests.Merge(pullSpec, mergeRequest)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(mergeResult, wantMergeResult) {
		t.Errorf("got %+v, want %+v", mergeResult, wantMergeResult)
	}
}
