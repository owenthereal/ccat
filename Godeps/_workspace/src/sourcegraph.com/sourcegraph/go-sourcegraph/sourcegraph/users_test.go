package sourcegraph

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
)

func TestUserSpec(t *testing.T) {
	tests := []struct {
		str       string
		spec      UserSpec
		wantError bool
	}{
		{"a", UserSpec{Login: "a"}, false},
		{"$1", UserSpec{UID: 1}, false},
		{"a@a.com", UserSpec{}, true},
	}

	for _, test := range tests {
		spec, err := ParseUserSpec(test.str)
		if err != nil && !test.wantError {
			t.Errorf("%q: ParseUserSpec failed: %s", test.str, err)
		}
		if test.wantError && err == nil {
			t.Errorf("%q: ParseUserSpec returned nil error, want non-nil error", test.str)
			continue
		}
		if err != nil {
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
	}
}

func TestUsersService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := &User{UID: 1}

	var called bool
	mux.HandleFunc(urlPath(t, router.User, map[string]string{"UserSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	user_, _, err := client.Users.Get(UserSpec{Login: "a"}, nil)
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(user_, want) {
		t.Errorf("Users.Get returned %+v, want %+v", user_, want)
	}
}

func TestUsersService_GetSettings(t *testing.T) {
	setup()
	defer teardown()

	// Test success.
	want := &UserSettings{}

	var called bool
	mux.HandleFunc(urlPath(t, router.UserSettings, map[string]string{"UserSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	settings, _, err := client.Users.GetSettings(UserSpec{Login: "a"})
	if err != nil {
		t.Errorf("Users.GetSettings returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(settings, want) {
		t.Errorf("Users.GetSettings returned %+v, want %+v", settings, want)
	}

	// Test failure.
	expectErr := func(p UserSpec) {
		_, _, err = client.Users.GetSettings(p)
		if err == nil {
			t.Error("Expected GetSettings to error for %v.", p)
		}
	}
	expectErr(UserSpec{UID: 1000})
	expectErr(UserSpec{Login: "doesnotexist"})
}

func TestUsersService_UpdateSettings(t *testing.T) {
	setup()
	defer teardown()

	want := UserSettings{}

	var called bool
	mux.HandleFunc(urlPath(t, router.UserSettings, map[string]string{"UserSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
		wantBody, _ := json.Marshal(want)
		testBody(t, r, string(wantBody)+"\n")
	})

	_, err := client.Users.UpdateSettings(UserSpec{Login: "a"}, want)
	if err != nil {
		t.Errorf("Users.UpdateSettings returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}
}

func TestUsersService_ListEmails(t *testing.T) {
	setup()
	defer teardown()

	want := []*EmailAddr{{Email: "a@a.com", Verified: true, Primary: true}}

	var called bool
	mux.HandleFunc(urlPath(t, router.UserEmails, map[string]string{"UserSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	emails, _, err := client.Users.ListEmails(UserSpec{Login: "a"})
	if err != nil {
		t.Errorf("Users.ListEmails returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(emails, want) {
		t.Errorf("Users.ListEmails returned %+v, want %+v", emails, want)
	}
}

func TestUsersService_GetOrCreateFromGitHub(t *testing.T) {
	setup()
	defer teardown()

	want := &User{UID: 1, Login: "a"}

	var called bool
	mux.HandleFunc(urlPath(t, router.UserFromGitHub, map[string]string{"GitHubUserSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	user_, _, err := client.Users.GetOrCreateFromGitHub(GitHubUserSpec{Login: "a"}, nil)
	if err != nil {
		t.Errorf("Users.GetOrCreateFromGitHub returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(user_, want) {
		t.Errorf("Users.GetOrCreateFromGitHub returned %+v, want %+v", user_, want)
	}
}

func TestUsersService_RefreshProfile(t *testing.T) {
	setup()
	defer teardown()

	var called bool
	mux.HandleFunc(urlPath(t, router.UserRefreshProfile, map[string]string{"UserSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
	})

	_, err := client.Users.RefreshProfile(UserSpec{Login: "a"})
	if err != nil {
		t.Errorf("Users.RefreshProfile returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}
}

func TestUsersService_ComputeStats(t *testing.T) {
	setup()
	defer teardown()

	var called bool
	mux.HandleFunc(urlPath(t, router.UserComputeStats, map[string]string{"UserSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
	})

	_, err := client.Users.ComputeStats(UserSpec{Login: "a"})
	if err != nil {
		t.Errorf("Users.ComputeStats returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}
}

func TestUsersService_List(t *testing.T) {
	setup()
	defer teardown()

	want := []*User{{UID: 1}}

	var called bool
	mux.HandleFunc(urlPath(t, router.Users, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"Query":     "nl",
			"Sort":      "name",
			"Direction": "asc",
			"PerPage":   "1",
			"Page":      "2",
		})

		writeJSON(w, want)
	})

	users, _, err := client.Users.List(&UsersListOptions{
		Query:       "nl",
		Sort:        "name",
		Direction:   "asc",
		ListOptions: ListOptions{PerPage: 1, Page: 2},
	})
	if err != nil {
		t.Errorf("Users.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(users, want) {
		t.Errorf("Users.List returned %+v, want %+v", users, want)
	}
}

func TestUsersService_ListAuthors(t *testing.T) {
	setup()
	defer teardown()

	want := []*AugmentedPersonUsageByClient{{Author: &Person{FullName: "n"}}}

	var called bool
	mux.HandleFunc(urlPath(t, router.UserAuthors, map[string]string{"UserSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	authors, _, err := client.Users.ListAuthors(UserSpec{Login: "a"}, nil)
	if err != nil {
		t.Errorf("Users.ListAuthors returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(authors, want) {
		t.Errorf("Users.ListAuthors returned %+v, want %+v", authors, want)
	}
}

func TestUsersService_ListClients(t *testing.T) {
	setup()
	defer teardown()

	want := []*AugmentedPersonUsageOfAuthor{{Client: &Person{FullName: "n"}}}

	var called bool
	mux.HandleFunc(urlPath(t, router.UserClients, map[string]string{"UserSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	clients, _, err := client.Users.ListClients(UserSpec{Login: "a"}, nil)
	if err != nil {
		t.Errorf("Users.ListClients returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(clients, want) {
		t.Errorf("Users.ListClients returned %+v, want %+v", clients, want)
	}
}

func TestUsersService_ListOrgs(t *testing.T) {
	setup()
	defer teardown()

	want := []*Org{{User: User{Login: "o"}}}

	var called bool
	mux.HandleFunc(urlPath(t, router.UserOrgs, map[string]string{"UserSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	orgs, _, err := client.Users.ListOrgs(UserSpec{Login: "a"}, nil)
	if err != nil {
		t.Errorf("Users.ListOrgs returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(orgs, want) {
		t.Errorf("Users.ListOrgs returned %+v, want %+v", orgs, want)
	}
}
