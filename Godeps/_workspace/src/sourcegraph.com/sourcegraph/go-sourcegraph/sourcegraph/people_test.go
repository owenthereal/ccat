package sourcegraph

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
)

func TestPersonShortName(t *testing.T) {
	tests := []struct {
		person        Person
		wantShortName string
	}{
		{
			person:        Person{PersonSpec: PersonSpec{Login: "a"}},
			wantShortName: "a",
		},
		{
			person:        Person{PersonSpec: PersonSpec{Login: "a", Email: "x@x.com"}},
			wantShortName: "a",
		},
		{
			person:        Person{PersonSpec: PersonSpec{Email: "x@x.com"}},
			wantShortName: "x",
		},
		{
			person:        Person{PersonSpec: PersonSpec{Email: ""}},
			wantShortName: "(anonymous)",
		},
		{
			person:        Person{PersonSpec: PersonSpec{Email: "x"}},
			wantShortName: "(anonymous)",
		},
	}
	for _, test := range tests {
		n := test.person.ShortName()
		if n != test.wantShortName {
			t.Errorf("%v: got ShortName == %q, want %q", test.person, n, test.wantShortName)
		}
	}
}

func TestPersonSpec(t *testing.T) {
	tests := []struct {
		str  string
		spec PersonSpec
	}{
		{"a", PersonSpec{Login: "a"}},
		{"a@a.com", PersonSpec{Email: "a@a.com"}},
		{"$1", PersonSpec{UID: 1}},
	}

	for _, test := range tests {
		spec, err := ParsePersonSpec(test.str)
		if err != nil {
			t.Errorf("%q: ParsePersonSpec failed: %s", test.str, err)
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

func TestPeopleService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := &Person{FullName: "n"}

	var called bool
	mux.HandleFunc(urlPath(t, router.Person, map[string]string{"PersonSpec": "a"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	person_, _, err := client.People.Get(PersonSpec{Login: "a"})
	if err != nil {
		t.Errorf("People.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(person_, want) {
		t.Errorf("People.Get returned %+v, want %+v", person_, want)
	}
}
