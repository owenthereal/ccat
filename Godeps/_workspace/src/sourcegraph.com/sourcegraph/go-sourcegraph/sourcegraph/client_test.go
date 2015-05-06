package sourcegraph

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
)

// Uses HTTP client testing code adapted from github.com/google/go-github.

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the Sourcegraph client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// sourcegraph client configured to use test server
	client = NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func urlPath(t *testing.T, routeName string, routeVars map[string]string) string {
	url, err := client.URL(routeName, routeVars, nil)
	if err != nil {
		t.Fatalf("Error constructing URL path for route %q with vars %+v: %s", routeName, routeVars, err)
	}
	return url.Path
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		panic("writeJSON: " + err.Error())
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}

	r.ParseForm()
	if !reflect.DeepEqual(want, r.Form) {
		t.Errorf("Request parameters = %v, want %v", r.Form, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Unable to read body")
	}
	str := string(b)
	if want != str {
		t.Errorf("Body = %s, want: %s", str, want)
	}
}

func TestClient_URL(t *testing.T) {
	tests := []struct {
		base      string
		route     string
		routeVars map[string]string
		opt       interface{}
		exp       string
	}{{
		base:      "https://sourcegraph.com/api/",
		route:     router.Repo,
		routeVars: map[string]string{"RepoSpec": "github.com/gorilla/mux"},
		exp:       "https://sourcegraph.com/api/repos/github.com/gorilla/mux",
	}, {
		base:      "https://sourcegraph.com/api",
		route:     router.Repo,
		routeVars: map[string]string{"RepoSpec": "github.com/gorilla/mux"},
		exp:       "https://sourcegraph.com/api/repos/github.com/gorilla/mux",
	}, {
		base:      "http://localhost:3000/api/",
		route:     router.Repo,
		routeVars: map[string]string{"RepoSpec": "github.com/gorilla/mux"},
		exp:       "http://localhost:3000/api/repos/github.com/gorilla/mux",
	}, {
		base:      "http://localhost:3000/api",
		route:     router.Repo,
		routeVars: map[string]string{"RepoSpec": "github.com/gorilla/mux"},
		exp:       "http://localhost:3000/api/repos/github.com/gorilla/mux",
	}}
	for _, test := range tests {
		func() {
			c := NewClient(nil)
			baseURL, err := url.Parse(test.base)
			if err != nil {
				t.Fatal(err)
			}
			c.BaseURL = baseURL
			url, err := c.URL(test.route, test.routeVars, test.opt)
			if err != nil {
				t.Errorf("Error generating URL: %s", err)
				return
			}
			if url.String() != test.exp {
				t.Errorf("Expected %s, got %s on test case %+v", test.exp, url.String(), test)
				return
			}
		}()
	}
}

func normalizeTime(tm *time.Time) {
	*tm = tm.In(time.UTC)
}
