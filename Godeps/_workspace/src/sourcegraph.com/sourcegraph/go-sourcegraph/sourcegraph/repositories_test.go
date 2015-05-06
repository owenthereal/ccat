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

func TestRepoSpec(t *testing.T) {
	tests := []struct {
		str  string
		spec RepoSpec
	}{
		{"a.com/x", RepoSpec{URI: "a.com/x"}},
		{"R$1", RepoSpec{RID: 1}},
	}

	for _, test := range tests {
		spec, err := ParseRepoSpec(test.str)
		if err != nil {
			t.Errorf("%q: ParseRepoSpec failed: %s", test.str, err)
			continue
		}
		if spec != test.spec {
			t.Errorf("%q: got spec %+v, want %+v", test.str, spec, test.spec)
			continue
		}

		str := test.spec.PathComponent()
		if str != test.str {
			t.Errorf("%+v: got str %q, want %q", test.spec, str, test.str)
			continue
		}

		spec2, err := UnmarshalRepoSpec(test.spec.RouteVars())
		if err != nil {
			t.Errorf("%+v: UnmarshalRepoSpec: %s", test.spec, err)
			continue
		}
		if spec2 != test.spec {
			t.Errorf("%q: got spec %+v, want %+v", test.str, spec, test.spec)
			continue
		}
	}
}

func TestRepoRevSpec(t *testing.T) {
	tests := []struct {
		spec      RepoRevSpec
		routeVars map[string]string
	}{
		{RepoRevSpec{RepoSpec: RepoSpec{URI: "a.com/x"}, Rev: "r"}, map[string]string{"RepoSpec": "a.com/x", "Rev": "r"}},
		{RepoRevSpec{RepoSpec: RepoSpec{RID: 123}, Rev: "r"}, map[string]string{"RepoSpec": "R$123", "Rev": "r"}},
		{RepoRevSpec{RepoSpec: RepoSpec{URI: "a.com/x"}, Rev: "r", CommitID: "c"}, map[string]string{"RepoSpec": "a.com/x", "Rev": "r===c"}},
	}

	for _, test := range tests {
		routeVars := test.spec.RouteVars()
		if !reflect.DeepEqual(routeVars, test.routeVars) {
			t.Errorf("got route vars %+v, want %+v", routeVars, test.routeVars)
		}
		spec, err := UnmarshalRepoRevSpec(routeVars)
		if err != nil {
			t.Errorf("UnmarshalRepoRevSpec(%+v): %s", routeVars, err)
			continue
		}
		if spec != test.spec {
			t.Errorf("got spec %+v, want %+v", spec, test.spec)
		}
	}
}

func TestReposService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := &Repo{RID: 1}

	var called bool
	mux.HandleFunc(urlPath(t, router.Repo, map[string]string{"RepoSpec": "r.com/x"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	repo_, _, err := client.Repos.Get(RepoSpec{URI: "r.com/x"}, nil)
	if err != nil {
		t.Errorf("Repos.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normRepo(want)
	if !reflect.DeepEqual(repo_, want) {
		t.Errorf("Repos.Get returned %+v, want %+v", repo_, want)
	}
}

func TestReposService_GetStats(t *testing.T) {
	setup()
	defer teardown()

	want := RepoStats{"x": 1, "y": 2}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoStats, map[string]string{"RepoSpec": "r.com/x", "Rev": "c"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	stats, _, err := client.Repos.GetStats(RepoRevSpec{RepoSpec: RepoSpec{URI: "r.com/x"}, Rev: "c"})
	if err != nil {
		t.Errorf("Repos.GetStats returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(stats, want) {
		t.Errorf("Repos.GetStats returned %+v, want %+v", stats, want)
	}
}

func TestReposService_GetOrCreate(t *testing.T) {
	setup()
	defer teardown()

	want := &Repo{RID: 1}

	var called bool
	mux.HandleFunc(urlPath(t, router.ReposGetOrCreate, map[string]string{"RepoSpec": "r.com/x"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")

		writeJSON(w, want)
	})

	repo_, _, err := client.Repos.GetOrCreate(RepoSpec{URI: "r.com/x"}, nil)
	if err != nil {
		t.Errorf("Repos.GetOrCreate returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normRepo(want)
	if !reflect.DeepEqual(repo_, want) {
		t.Errorf("Repos.GetOrCreate returned %+v, want %+v", repo_, want)
	}
}

func TestReposService_GetSettings(t *testing.T) {
	setup()
	defer teardown()

	want := &RepoSettings{}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoSettings, map[string]string{"RepoSpec": "r.com/x"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	settings, _, err := client.Repos.GetSettings(RepoSpec{URI: "r.com/x"})
	if err != nil {
		t.Errorf("Repos.GetSettings returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(settings, want) {
		t.Errorf("Repos.GetSettings returned %+v, want %+v", settings, want)
	}
}

func TestReposService_UpdateSettings(t *testing.T) {
	setup()
	defer teardown()

	want := RepoSettings{}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoSettings, map[string]string{"RepoSpec": "r.com/x"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
		testBody(t, r, `{}`+"\n")

		writeJSON(w, want)
	})

	_, err := client.Repos.UpdateSettings(RepoSpec{URI: "r.com/x"}, want)
	if err != nil {
		t.Errorf("Repos.UpdateSettings returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}
}

func TestReposService_RefreshProfile(t *testing.T) {
	setup()
	defer teardown()

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoRefreshProfile, map[string]string{"RepoSpec": "r.com/x"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
	})

	_, err := client.Repos.RefreshProfile(RepoSpec{URI: "r.com/x"})
	if err != nil {
		t.Errorf("Repos.RefreshProfile returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}
}

func TestReposService_RefreshVCSData(t *testing.T) {
	setup()
	defer teardown()

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoRefreshVCSData, map[string]string{"RepoSpec": "r.com/x"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
	})

	_, err := client.Repos.RefreshVCSData(RepoSpec{URI: "r.com/x"})
	if err != nil {
		t.Errorf("Repos.RefreshVCSData returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}
}

func TestReposService_ComputeStats(t *testing.T) {
	setup()
	defer teardown()

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoComputeStats, map[string]string{"RepoSpec": "r.com/x", "Rev": "c"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
	})

	_, err := client.Repos.ComputeStats(RepoRevSpec{RepoSpec: RepoSpec{URI: "r.com/x"}, Rev: "c"})
	if err != nil {
		t.Errorf("Repos.ComputeStats returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}
}

func TestReposService_GetBuild(t *testing.T) {
	setup()
	defer teardown()

	want := &RepoBuildInfo{
		Exact:                &Build{BID: 1},
		LastSuccessful:       &Build{BID: 2},
		CommitsBehind:        3,
		LastSuccessfulCommit: &Commit{Commit: &vcs.Commit{Message: "m"}},
	}
	normalizeTime(&want.LastSuccessfulCommit.Author.Date)
	normalizeBuildTime(want.Exact)
	normalizeBuildTime(want.LastSuccessful)

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoBuild, map[string]string{"RepoSpec": "r.com/x", "Rev": "r"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	buildInfo, _, err := client.Repos.GetBuild(RepoRevSpec{RepoSpec: RepoSpec{URI: "r.com/x"}, Rev: "r"}, nil)
	if err != nil {
		t.Errorf("Repos.GetBuild returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&buildInfo.LastSuccessfulCommit.Author.Date)
	normalizeBuildTime(buildInfo.Exact)
	normalizeBuildTime(buildInfo.LastSuccessful)
	if !reflect.DeepEqual(buildInfo.Exact, want.Exact) {
		t.Errorf("Repos.GetBuild returned %+v, want %+v", buildInfo.Exact, want.Exact)
	}
}

func TestReposService_Create(t *testing.T) {
	setup()
	defer teardown()

	newRepo := NewRepoSpec{Type: "git", CloneURLStr: "http://r.com/x"}
	want := &Repo{RID: 1}

	var called bool
	mux.HandleFunc(urlPath(t, router.ReposCreate, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")
		testBody(t, r, `{"Type":"git","CloneURL":"http://r.com/x"}`+"\n")

		writeJSON(w, want)
	})

	repo_, _, err := client.Repos.Create(newRepo)
	if err != nil {
		t.Errorf("Repos.Create returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normRepo(want)
	if !reflect.DeepEqual(repo_, want) {
		t.Errorf("Repos.Create returned %+v, want %+v", repo_, want)
	}
}

func TestReposService_GetReadme(t *testing.T) {
	setup()
	defer teardown()

	want := &vcsclient.TreeEntry{Name: "hello"}
	want.ModTime = want.ModTime.In(time.UTC)

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoReadme, map[string]string{"RepoSpec": "r.com/x"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	readme, _, err := client.Repos.GetReadme(RepoRevSpec{RepoSpec: RepoSpec{URI: "r.com/x"}})
	if err != nil {
		t.Errorf("Repos.GetReadme returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(readme, want) {
		t.Errorf("Repos.GetReadme returned %+v, want %+v", readme, want)
	}
}

func TestReposService_List(t *testing.T) {
	setup()
	defer teardown()

	want := []*Repo{&Repo{RID: 1}}

	var called bool
	mux.HandleFunc(urlPath(t, router.Repos, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"URIs":      "a,b",
			"Name":      "n",
			"Owner":     "o",
			"Sort":      "name",
			"Direction": "asc",
			"NoFork":    "true",
			"PerPage":   "1",
			"Page":      "2",
		})

		writeJSON(w, want)
	})

	repos, _, err := client.Repos.List(&RepoListOptions{
		URIs:        []string{"a", "b"},
		Name:        "n",
		Owner:       "o",
		Sort:        "name",
		Direction:   "asc",
		NoFork:      true,
		ListOptions: ListOptions{PerPage: 1, Page: 2},
	})
	if err != nil {
		t.Errorf("Repos.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normRepo(want...)
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Repos.List: got %+v, want %+v", repos, want)
	}
}

func TestReposService_ListCommits(t *testing.T) {
	setup()
	defer teardown()

	want := []*Commit{{Commit: &vcs.Commit{Message: "m"}}}
	normTime(want[0])

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoCommits, map[string]string{"RepoSpec": "r.com/x"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"Head": "myhead"})

		writeJSON(w, want)
	})

	commits, _, err := client.Repos.ListCommits(RepoSpec{URI: "r.com/x"}, &RepoListCommitsOptions{Head: "myhead"})
	if err != nil {
		t.Errorf("Repos.ListCommits returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(commits, want) {
		t.Errorf("Repos.ListCommits returned %+v, want %+v", commits, want)
	}
}

func TestReposService_GetCommit(t *testing.T) {
	setup()
	defer teardown()

	want := &Commit{Commit: &vcs.Commit{Message: "m"}}
	normTime(want)

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoCommit, map[string]string{"RepoSpec": "r.com/x", "Rev": "r"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	commit, _, err := client.Repos.GetCommit(RepoRevSpec{RepoSpec: RepoSpec{URI: "r.com/x"}, Rev: "r"}, nil)
	if err != nil {
		t.Errorf("Repos.GetCommit returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(commit, want) {
		t.Errorf("Repos.GetCommit returned %+v, want %+v", commit, want)
	}
}

func TestReposService_ListBranches(t *testing.T) {
	setup()
	defer teardown()

	want := []*vcs.Branch{{Name: "b", Head: "c"}}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoBranches, map[string]string{"RepoSpec": "r.com/x"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	branches, _, err := client.Repos.ListBranches(RepoSpec{URI: "r.com/x"}, nil)
	if err != nil {
		t.Errorf("Repos.ListBranches returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(branches, want) {
		t.Errorf("Repos.ListBranches returned %+v, want %+v", branches, want)
	}
}

func TestReposService_ListTags(t *testing.T) {
	setup()
	defer teardown()

	want := []*vcs.Tag{{Name: "t", CommitID: "c"}}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoTags, map[string]string{"RepoSpec": "r.com/x"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	tags, _, err := client.Repos.ListTags(RepoSpec{URI: "r.com/x"}, nil)
	if err != nil {
		t.Errorf("Repos.ListTags returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(tags, want) {
		t.Errorf("Repos.ListTags returned %+v, want %+v", tags, want)
	}
}

func TestReposService_ListBadges(t *testing.T) {
	setup()
	defer teardown()

	want := []*Badge{{Name: "b"}}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoBadges, map[string]string{"RepoSpec": "r.com/x"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	badges, _, err := client.Repos.ListBadges(RepoSpec{URI: "r.com/x"})
	if err != nil {
		t.Errorf("Repos.ListBadges returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(badges, want) {
		t.Errorf("Repos.ListBadges returned %+v, want %+v", badges, want)
	}
}

func TestReposService_ListCounters(t *testing.T) {
	setup()
	defer teardown()

	want := []*Counter{{Name: "b"}}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoCounters, map[string]string{"RepoSpec": "r.com/x"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	counters, _, err := client.Repos.ListCounters(RepoSpec{URI: "r.com/x"})
	if err != nil {
		t.Errorf("Repos.ListCounters returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(counters, want) {
		t.Errorf("Repos.ListCounters returned %+v, want %+v", counters, want)
	}
}

func TestReposService_ListAuthors(t *testing.T) {
	setup()
	defer teardown()

	want := []*AugmentedRepoAuthor{{Person: &Person{FullName: "b"}}}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoAuthors, map[string]string{"RepoSpec": "r.com/x", "Rev": "c"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	authors, _, err := client.Repos.ListAuthors(RepoRevSpec{RepoSpec: RepoSpec{URI: "r.com/x"}, Rev: "c"}, nil)
	if err != nil {
		t.Errorf("Repos.ListAuthors returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(authors, want) {
		t.Errorf("Repos.ListAuthors returned %+v, want %+v", authors, want)
	}
}

func TestReposService_ListClients(t *testing.T) {
	setup()
	defer teardown()

	want := []*AugmentedRepoClient{{Person: &Person{FullName: "b"}}}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoClients, map[string]string{"RepoSpec": "r.com/x"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	clients, _, err := client.Repos.ListClients(RepoSpec{URI: "r.com/x"}, nil)
	if err != nil {
		t.Errorf("Repos.ListClients returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(clients, want) {
		t.Errorf("Repos.ListClients returned %+v, want %+v", clients, want)
	}
}

func TestReposService_ListDependents(t *testing.T) {
	setup()
	defer teardown()

	want := []*AugmentedRepoDependent{{Repo: &Repo{URI: "r2"}}}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoDependents, map[string]string{"RepoSpec": "r.com/x"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	dependents, _, err := client.Repos.ListDependents(RepoSpec{URI: "r.com/x"}, nil)
	if err != nil {
		t.Errorf("Repos.ListDependents returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normRepo(want[0].Repo)
	if !reflect.DeepEqual(dependents, want) {
		t.Errorf("Repos.ListDependents returned %+v, want %+v", dependents, want)
	}
}

func TestReposService_ListDependencies(t *testing.T) {
	setup()
	defer teardown()

	want := []*AugmentedRepoDependency{{Repo: &Repo{URI: "r2"}}}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoDependencies, map[string]string{"RepoSpec": "r.com/x", "Rev": "c"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	dependencies, _, err := client.Repos.ListDependencies(RepoRevSpec{RepoSpec: RepoSpec{URI: "r.com/x"}, Rev: "c"}, nil)
	if err != nil {
		t.Errorf("Repos.ListDependencies returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normRepo(want[0].Repo)
	if !reflect.DeepEqual(dependencies, want) {
		t.Errorf("Repos.ListDependencies returned %+v, want %+v", dependencies, want)
	}
}

func TestReposService_ListByContributor(t *testing.T) {
	setup()
	defer teardown()

	want := []*AugmentedRepoContribution{{Repo: &Repo{URI: "r.com/x"}}}

	var called bool
	mux.HandleFunc(urlPath(t, router.UserRepoContributions, map[string]string{"UserSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"NoFork": "true"})

		writeJSON(w, want)
	})

	repos, _, err := client.Repos.ListByContributor(UserSpec{Login: "a"}, &RepoListByContributorOptions{NoFork: true})
	if err != nil {
		t.Errorf("Repos.ListByContributor returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normRepo(want[0].Repo)
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Repos.ListByContributor returned %+v, want %+v", repos, want)
	}
}

func TestReposService_ListByClient(t *testing.T) {
	setup()
	defer teardown()

	want := []*AugmentedRepoUsageByClient{{DefRepo: &Repo{URI: "r.com/x"}}}

	var called bool
	mux.HandleFunc(urlPath(t, router.UserRepoDependencies, map[string]string{"UserSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	repos, _, err := client.Repos.ListByClient(UserSpec{Login: "a"}, nil)
	if err != nil {
		t.Errorf("Repos.ListByClient returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normRepo(want[0].DefRepo)
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Repos.ListByClient returned %+v, want %+v", repos, want)
	}
}

func TestReposService_ListByRefdAuthor(t *testing.T) {
	setup()
	defer teardown()

	want := []*AugmentedRepoUsageOfAuthor{{Repo: &Repo{URI: "r.com/x"}}}

	var called bool
	mux.HandleFunc(urlPath(t, router.UserRepoDependents, map[string]string{"UserSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	repos, _, err := client.Repos.ListByRefdAuthor(UserSpec{Login: "a"}, nil)
	if err != nil {
		t.Errorf("Repos.ListByRefdAuthor returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normRepo(want[0].Repo)
	if !reflect.DeepEqual(repos, want) {
		t.Errorf("Repos.ListByRefdAuthor returned %+v, want %+v", repos, want)
	}
}

func normTime(c *Commit) {
	c.Author.Date = c.Author.Date.In(time.UTC)
	if c.Committer != nil {
		c.Committer.Date = c.Committer.Date.In(time.UTC)
	}
}

func normRepo(r ...*Repo) {
	for _, r := range r {
		r.CreatedAt = r.CreatedAt.UTC()
		r.UpdatedAt = r.UpdatedAt.UTC()
		r.PushedAt = r.PushedAt.UTC()
	}
}
