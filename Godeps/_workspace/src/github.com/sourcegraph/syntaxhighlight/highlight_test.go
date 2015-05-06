package syntaxhighlight

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/sourcegraph/annotate"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/vcsstore/vcsclient"
	"github.com/kr/pretty"
)

var saveExp = flag.Bool("exp", false, "overwrite all expected output files with actual output (returning a failure)")
var match = flag.String("m", "", "only run tests whose name contains this string")

func TestAsHTML(t *testing.T) {
	dir := "testdata"
	tests, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		name := test.Name()
		if !strings.Contains(name, *match) {
			continue
		}
		if strings.HasSuffix(name, ".html") {
			continue
		}
		if name == "net_http_client.go" {
			// only use this file for benchmarking
			continue
		}
		path := filepath.Join(dir, name)
		input, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatal(err)
			continue
		}

		got, err := AsHTML(input)
		if err != nil {
			t.Errorf("%s: AsHTML: %s", name, err)
			continue
		}

		expPath := path + ".html"
		if *saveExp {
			err = ioutil.WriteFile(expPath, got, 0700)
			if err != nil {
				t.Fatal(err)
			}
			continue
		}

		want, err := ioutil.ReadFile(expPath)
		if err != nil {
			t.Fatal(err)
		}

		want = bytes.TrimSpace(want)
		got = bytes.TrimSpace(got)

		if !bytes.Equal(want, got) {
			t.Errorf("%s:\nwant ==========\n%q\ngot ===========\n%q", name, want, got)
			continue
		}
	}

	if *saveExp {
		t.Fatal("overwrote all expected output files with actual output (run tests again without -exp)")
	}
}

func TestAnnotate(t *testing.T) {
	src := []byte(`a:=2`)
	want := annotate.Annotations{
		{Start: 0, End: 1, Left: []byte(`<span class="pln">`), Right: []byte("</span>")},
		{Start: 1, End: 2, Left: []byte(`<span class="pun">`), Right: []byte("</span>")},
		{Start: 2, End: 3, Left: []byte(`<span class="pun">`), Right: []byte("</span>")},
		{Start: 3, End: 4, Left: []byte(`<span class="dec">`), Right: []byte("</span>")},
	}
	got, err := Annotate(src, HTMLAnnotator(DefaultHTMLConfig))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %# v, got %# v\n\ndiff:\n%v", pretty.Formatter(want), pretty.Formatter(got), strings.Join(pretty.Diff(got, want), "\n"))
		for _, g := range got {
			t.Logf("%+v  %q  LEFT=%q RIGHT=%q", g, src[g.Start:g.End], g.Left, g.Right)
		}
	}
}

// codeEquals tests the equality between the given SourceCode entry and an
// array of lines containing arrays of tokens as their string representation.
func codeEquals(code *sourcegraph.SourceCode, want [][]string) bool {
	if len(code.Lines) != len(want) {
		return false
	}
	for i, line := range code.Lines {
		for j, t := range line.Tokens {
			switch t := t.(type) {
			case *sourcegraph.SourceCodeToken:
				if t.Label != want[i][j] {
					return false
				}
			case string:
				if t != want[i][j] {
					return false
				}
			}
		}
	}
	return true
}

func TestCodeEquals(t *testing.T) {
	for _, tt := range []struct {
		code *sourcegraph.SourceCode
		want [][]string
	}{
		{
			code: &sourcegraph.SourceCode{
				Lines: []*sourcegraph.SourceCodeLine{
					&sourcegraph.SourceCodeLine{
						Tokens: []interface{}{
							&sourcegraph.SourceCodeToken{Label: "a"},
							&sourcegraph.SourceCodeToken{Label: "b"},
							"c",
							&sourcegraph.SourceCodeToken{Label: "d"},
							"e",
						},
					},
					&sourcegraph.SourceCodeLine{},
					&sourcegraph.SourceCodeLine{
						Tokens: []interface{}{
							"c",
						},
					},
				},
			},
			want: [][]string{[]string{"a", "b", "c", "d", "e"}, []string{}, []string{"c"}},
		},
	} {
		if !codeEquals(tt.code, tt.want) {
			t.Errorf("Expected: %# v, Got: %# v\n", tt.code, tt.want)
		}
	}
}

func newFileWithRange(src []byte) *vcsclient.FileWithRange {
	return &vcsclient.FileWithRange{
		TreeEntry: &vcsclient.TreeEntry{Contents: []byte(src)},
		FileRange: vcsclient.FileRange{StartByte: 0, EndByte: int64(len(src))},
	}
}

func TestNilAnnotator_multiLineTokens(t *testing.T) {
	for _, tt := range []struct {
		src  string
		want [][]string
	}{
		{
			src: "/* I am\na multiline\ncomment\n*/",
			want: [][]string{
				[]string{"/* I am"},
				[]string{"a multiline"},
				[]string{"comment"},
				[]string{"*/"},
			},
		},
		{
			src: "a := `I am\na multiline\nstring literal\n`",
			want: [][]string{
				[]string{"a", " ", ":", "=", " ", "`I am"},
				[]string{"a multiline"},
				[]string{"string literal"},
				[]string{"`"},
			},
		},
	} {
		e := newFileWithRange([]byte(tt.src))
		ann := NewNilAnnotator(e)
		_, err := Annotate(e.Contents, ann)
		if err != nil {
			t.Fatal(err)
		}
		if !codeEquals(ann.Code, tt.want) {
			t.Errorf("Expected %# v\n\nGot %# v", tt.want, pretty.Formatter(ann.Code.Lines))
		}
	}
}

func BenchmarkAnnotate(b *testing.B) {
	input, err := ioutil.ReadFile("testdata/net_http_client.go")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Annotate(input[:2000], HTMLAnnotator(DefaultHTMLConfig))
		if err != nil {
			b.Fatal(err)
		}
	}
}
