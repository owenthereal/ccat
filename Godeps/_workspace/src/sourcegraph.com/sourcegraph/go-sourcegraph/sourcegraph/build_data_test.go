package sourcegraph

import (
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"testing"

	"sort"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/rwvfs"
)

func TestBuildDataService_GetBuildDataFile(t *testing.T) {
	setup()
	defer teardown()

	want := []byte("hello")

	var called int
	mux.HandleFunc(urlPath(t, router.RepoBuildDataEntry, map[string]string{"RepoSpec": "r.com/x", "Rev": "c", "Path": "a/b"}), func(w http.ResponseWriter, r *http.Request) {
		called++

		switch r.Method {
		case "GET":
			w.Write(want)
		case "HEAD":
		}
	})

	file, _, err := GetBuildDataFile(client.BuildData, BuildDataFileSpec{RepoRev: RepoRevSpec{RepoSpec: RepoSpec{URI: "r.com/x"}, Rev: "c"}, Path: "a/b"})
	if err != nil {
		t.Fatalf("GetBuildDataFile returned error: %v", err)
	}
	defer file.Close()

	if called != 2 {
		t.Fatalf("got called == %d, want 2", called)
	}

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(fileData, want) {
		t.Errorf("GetBuildDataFile returned file data %+v, want %+v", fileData, want)
	}
}

func TestBuildDataService_ListAll(t *testing.T) {
	setup()
	defer teardown()

	pathPrefix := urlPath(t, router.RepoBuildDataEntry, map[string]string{"RepoSpec": "r.com/x", "Rev": "c", "Path": "."})
	fs := rwvfs.Map(map[string]string{
		"a":     "a",
		"b/c":   "c",
		"b/d/e": "e",
	})
	mux.Handle(pathPrefix+"/", http.StripPrefix(pathPrefix, rwvfs.HTTPHandler(fs, nil)))

	fs, err := client.BuildData.FileSystem(RepoRevSpec{RepoSpec: RepoSpec{URI: "r.com/x"}, Rev: "c"})
	if err != nil {
		t.Fatal(err)
	}

	entries, err := rwvfs.StatAllRecursive(".", rwvfs.Walkable(fs))
	if err != nil {
		t.Fatalf("StatAllRecursive returned error: %v", err)
	}

	names := fileInfoNames(entries)
	wantNames := []string{".", "a", "b", "b/c", "b/d", "b/d/e"}
	sort.Strings(names)
	sort.Strings(wantNames)
	if !reflect.DeepEqual(names, wantNames) {
		t.Errorf("got entry names %v, want %v", names, wantNames)
	}
}

func fileInfoNames(fis []os.FileInfo) []string {
	names := make([]string, len(fis))
	for i, fi := range fis {
		names[i] = fi.Name()
	}
	return names
}
