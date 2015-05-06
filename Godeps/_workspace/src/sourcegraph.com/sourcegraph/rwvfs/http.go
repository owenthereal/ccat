package rwvfs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	pathpkg "path"
	"strconv"
	"strings"
	"syscall"

	"github.com/jingweno/ccat/Godeps/_workspace/src/golang.org/x/tools/godoc/vfs"
)

// HTTP creates a new VFS that accesses paths on an HTTP server.
func HTTP(base *url.URL, httpClient *http.Client) FileSystem {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &httpFS{
		baseURL:    base,
		httpClient: httpClient,
	}
}

type httpFS struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func (c *httpFS) String() string { return "http(" + c.baseURL.String() + ")" }

func (c *httpFS) Lstat(path string) (os.FileInfo, error) {
	fi, err := c.stat(c.httpClient, path)
	if err != nil {
		err = &os.PathError{"lstat", path, err}
	}
	return fi, err
}

func (c *httpFS) Stat(path string) (os.FileInfo, error) {
	fi, err := c.stat(c.httpClient, path)
	if err != nil {
		err = &os.PathError{"stat", path, err}
	}
	return fi, err
}

func (c *httpFS) stat(httpClient *http.Client, path string) (os.FileInfo, error) {
	fi := fileInfo{name: pathpkg.Base(path)}
	req, err := c.newRequest("HEAD", path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.send(httpClient, req)
	if resp != nil {
		defer resp.Body.Close()
		// Don't check for errors, so that this VFS can be used
		// against HTTP endpoints other than just those created by
		// HTTPHandler.
		setHTTPResponseFileInfo(resp, &fi)
	}
	if err != nil && !isIgnoredRedirectErr(err) {
		return nil, err
	}
	return fi, nil
}

func setHTTPResponseFileInfo(resp *http.Response, fi *fileInfo) error {
	if lastMod := resp.Header.Get("last-modified"); lastMod != "" {
		var err error
		fi.modTime, err = http.ParseTime(lastMod)
		if err != nil {
			return err
		}
	}
	if resp.StatusCode >= 300 && resp.StatusCode <= 399 {
		fi.symlink = true
	}
	switch resp.Header.Get("content-type") {
	case httpFileContentType: // default, nothing to do
	case httpDirContentType:
		fi.dir = true
	case httpSymlinkContentType:
		fi.symlink = true
	}
	fi.size = resp.ContentLength
	return nil
}

const (
	// MIME types for the HTTP response Content-Type header to
	// indicate which type of resource exists at a path.
	httpFileContentType    = "application/vnd.rwvfs.file"
	httpDirContentType     = "application/vnd.rwvfs.dir"
	httpSymlinkContentType = "application/vnd.rwvfs.symlink"
)

func (c *httpFS) ReadDir(path string) ([]os.FileInfo, error) {
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", httpDirContentType)

	resp, err := c.send(nil, req)
	if err != nil {
		return nil, &os.PathError{"readdir", path, err}
	}
	defer resp.Body.Close()

	if contentType := resp.Header.Get("content-type"); contentType != httpDirContentType {
		return nil, &os.PathError{"readdir", path, syscall.ENOTDIR}
	}

	var fis []*fileInfoJSON
	if err := json.NewDecoder(resp.Body).Decode(&fis); err != nil {
		return nil, err
	}
	fis2 := make([]os.FileInfo, len(fis))
	for i, fi := range fis {
		fis2[i] = fi
	}
	return fis2, nil
}

func (c *httpFS) Open(name string) (vfs.ReadSeekCloser, error) {
	req, err := c.newRequest("GET", name, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", httpFileContentType)

	resp, err := c.send(nil, req)
	if err != nil {
		return nil, &os.PathError{"open", name, err}
	}
	return failSeeker{resp.Body}, nil
}

type failSeeker struct{ io.ReadCloser }

func (failSeeker) Seek(offset int64, whence int) (int64, error) {
	// TODO(sqs): is Seek used by any clients of rwvfs? if so,
	// consider buffering the HTTP response so it can actually
	// implement Seek.
	return 0, errors.New("rwvfs.HTTP VFS does not support seeking")
}

func (c *httpFS) Create(path string) (io.WriteCloser, error) {
	return &httpFilePost{c: c, path: path}, nil
}

type httpFilePost struct {
	bytes.Buffer
	c    *httpFS
	path string

	closed bool
}

func (f *httpFilePost) Close() error {
	if f.closed {
		return nil
	}
	f.closed = true
	req, err := f.c.newRequest("PUT", f.path, ioutil.NopCloser(&f.Buffer))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", httpFileContentType)
	req.ContentLength = int64(f.Buffer.Len())
	resp, err := f.c.send(nil, req)
	if err != nil {
		return &os.PathError{"create", f.path, err}
	}
	return resp.Body.Close()
}

func (c *httpFS) Mkdir(name string) error {
	req, err := c.newRequest("PUT", name, nil)
	if err != nil {
		return err
	}
	req.Header.Set("content-type", httpDirContentType)

	resp, err := c.send(nil, req)
	if err != nil {
		return &os.PathError{"mkdir", name, err}
	}
	return resp.Body.Close()
}

func (c *httpFS) Remove(name string) error {
	req, err := c.newRequest("DELETE", name, nil)
	if err != nil {
		return err
	}
	resp, err := c.send(nil, req)
	if err != nil {
		return &os.PathError{"remove", name, err}
	}
	return resp.Body.Close()
}

// newRequest creates a new (unsent) HTTP request.
func (c *httpFS) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	urlPath := pathpkg.Join(c.baseURL.Path, strings.TrimPrefix(path, "/"))
	url := c.baseURL.ResolveReference(&url.URL{Path: urlPath})
	return http.NewRequest(method, url.String(), body)
}

// send issues a request using the provided HTTP client (or
// c.httpClient if nil). If the response has a non-20x (200-299) HTTP
// status code, a non-nil error is returned. Callers are responsible
// for closing the HTTP response body unless a non-nil error is
// returned.
func (c *httpFS) send(httpClient *http.Client, req *http.Request) (*http.Response, error) {
	if httpClient == nil {
		httpClient = c.httpClient
	}
	resp, err := httpClient.Do(req)
	isHTTP20x := resp != nil && (resp.StatusCode >= 200 && resp.StatusCode <= 299)
	if resp != nil && (err != nil || !isHTTP20x) {
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	}
	if err != nil {
		return resp, err
	}
	if resp != nil && !isHTTP20x {
		resp.Body.Close()
		switch resp.StatusCode {
		case http.StatusNotFound:
			err = os.ErrNotExist
		case http.StatusConflict:
			err = os.ErrExist
		default:
			err = fmt.Errorf("http status %d", resp.StatusCode)
		}
	}
	return resp, err
}

// HTTPHandler creates an http.Handler that allows HTTP clients to
// access fs. It should be accessed by clients created using this
// package's HTTP func.
func HTTPHandler(fs FileSystem, logTo io.Writer) http.Handler {
	if logTo == nil {
		logTo = ioutil.Discard
	}
	return &httpFSHandler{fs, log.New(logTo, "rwvfs HTTP: ", log.Flags())}
}

type httpFSHandler struct {
	fs  FileSystem
	log *log.Logger
}

func (h *httpFSHandler) ServeHTTPAndReturnError(w http.ResponseWriter, r *http.Request) error {
	var err error
	switch r.Method {
	case "GET":
		err = h.get(w, r)
	case "HEAD":
		err = h.stat(w, r)
	case "PUT":
		err = h.put(w, r)
	case "DELETE":
		err = h.remove(w, r)
	}
	return err
}

func (h *httpFSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.ServeHTTPAndReturnError(w, r)
	var status int
	if os.IsNotExist(err) {
		status = http.StatusNotFound
	} else if os.IsExist(err) {
		status = http.StatusConflict
	} else {
		status = http.StatusInternalServerError
	}
	if err != nil {
		h.log.Printf("rwvfs %s %s: %s", r.Method, r.URL, err)
		http.Error(w, "", status)
	}
}

func (h *httpFSHandler) get(w http.ResponseWriter, r *http.Request) error {
	switch r.Header.Get("accept") {
	case httpFileContentType:
		return h.open(w, r)
	case httpDirContentType:
		return h.readDir(w, r)
	}

	fi, err := h.fs.Stat(r.URL.Path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return h.readDir(w, r)
	}
	return h.open(w, r)
}

func (h *httpFSHandler) open(w http.ResponseWriter, r *http.Request) error {
	fi, err := h.fs.Stat(r.URL.Path)
	if err != nil {
		return err
	}

	notMod := false
	if ifModSince, err := http.ParseTime(r.Header.Get("if-modified-since")); err == nil {
		if !fi.ModTime().IsZero() && fi.ModTime().Before(ifModSince) {
			w.WriteHeader(http.StatusNotModified)
			notMod = true
		}
	}

	writeFileInfoHeaders(w, fi, true)

	if notMod {
		return nil
	}

	f, err := h.fs.Open(r.URL.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(w, f)
	return err
}

func (h *httpFSHandler) readDir(w http.ResponseWriter, r *http.Request) error {
	fis, err := h.fs.ReadDir(r.URL.Path)
	if err != nil {
		return err
	}

	fi, err := h.fs.Stat(r.URL.Path)
	if err != nil {
		return err
	}
	writeFileInfoHeaders(w, fi, false)
	w.WriteHeader(http.StatusOK)

	jsonFIs := make([]*fileInfoJSON, len(fis))
	for i, fi := range fis {
		jsonFIs[i] = &fileInfoJSON{fi}
	}

	return json.NewEncoder(w).Encode(jsonFIs)
}

func (h *httpFSHandler) stat(w http.ResponseWriter, r *http.Request) error {
	fi, err := h.fs.Stat(r.URL.Path)
	if err != nil {
		return err
	}
	writeFileInfoHeaders(w, fi, true)
	w.WriteHeader(http.StatusOK)
	return nil
}

func writeFileInfoHeaders(w http.ResponseWriter, fi os.FileInfo, writeContentLength bool) {
	if writeContentLength {
		w.Header().Set("content-length", strconv.FormatInt(fi.Size(), 10))
	}
	if !fi.ModTime().IsZero() {
		w.Header().Set("last-modified", fi.ModTime().Format(http.TimeFormat))
	}
	if fi.IsDir() {
		w.Header().Set("content-type", httpDirContentType)
	} else if fi.Mode()&os.ModeSymlink > 0 {
		// TODO(sqs): get link dest (requires adding Readlink to VFS
		// interface?
		w.Header().Set("content-type", httpSymlinkContentType)
	} else {
		w.Header().Set("content-type", httpFileContentType)
	}
}

func (h *httpFSHandler) put(w http.ResponseWriter, r *http.Request) error {
	switch r.Header.Get("content-type") {
	case httpFileContentType:
		return h.create(w, r)
	case httpDirContentType:
		return h.mkdir(w, r)
	}
	http.Error(w, "", http.StatusBadRequest)
	return nil
}

func (h *httpFSHandler) create(w http.ResponseWriter, r *http.Request) error {
	if r.ContentLength == -1 {
		http.Error(w, "", http.StatusLengthRequired)
		return nil
	}
	f, err := h.fs.Create(slash(r.URL.Path))
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(f, r.Body); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *httpFSHandler) mkdir(w http.ResponseWriter, r *http.Request) error {
	return h.fs.Mkdir(slash(r.URL.Path))
}

func (h *httpFSHandler) remove(w http.ResponseWriter, r *http.Request) error {
	return h.fs.Remove(r.URL.Path)
}

// ignoreRedirectsHTTPClient returns an HTTP client identical to c
// (using the same Transport, etc.)  except that when it encounters a
// redirect, it returns errIgnoredRedirect.
func ignoreRedirectsHTTPClient(c *http.Client) *http.Client {
	c2 := *c
	c2.CheckRedirect = func(r *http.Request, via []*http.Request) error { return errIgnoredRedirect }
	return &c2
}

var errIgnoredRedirect = errors.New("not following redirect")

func isIgnoredRedirectErr(err error) bool {
	if err, ok := err.(*url.Error); ok && err.Err == errIgnoredRedirect {
		return true
	}
	return false
}

func asJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
