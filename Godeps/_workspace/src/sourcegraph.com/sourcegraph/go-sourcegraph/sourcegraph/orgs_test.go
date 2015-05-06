package sourcegraph

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
)

func TestOrgSpec(t *testing.T) {
	tests := []struct {
		str  string
		spec OrgSpec
	}{
		{"a", OrgSpec{Org: "a"}},
		{"$1", OrgSpec{UID: 1}},
	}

	for _, test := range tests {
		spec, err := ParseOrgSpec(test.str)
		if err != nil {
			t.Errorf("%q: ParseOrgSpec failed: %s", test.str, err)
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

func TestOrgsService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := &Org{User: User{UID: 1}}

	var called bool
	mux.HandleFunc(urlPath(t, router.Org, map[string]string{"OrgSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	org_, _, err := client.Orgs.Get(OrgSpec{Org: "a"})
	if err != nil {
		t.Errorf("Orgs.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(org_, want) {
		t.Errorf("Orgs.Get returned %+v, want %+v", org_, want)
	}
}

func TestOrgsService_ListMembers(t *testing.T) {
	setup()
	defer teardown()

	want := []*User{{UID: 1}}

	var called bool
	mux.HandleFunc(urlPath(t, router.OrgMembers, map[string]string{"OrgSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	members, _, err := client.Orgs.ListMembers(OrgSpec{Org: "a"}, nil)
	if err != nil {
		t.Errorf("Orgs.ListMembers returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(members, want) {
		t.Errorf("Orgs.ListMembers returned %+v, want %+v", members, want)
	}
}

func TestOrgsService_GetSettings(t *testing.T) {
	setup()
	defer teardown()

	want := &OrgSettings{}

	var called bool
	mux.HandleFunc(urlPath(t, router.OrgSettings, map[string]string{"OrgSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	settings, _, err := client.Orgs.GetSettings(OrgSpec{Org: "a"})
	if err != nil {
		t.Errorf("Orgs.GetSettings returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(settings, want) {
		t.Errorf("Orgs.GetSettings returned %+v, want %+v", settings, want)
	}
}

func TestOrgsService_UpdateSettings(t *testing.T) {
	setup()
	defer teardown()

	want := OrgSettings{}

	var called bool
	mux.HandleFunc(urlPath(t, router.OrgSettings, map[string]string{"OrgSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
		testBody(t, r, `{}`+"\n")

		writeJSON(w, want)
	})

	_, err := client.Orgs.UpdateSettings(OrgSpec{Org: "a"}, want)
	if err != nil {
		t.Errorf("Orgs.UpdateSettings returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}
}
