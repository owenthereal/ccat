package rwvfs_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/kr/fs"
	"github.com/jingweno/ccat/Godeps/_workspace/src/golang.org/x/tools/godoc/vfs/mapfs"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/rwvfs"
)

func TestSub(t *testing.T) {
	m := rwvfs.Map(map[string]string{})
	sub := rwvfs.Sub(m, "/sub")

	err := sub.Mkdir("/")
	if err != nil {
		t.Fatal(err)
	}
	testIsDir(t, "sub", m, "/sub")

	f, err := sub.Create("f1")
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	testIsFile(t, "sub", m, "/sub/f1")

	f, err = sub.Create("/f2")
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	testIsFile(t, "sub", m, "/sub/f2")

	err = sub.Mkdir("/d1")
	if err != nil {
		t.Fatal(err)
	}
	testIsDir(t, "sub", m, "/sub/d1")

	err = sub.Mkdir("/d2")
	if err != nil {
		t.Fatal(err)
	}
	testIsDir(t, "sub", m, "/sub/d2")
}

func TestRWVFS(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "rwvfs-test-")
	if err != nil {
		t.Fatal("TempDir", err)
	}
	defer os.RemoveAll(tmpdir)

	h := http.Handler(rwvfs.HTTPHandler(rwvfs.Map(map[string]string{}), nil))
	httpServer := httptest.NewServer(h)
	defer httpServer.Close()
	httpURL, err := url.Parse(httpServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		fs   rwvfs.FileSystem
		path string
	}{
		{rwvfs.OS(tmpdir), "/foo"},
		{rwvfs.Map(map[string]string{}), "/foo"},
		{rwvfs.Sub(rwvfs.Map(map[string]string{}), "/x"), "/foo"},
		{rwvfs.HTTP(httpURL, nil), "/foo"},
	}
	for _, test := range tests {
		testWrite(t, test.fs, test.path)
		testMkdir(t, test.fs)
		testMkdirAll(t, test.fs)
		testGlob(t, test.fs)
	}
}

func testGlob(t *testing.T, fs rwvfs.FileSystem) {
	label := fmt.Sprintf("%T", fs)

	files := []string{"x/y/0.txt", "x/y/1.txt", "x/2.txt"}
	for _, file := range files {
		err := rwvfs.MkdirAll(fs, filepath.Dir(file))
		if err != nil {
			t.Fatalf("%s: MkdirAll: %s", label, err)
		}
		w, err := fs.Create(file)
		if err != nil {
			t.Errorf("%s: Create(%q): %s", label, file, err)
			return
		}
		w.Close()
	}

	globTests := []struct {
		prefix  string
		pattern string
		matches []string
	}{
		{"", "x/y/*.txt", []string{"x/y/0.txt", "x/y/1.txt"}},
		{"x/y", "x/y/*.txt", []string{"x/y/0.txt", "x/y/1.txt"}},
		{"", "x/*", []string{"x/y", "x/2.txt"}},
	}
	for _, test := range globTests {
		matches, err := rwvfs.Glob(rwvfs.Walkable(fs), test.prefix, test.pattern)
		if err != nil {
			t.Errorf("%s: Glob(prefix=%q, pattern=%q): %s", label, test.prefix, test.pattern, err)
			continue
		}
		sort.Strings(test.matches)
		sort.Strings(matches)
		if !reflect.DeepEqual(matches, test.matches) {
			t.Errorf("%s: Glob(prefix=%q, pattern=%q): got %v, want %v", label, test.prefix, test.pattern, matches, test.matches)
		}
	}
}

func testWrite(t *testing.T, fs rwvfs.FileSystem, path string) {
	label := fmt.Sprintf("%T", fs)

	w, err := fs.Create(path)
	if err != nil {
		t.Fatalf("%s: WriterOpen: %s", label, err)
	}

	input := []byte("qux")
	_, err = w.Write(input)
	if err != nil {
		t.Fatalf("%s: Write: %s", label, err)
	}

	err = w.Close()
	if err != nil {
		t.Fatalf("%s: w.Close: %s", label, err)
	}

	var r io.ReadCloser
	r, err = fs.Open(path)
	if err != nil {
		t.Fatalf("%s: Open: %s", label, err)
	}
	var output []byte
	output, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("%s: ReadAll: %s", label, err)
	}
	if !bytes.Equal(output, input) {
		t.Errorf("%s: got output %q, want %q", label, output, input)
	}

	r, err = fs.Open(path)
	if err != nil {
		t.Fatalf("%s: Open: %s", label, err)
	}
	output, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("%s: ReadAll: %s", label, err)
	}
	if !bytes.Equal(output, input) {
		t.Errorf("%s: got output %q, want %q", label, output, input)
	}

	err = fs.Remove(path)
	if err != nil {
		t.Errorf("%s: Remove(%q): %s", label, path, err)
	}
	testPathDoesNotExist(t, label, fs, path)
}

func testMkdir(t *testing.T, fs rwvfs.FileSystem) {
	label := fmt.Sprintf("%T", fs)

	if strings.Contains(label, "subFS") {
		if err := fs.Mkdir("/"); err != nil && !os.IsExist(err) {
			t.Fatalf("%s: subFS Mkdir(/): %s", label, err)
		}
	}
	if strings.Contains(label, "mapFS") {
		if err := fs.Mkdir("/"); err != nil && !os.IsExist(err) {
			t.Fatalf("%s: mapFS Mkdir(/): %s", label, err)
		}
	}

	fi, err := fs.Stat(".")
	if err != nil {
		t.Fatalf("%s: Stat(.): %s", label, err)
	}
	if !fi.Mode().IsDir() {
		t.Fatalf("%s: got Stat(.) FileMode %o, want IsDir", label, fi.Mode())
	}

	fi, err = fs.Stat("/")
	if err != nil {
		t.Fatalf("%s: Stat(/): %s", label, err)
	}
	if !fi.Mode().IsDir() {
		t.Fatalf("%s: got Stat(/) FileMode %o, want IsDir", label, fi.Mode())
	}

	if _, err := fs.ReadDir("."); err != nil {
		t.Fatalf("%s: ReadDir(.): %s", label, err)
	}
	if _, err := fs.ReadDir("/"); err != nil {
		t.Fatalf("%s: ReadDir(/): %s", label, err)
	}

	fis, err := fs.ReadDir("/")
	if err != nil {
		t.Fatalf("%s: ReadDir(/): %s", label, err)
	}
	if len(fis) != 0 {
		t.Fatalf("%s: ReadDir(/): got %d file infos (%v), want none (is it including .?)", label, len(fis), fis)
	}

	err = fs.Mkdir("dir0")
	if err != nil {
		t.Fatalf("%s: Mkdir(dir0): %s", label, err)
	}
	testIsDir(t, label, fs, "dir0")
	testIsDir(t, label, fs, "/dir0")

	err = fs.Mkdir("/dir1")
	if err != nil {
		t.Fatalf("%s: Mkdir(/dir1): %s", label, err)
	}
	testIsDir(t, label, fs, "dir1")
	testIsDir(t, label, fs, "/dir1")

	err = fs.Mkdir("/dir1")
	if !os.IsExist(err) {
		t.Errorf("%s: Mkdir(/dir1) again: got err %v, want os.IsExist-satisfying error", label, err)
	}

	err = fs.Mkdir("/parent-doesnt-exist/dir2")
	if !os.IsNotExist(err) {
		t.Errorf("%s: Mkdir(/parent-doesnt-exist/dir2): got error %v, want os.IsNotExist-satisfying error", label, err)
	}

	err = fs.Remove("/dir1")
	if err != nil {
		t.Errorf("%s: Remove(/dir1): %s", label, err)
	}
	testPathDoesNotExist(t, label, fs, "/dir1")
}

func testMkdirAll(t *testing.T, fs rwvfs.FileSystem) {
	label := fmt.Sprintf("%T", fs)

	err := rwvfs.MkdirAll(fs, "/a/b/c")
	if err != nil {
		t.Fatalf("%s: MkdirAll: %s", label, err)
	}
	testIsDir(t, label, fs, "/a")
	testIsDir(t, label, fs, "/a/b")
	testIsDir(t, label, fs, "/a/b/c")

	err = rwvfs.MkdirAll(fs, "/a/b/c")
	if err != nil {
		t.Fatalf("%s: MkdirAll again: %s", label, err)
	}
}

func testIsDir(t *testing.T, label string, fs rwvfs.FileSystem, path string) {
	fi, err := fs.Stat(path)
	if err != nil {
		t.Fatalf("%s: Stat(%q): %s", label, path, err)
	}

	if fi == nil {
		t.Fatalf("%s: FileInfo (%q) == nil", label, path)
	}

	if !fi.IsDir() {
		t.Errorf("%s: got fs.Stat(%q) IsDir() == false, want true", label, path)
	}
}

func testIsFile(t *testing.T, label string, fs rwvfs.FileSystem, path string) {
	fi, err := fs.Stat(path)
	if err != nil {
		t.Fatalf("%s: Stat(%q): %s", label, path, err)
	}

	if !fi.Mode().IsRegular() {
		t.Errorf("%s: got fs.Stat(%q) Mode().IsRegular() == false, want true", label, path)
	}
}

func testPathDoesNotExist(t *testing.T, label string, fs rwvfs.FileSystem, path string) {
	fi, err := fs.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		t.Errorf("%s: Stat(%q): want os.IsNotExist-satisfying error, got %q", label, path, err)
	} else if err == nil {
		t.Errorf("%s: Stat(%q): want file to not exist, got existing file with FileInfo %+v", label, path, fi)
	}
}

func TestMap_MkdirAllWithRootNotExists(t *testing.T) {
	m := map[string]string{}
	fs := rwvfs.Sub(rwvfs.Map(m), "x")

	paths := []string{"a/b", "/c/d"}
	for _, path := range paths {
		if err := rwvfs.MkdirAll(fs, path); err != nil {
			t.Errorf("MkdirAll %q: %s", path, err)
		}
	}
}

func TestHTTP_BaseURL(t *testing.T) {
	m := map[string]string{"b/c": "c"}
	mapFS := rwvfs.Map(m)

	prefix := "/foo/bar/baz"

	h := http.Handler(http.StripPrefix(prefix, rwvfs.HTTPHandler(mapFS, nil)))
	httpServer := httptest.NewServer(h)
	defer httpServer.Close()
	httpURL, err := url.Parse(httpServer.URL + prefix)
	if err != nil {
		t.Fatal(err)
	}

	fs := rwvfs.HTTP(httpURL, nil)

	if err := rwvfs.MkdirAll(fs, "b"); err != nil {
		t.Errorf("MkdirAll %q: %s", "b", err)
	}

	fis, err := fs.ReadDir("b")
	if err != nil {
		t.Fatal(err)
	}
	if len(fis) != 1 {
		t.Errorf("got len(fis) == %d, want 1", len(fis))
	}
	if wantName := "c"; fis[0].Name() != wantName {
		t.Errorf("got name == %q, want %q", fis[0].Name(), wantName)
	}
}

func TestMap_Walk(t *testing.T) {
	m := map[string]string{"a": "a", "b/c": "c", "b/x/y/z": "z"}
	mapFS := rwvfs.Map(m)

	var names []string
	w := fs.WalkFS(".", rwvfs.Walkable(mapFS))
	for w.Step() {
		if err := w.Err(); err != nil {
			t.Fatalf("walk path %q: %s", w.Path(), err)
		}
		names = append(names, w.Path())
	}

	wantNames := []string{".", "a", "b", "b/c", "b/x", "b/x/y", "b/x/y/z"}
	sort.Strings(names)
	sort.Strings(wantNames)
	if !reflect.DeepEqual(names, wantNames) {
		t.Errorf("got entry names %v, want %v", names, wantNames)
	}
}

func TestMap_Walk2(t *testing.T) {
	m := map[string]string{"a/b/c/d": "a"}
	mapFS := rwvfs.Map(m)

	var names []string
	w := fs.WalkFS(".", rwvfs.Walkable(rwvfs.Sub(mapFS, "a/b")))
	for w.Step() {
		if err := w.Err(); err != nil {
			t.Fatalf("walk path %q: %s", w.Path(), err)
		}
		names = append(names, w.Path())
	}

	wantNames := []string{".", "c", "c/d"}
	sort.Strings(names)
	sort.Strings(wantNames)
	if !reflect.DeepEqual(names, wantNames) {
		t.Errorf("got entry names %v, want %v", names, wantNames)
	}
}

func TestReadOnly(t *testing.T) {
	m := map[string]string{"x": "y"}
	rfs := mapfs.New(m)
	wfs := rwvfs.ReadOnly(rfs)

	if _, err := rfs.Stat("/x"); err != nil {
		t.Error(err)
	}

	_, err := wfs.Create("/y")
	if want := (&os.PathError{"create", "/y", rwvfs.ErrReadOnly}); !reflect.DeepEqual(err, want) {
		t.Errorf("Create: got err %v, want %v", err, want)
	}

	err = wfs.Mkdir("/y")
	if want := (&os.PathError{"mkdir", "/y", rwvfs.ErrReadOnly}); !reflect.DeepEqual(err, want) {
		t.Errorf("Mkdir: got err %v, want %v", err, want)
	}

	err = wfs.Remove("/y")
	if want := (&os.PathError{"remove", "/y", rwvfs.ErrReadOnly}); !reflect.DeepEqual(err, want) {
		t.Errorf("Remove: got err %v, want %v", err, want)
	}

}

func TestOS_Symlink(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "rwvfs-test-")
	if err != nil {
		t.Fatal("TempDir", err)
	}
	defer os.RemoveAll(tmpdir)
	want := "hello"

	if err := ioutil.WriteFile(filepath.Join(tmpdir, "myfile"), []byte(want), 0600); err != nil {
		t.Fatal(err)
	}

	osfs := rwvfs.OS(tmpdir)
	if err := osfs.(rwvfs.LinkFS).Symlink("myfile", "mylink"); err != nil {
		t.Fatal(err)
	}
	got, err := ioutil.ReadFile(filepath.Join(tmpdir, "mylink"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("%s: ReadLink: got %q, want %q", osfs, string(got), want)
	}
}

func TestOS_Symlink_walkable(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "rwvfs-test-")
	if err != nil {
		t.Fatal("TempDir", err)
	}
	defer os.RemoveAll(tmpdir)
	want := "hello"

	if err := ioutil.WriteFile(filepath.Join(tmpdir, "myfile"), []byte(want), 0600); err != nil {
		t.Fatal(err)
	}

	osfs := rwvfs.OS(tmpdir)
	if err := rwvfs.Walkable(osfs).(rwvfs.LinkFS).Symlink("myfile", "mylink"); err != nil {
		t.Fatal(err)
	}
	got, err := ioutil.ReadFile(filepath.Join(tmpdir, "mylink"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("%s: ReadLink: got %q, want %q", osfs, string(got), want)
	}
}

func TestSub_Symlink(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "rwvfs-test-")
	if err != nil {
		t.Fatal("TempDir", err)
	}
	//defer os.RemoveAll(tmpdir)
	want := "hello"

	if err := os.Mkdir(filepath.Join(tmpdir, "mydir"), 0700); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(filepath.Join(tmpdir, "mydir", "myfile"), []byte(want), 0600); err != nil {
		t.Fatal(err)
	}

	osfs := rwvfs.OS(tmpdir)
	sub := rwvfs.Sub(osfs, "mydir")
	if err := sub.(rwvfs.LinkFS).Symlink("myfile", "mylink"); err != nil {
		t.Fatal(err, osfs)
	}
	got, err := ioutil.ReadFile(filepath.Join(tmpdir, "mydir", "mylink"))
	if err != nil {
		t.Fatal(err, osfs, sub)
	}
	if string(got) != want {
		t.Errorf("%s: ReadLink: got %q, want %q", osfs, string(got), want)
	}
}

func TestOS_ReadLink(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "rwvfs-test-")
	if err != nil {
		t.Fatal("TempDir", err)
	}
	defer os.RemoveAll(tmpdir)

	if err := ioutil.WriteFile(filepath.Join(tmpdir, "myfile"), []byte("hello"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(tmpdir, "myfile"), filepath.Join(tmpdir, "mylink")); err != nil {
		t.Fatal(err)
	}

	osfs := rwvfs.OS(tmpdir)
	dst, err := osfs.(rwvfs.LinkFS).ReadLink("mylink")
	if err != nil {
		t.Fatal(err)
	}
	if want := "myfile"; dst != want {
		t.Errorf("%s: ReadLink: got %q, want %q", osfs, dst, want)
	}
}

func TestOS_ReadLink_walkable(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "rwvfs-test-")
	if err != nil {
		t.Fatal("TempDir", err)
	}
	defer os.RemoveAll(tmpdir)

	if err := ioutil.WriteFile(filepath.Join(tmpdir, "myfile"), []byte("hello"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(tmpdir, "myfile"), filepath.Join(tmpdir, "mylink")); err != nil {
		t.Fatal(err)
	}

	osfs := rwvfs.OS(tmpdir)
	dst, err := rwvfs.Walkable(osfs).(rwvfs.LinkFS).ReadLink("mylink")
	if err != nil {
		t.Fatal(err)
	}
	if want := "myfile"; dst != want {
		t.Errorf("%s: ReadLink: got %q, want %q", osfs, dst, want)
	}
}

func TestSub_ReadLink(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "rwvfs-test-")
	if err != nil {
		t.Fatal("TempDir", err)
	}
	defer os.RemoveAll(tmpdir)

	if err := os.Mkdir(filepath.Join(tmpdir, "mydir"), 0700); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(filepath.Join(tmpdir, "mydir", "myfile"), []byte("hello"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(tmpdir, "mydir", "myfile"), filepath.Join(tmpdir, "mydir", "mylink")); err != nil {
		t.Fatal(err)
	}

	osfs := rwvfs.OS(tmpdir)
	sub := rwvfs.Sub(osfs, "mydir")
	dst, err := sub.(rwvfs.LinkFS).ReadLink("mylink")
	if err != nil {
		t.Fatal(err)
	}
	if want := "myfile"; dst != want {
		t.Errorf("%s: ReadLink: got %q, want %q", osfs, dst, want)
	}
}

func TestOS_ReadLink_ErrOutsideRoot(t *testing.T) {
	tmpdir1, err := ioutil.TempDir("", "rwvfs-test-")
	if err != nil {
		t.Fatal("TempDir", err)
	}
	defer os.RemoveAll(tmpdir1)

	tmpdir2, err := ioutil.TempDir("", "rwvfs-test-")
	if err != nil {
		t.Fatal("TempDir", err)
	}
	defer os.RemoveAll(tmpdir2)

	if err := ioutil.WriteFile(filepath.Join(tmpdir1, "myfile"), []byte("hello"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(tmpdir1, "myfile"), filepath.Join(tmpdir2, "mylink")); err != nil {
		t.Fatal(err)
	}

	osfs := rwvfs.OS(tmpdir2)
	dst, err := osfs.(rwvfs.LinkFS).ReadLink("mylink")
	if want := rwvfs.ErrOutsideRoot; err != want {
		t.Fatalf("%s: ReadLink: got err %v, want %v", osfs, err, want)
	}
	if want := filepath.Join(tmpdir1, "myfile"); dst != want {
		t.Errorf("%s: ReadLink: got %q, want %q", osfs, dst, want)
	}
}
