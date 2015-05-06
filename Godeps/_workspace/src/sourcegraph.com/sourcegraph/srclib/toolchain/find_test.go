package toolchain

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib"
)

func TestList_program(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "srclib-toolchain-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	defer func(orig string) {
		srclib.Path = orig
	}(srclib.Path)
	srclib.Path = tmpdir

	files := map[string]os.FileMode{
		// ok
		"a/a/.bin/a":          0700,
		"a/a/Srclibtoolchain": 0700,

		// not executable
		"b/b/.bin/z":          0600,
		"b/b/Srclibtoolchain": 0600,

		// not in .bin
		"c/c/c":               0700,
		"c/c/Srclibtoolchain": 0700,
	}
	for f, mode := range files {
		f = filepath.Join(tmpdir, f)
		if err := os.MkdirAll(filepath.Dir(f), 0700); err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(f, nil, mode); err != nil {
			t.Fatal(err)
		}
	}

	// Put a file symlink in srclib DIR path.
	oldp := filepath.Join(tmpdir, "a/a/.bin/a")
	newp := filepath.Join(tmpdir, "link")
	if err := os.Symlink(oldp, newp); err != nil {
		t.Fatal(err)
	}

	toolchains, err := List()
	if err != nil {
		t.Fatal(err)
	}

	got := toolchainPathsWithProgramOrDockerfile(toolchains)
	want := []string{"a/a"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got toolchains %v, want %v", got, want)
	}
}

func TestList_docker(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "srclib-toolchain-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	defer func(orig string) {
		srclib.Path = orig
	}(srclib.Path)
	srclib.Path = tmpdir

	files := map[string]struct{}{
		// ok
		"a/a/Dockerfile": struct{}{}, "a/a/Srclibtoolchain": struct{}{},

		// no Srclibtoolchain
		"b/b/Dockerfile": struct{}{},

		// ok (no Dockerfile)
		"c/c/Srclibtoolchain": struct{}{},
	}
	for f := range files {
		if err := os.MkdirAll(filepath.Join(tmpdir, filepath.Dir(f)), 0700); err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(filepath.Join(tmpdir, f), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}

	toolchains, err := List()
	if err != nil {
		t.Fatal(err)
	}

	got := toolchainPathsWithProgramOrDockerfile(toolchains)
	want := []string{"a/a"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got toolchains %v, want %v", got, want)
	}
}

func toolchainPathsWithProgramOrDockerfile(toolchains []*Info) []string {
	paths := make([]string, 0, len(toolchains))
	for _, toolchain := range toolchains {
		if toolchain.Program != "" || toolchain.Dockerfile != "" {
			paths = append(paths, toolchain.Path)
		}
	}
	return paths
}
