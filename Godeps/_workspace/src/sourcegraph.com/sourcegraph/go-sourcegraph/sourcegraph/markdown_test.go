package sourcegraph

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
)

func TestMarkdown(t *testing.T) {
	setup()
	defer teardown()

	input := MarkdownRequestBody{
		Markdown:    []byte(`raw markdown`),
		MarkdownOpt: MarkdownOpt{EnableCheckboxes: true},
	}
	want := &MarkdownData{
		Rendered: []byte(`i am rendered`),
		Checklist: &Checklist{
			Todo: 2,
			Done: 1,
		},
	}

	var called bool
	mux.HandleFunc(urlPath(t, router.Markdown, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")

		var m MarkdownRequestBody
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(m, input) {
			t.Fatalf("got input %+v, expected %+v", m, input)
		}

		writeJSON(w, want)
	})

	got, _, err := client.Markdown.Render(input.Markdown, input.MarkdownOpt)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Errorf("!called")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
