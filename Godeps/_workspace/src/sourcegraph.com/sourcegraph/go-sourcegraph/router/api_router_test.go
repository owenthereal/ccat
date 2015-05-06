package router

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/sourcegraph/mux"
	"github.com/kr/pretty"
)

func TestMatch(t *testing.T) {
	router := NewAPIRouter(nil)
	tests := []struct {
		path          string
		wantNoMatch   bool
		wantRouteName string
		wantVars      map[string]string
		wantPath      string
	}{
		// Repo
		{
			path:          "/repos/repohost.com/foo",
			wantRouteName: Repo,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo"},
		},
		{
			path:          "/repos/a/b/c",
			wantRouteName: Repo,
			wantVars:      map[string]string{"RepoSpec": "a/b/c"},
		},
		{
			path:          "/repos/a.com/b",
			wantRouteName: Repo,
			wantVars:      map[string]string{"RepoSpec": "a.com/b"},
		},
		{
			path:          "/repos/R$123",
			wantRouteName: Repo,
			wantVars:      map[string]string{"RepoSpec": "R$123"},
		},
		{
			path:        "/repos/a.com/b@mycommitid", // doesn't accept a commit ID
			wantNoMatch: true,
		},
		{
			path:        "/repos/.invalidrepo",
			wantNoMatch: true,
		},

		// Repo sub-routes
		{
			path:          "/repos/repohost.com/foo/.authors",
			wantRouteName: RepoAuthors,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo"},
		},
		{
			path:          "/repos/R$123/.authors",
			wantRouteName: RepoAuthors,
			wantVars:      map[string]string{"RepoSpec": "R$123"},
		},
		{
			path:          "/repos/repohost.com/foo@myrev/.authors",
			wantRouteName: RepoAuthors,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "myrev"},
		},
		{
			path:          "/repos/repohost.com/foo@myrev==abcd/.authors",
			wantRouteName: RepoAuthors,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "myrev==abcd"},
		},
		{
			path:          "/repos/repohost.com/foo@myrev/subrev/.authors",
			wantRouteName: RepoAuthors,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "myrev/subrev"},
		},
		{
			path:          "/repos/repohost.com/foo@myrev/subrev1/subrev2/.authors",
			wantRouteName: RepoAuthors,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "myrev/subrev1/subrev2"},
		},
		{
			path:          "/repos/repohost.com/foo@myrev/subrev==abcd/.authors",
			wantRouteName: RepoAuthors,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "myrev/subrev==abcd"},
		},
		{
			path:          "/repos/repohost.com/foo@releases/1.0rc/.authors",
			wantRouteName: RepoAuthors,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "releases/1.0rc"},
		},

		// Repo sub-routes that don't allow an "@REVSPEC" revision.
		{
			path:          "/repos/repohost.com/foo/.dependents",
			wantRouteName: RepoDependents,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo"},
		},
		{
			path:        "/repos/repohost.com/foo@myrevspec/.dependents", // no @REVSPEC match
			wantNoMatch: true,
		},
		{
			path:          "/repos/repohost.com/foo/.commits",
			wantRouteName: RepoCommits,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo"},
		},
		{
			path:          "/repos/repohost.com/foo/.commits/123abc",
			wantRouteName: RepoCommit,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "123abc"},
		},
		{
			path:          "/repos/repohost.com/foo/.commits/123abc/.compare",
			wantRouteName: RepoCompareCommits,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "123abc"},
		},
		{
			path:          "/repos/repohost.com/foo/.commits/123abc/xyz/.compare",
			wantRouteName: RepoCompareCommits,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "123abc/xyz"},
		},

		// Repo tree
		{
			path:          "/repos/repohost.com/foo@mycommitid/.tree",
			wantRouteName: RepoTreeEntry,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "mycommitid", "Path": "."},
		},
		{
			path:          "/repos/repohost.com/foo@mycommitid/.tree/",
			wantRouteName: RepoTreeEntry,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "mycommitid", "Path": "."},
			wantPath:      "/repos/repohost.com/foo@mycommitid/.tree",
		},
		{
			path:          "/repos/repohost.com/foo@mycommitid/.tree/my/file",
			wantRouteName: RepoTreeEntry,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "mycommitid", "Path": "my/file"},
		},
		{
			path:          "/repos/repohost.com/foo@myrev/subrev/.tree/my/file",
			wantRouteName: RepoTreeEntry,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "myrev/subrev", "Path": "my/file"},
		},

		// Units
		{
			path:          "/.units",
			wantRouteName: Units,
			wantVars:      map[string]string{},
		},
		{
			path:          "/repos/repohost.com/foo@mycommitid/.units/t/u",
			wantRouteName: Unit,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "mycommitid", "UnitType": "t", "Unit": "u"},
		},

		// Repo build data
		{
			path:          "/repos/repohost.com/foo/.build-data",
			wantRouteName: RepoBuildDataEntry,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Path": "."},
		},
		{
			path:          "/repos/repohost.com/foo@mycommitid/.build-data/",
			wantRouteName: RepoBuildDataEntry,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "mycommitid", "Path": "."},
			wantPath:      "/repos/repohost.com/foo@mycommitid/.build-data",
		},
		{
			path:          "/repos/repohost.com/foo@mycommitid/.build-data/my/file",
			wantRouteName: RepoBuildDataEntry,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "mycommitid", "Path": "my/file"},
		},

		// Defs
		{
			path:          "/repos/repohost.com/foo@mycommitid/.defs/.t/.def/p",
			wantRouteName: Def,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "UnitType": "t", "Unit": ".", "Path": "p", "Rev": "mycommitid"},
		},
		{
			path:          "/repos/repohost.com/foo@myrev/mysubrev/.defs/.t/.def/p",
			wantRouteName: Def,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "UnitType": "t", "Unit": ".", "Path": "p", "Rev": "myrev/mysubrev"},
		},
		{
			path:          "/repos/repohost.com/foo/.defs/.t/.def/p",
			wantRouteName: Def,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "UnitType": "t", "Unit": ".", "Path": "p"},
		},
		{
			path:          "/repos/repohost.com/foo/.defs/.t/.def", // empty path
			wantRouteName: Def,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "UnitType": "t", "Unit": ".", "Path": "."},
		},
		{
			path:          "/repos/repohost.com/foo/.defs/.t/u1/.def/p",
			wantRouteName: Def,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "UnitType": "t", "Unit": "u1", "Path": "p"},
		},
		{
			path:          "/repos/repohost.com/foo/.defs/.t/u1/u2/.def/p1/p2",
			wantRouteName: Def,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "UnitType": "t", "Unit": "u1/u2", "Path": "p1/p2"},
		},

		// Def sub-routes
		{
			path:          "/repos/repohost.com/foo/.defs/.t/.def/p/.authors",
			wantRouteName: DefAuthors,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "UnitType": "t", "Unit": ".", "Path": "p"},
		},
		{
			path:          "/repos/repohost.com/foo/.defs/.t/.def/.authors", // empty path
			wantRouteName: DefAuthors,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "UnitType": "t", "Unit": ".", "Path": "."},
		},
		{
			path:          "/repos/repohost.com/foo/.defs/.t/u1/u2/.def/p1/p2/.authors",
			wantRouteName: DefAuthors,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "UnitType": "t", "Unit": "u1/u2", "Path": "p1/p2"},
		},

		// Deltas
		{
			path:          "/repos/repohost.com/foo/.deltas/branch1..branch2",
			wantRouteName: Delta,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "branch1", "DeltaHeadRev": "branch2"},
		},
		{
			path:          "/repos/repohost.com/foo/.deltas/a/b/c..x/y/z",
			wantRouteName: Delta,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "a/b/c", "DeltaHeadRev": "x/y/z"},
		},
		{
			path:          "/repos/repohost.com/foo/.deltas/branch1..branch2/.reviewers",
			wantRouteName: DeltaReviewers,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "branch1", "DeltaHeadRev": "branch2"},
		},
		{
			path:          "/repos/repohost.com/foo/.deltas/branch1..branch2===4739/.reviewers",
			wantRouteName: DeltaReviewers,
			wantVars:      map[string]string{"RepoSpec": "repohost.com/foo", "Rev": "branch1", "DeltaHeadRev": "branch2===4739"},
		},

		// User
		{
			path:          "/users/alice",
			wantRouteName: User,
			wantVars:      map[string]string{"UserSpec": "alice"},
		},
		{
			path:          "/users/$1",
			wantRouteName: User,
			wantVars:      map[string]string{"UserSpec": "$1"},
		},
		{
			path:        "/users/alice@example.com",
			wantNoMatch: true,
		},
		{
			path:        "/users/alice@-x-yJAANTud-iAVVw==",
			wantNoMatch: true,
		},

		// Person
		{
			path:          "/people/alice",
			wantRouteName: Person,
			wantVars:      map[string]string{"PersonSpec": "alice"},
		},
		{
			path:          "/people/$1",
			wantRouteName: Person,
			wantVars:      map[string]string{"PersonSpec": "$1"},
		},
		{
			path:          "/people/alice@example.com",
			wantRouteName: Person,
			wantVars:      map[string]string{"PersonSpec": "alice@example.com"},
		},
		{
			path:          "/people/alice@-x-yJAANTud-iAVVw==",
			wantRouteName: Person,
			wantVars:      map[string]string{"PersonSpec": "alice@-x-yJAANTud-iAVVw=="},
		},
	}
	for _, test := range tests {
		var routeMatch mux.RouteMatch
		match := router.Match(&http.Request{Method: "GET", URL: &url.URL{Path: test.path}}, &routeMatch)

		if match && test.wantNoMatch {
			t.Errorf("%s: got match (route %q), want no match", test.path, routeMatch.Route.GetName())
		}
		if !match && !test.wantNoMatch {
			t.Errorf("%s: got no match, wanted match", test.path)
		}
		if !match || test.wantNoMatch {
			continue
		}

		if routeName := routeMatch.Route.GetName(); routeName != test.wantRouteName {
			t.Errorf("%s: got matched route %q, want %q", test.path, routeName, test.wantRouteName)
		}

		if diff := pretty.Diff(routeMatch.Vars, test.wantVars); len(diff) > 0 {
			t.Errorf("%s: vars don't match expected:\n%s", test.path, strings.Join(diff, "\n"))
		}

		// Check that building the URL yields the original path.
		var pairs []string
		for k, v := range test.wantVars {
			pairs = append(pairs, k, v)
		}
		path, err := routeMatch.Route.URLPath(pairs...)
		if err != nil {
			t.Errorf("%s: URLPath(%v) failed: %s", test.path, pairs, err)
			continue
		}
		var wantPath string
		if test.wantPath != "" {
			wantPath = test.wantPath
		} else {
			wantPath = test.path
		}
		if path.Path != wantPath {
			t.Errorf("got generated path %q, want %q", path, wantPath)
		}
	}
}
