package src

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/jingweno/ccat/Godeps/_workspace/src/code.google.com/p/rog-go/parallel"
	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/kr/fs"
	"github.com/jingweno/ccat/Godeps/_workspace/src/golang.org/x/tools/godoc/vfs"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/rwvfs"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/buildstore"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/store"
)

func init() {
	buildDataGroup, err := CLI.AddCommand("build-data",
		"build data operations",
		"The build-data command group contains subcommands for listing, displaying, uploading, and downloading build data.",
		&buildDataCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
	buildDataGroup.Aliases = []string{"bd"}

	c, err := buildDataGroup.AddCommand("ls",
		"list build data files and dirs",
		"The ls command lists build data files and directories for a repository at a specific commit.",
		&buildDataListCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
	setDefaultCommitIDOpt(c)

	c, err = buildDataGroup.AddCommand("cat",
		"display contents of build files",
		"The command displays the contents of a build data file for a repository at a specific commit.",
		&buildDataCatCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
	setDefaultCommitIDOpt(c)

	c, err = buildDataGroup.AddCommand("rm",
		"remove build data files and dirs",
		"The rm command removes a build data file or directory for a repository at a specific commit.",
		&buildDataRemoveCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
	setDefaultCommitIDOpt(c)

	c, err = buildDataGroup.AddCommand("fetch",
		"fetch remote build data to local dir",
		"The fetch command fetches remote build data for the current repository to the local .srclib-cache directory.",
		&buildDataFetchCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
	setDefaultCommitIDOpt(c)

	c, err = buildDataGroup.AddCommand("upload",
		"upload local build data to remote",
		"The upload command uploads local build data (in .srclib-cache) for the current repository to the remote.",
		&buildDataUploadCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
	setDefaultCommitIDOpt(c)
}

type BuildDataCmd struct {
}

var buildDataCmd BuildDataCmd

func (c *BuildDataCmd) Execute(args []string) error { return nil }

type buildDataSingleCommitCommonOpts struct {
	CommitID string `long:"commit" description:"commit ID of build data to operate on" required:"yes"`
}

func (c buildDataSingleCommitCommonOpts) getLocalFileSystem() (rwvfs.FileSystem, string, error) {
	return getLocalBuildDataFS(c.CommitID)
}

func (c buildDataSingleCommitCommonOpts) getRemoteFileSystem() (rwvfs.FileSystem, string, sourcegraph.RepoRevSpec, error) {
	cl := NewAPIClientWithAuthIfPresent()
	rrepo, err := getRemoteRepo(cl)
	if err != nil {
		return nil, "", sourcegraph.RepoRevSpec{}, err
	}

	repoRevSpec := sourcegraph.RepoRevSpec{RepoSpec: rrepo.RepoSpec(), Rev: c.CommitID, CommitID: c.CommitID}
	fs, err := cl.BuildData.FileSystem(repoRevSpec)
	return fs, fmt.Sprintf("remote repository (URI %s, commit %s)", rrepo.URI, c.CommitID), repoRevSpec, err
}

type buildDataSingleRepoCommonOpts struct {
	buildDataSingleCommitCommonOpts
	Local bool `long:"local" description:"execute build-data ls/cat/rm subcommands on local build data (.srclib-cache), not remote (sourcegraph.com)"`
}

// getFileSystem gets either the local or remote build data FS,
// depending on the value of c.Local.
func (c *buildDataSingleRepoCommonOpts) getFileSystem() (rwvfs.FileSystem, string, error) {
	lrepo, err := openLocalRepo()
	if err != nil {
		return nil, "", err
	}
	return getBuildDataFS(c.Local, lrepo.URI(), c.CommitID)
}

func getLocalBuildDataFS(commitID string) (rwvfs.FileSystem, string, error) {
	lrepo, err := openLocalRepo()
	if lrepo == nil || lrepo.RootDir == "" || commitID == "" {
		return nil, "", err
	}
	localStore, err := buildstore.LocalRepo(lrepo.RootDir)
	if err != nil {
		return nil, "", err
	}
	return localStore.Commit(commitID), fmt.Sprintf("local repository (root dir %s, commit %s)", lrepo.RootDir, commitID), nil
}

// getRemoteBuildDataFS gets the remote build data file system for
// repo at commitID. It returns an error if repo is empty.
func getRemoteBuildDataFS(repo, commitID string) (rwvfs.FileSystem, string, sourcegraph.RepoRevSpec, error) {
	if repo == "" {
		err := errors.New("getRemoteBuildDataFS: repo cannot be empty")
		return nil, "", sourcegraph.RepoRevSpec{}, err
	}
	cl := NewAPIClientWithAuthIfPresent()
	rrepo, _, err := cl.Repos.Get(sourcegraph.RepoSpec{URI: repo}, nil)
	if err != nil {
		return nil, "", sourcegraph.RepoRevSpec{}, err
	}

	repoRevSpec := sourcegraph.RepoRevSpec{RepoSpec: rrepo.RepoSpec(), Rev: commitID, CommitID: commitID}
	fs, err := cl.BuildData.FileSystem(repoRevSpec)
	return fs, fmt.Sprintf("remote repository (URI %s, commit %s)", rrepo.URI, commitID), repoRevSpec, err
}

// getBuildDataFS gets the build data file system for repo at
// commitID. If local is true, repo is ignored and build data is
// fetched for the local repo.
func getBuildDataFS(local bool, repo, commitID string) (rwvfs.FileSystem, string, error) {
	if local {
		return getLocalBuildDataFS(commitID)
	}
	fs, label, _, err := getRemoteBuildDataFS(repo, commitID)
	return fs, label, err
}

type BuildDataListCmd struct {
	buildDataSingleRepoCommonOpts

	Args struct {
		Dir string `name:"DIR" default:"." description:"list build data files in this dir"`
	} `positional-args:"yes"`

	Recursive bool   `short:"r" long:"recursive" description:"list recursively"`
	Long      bool   `short:"l" long:"long" description:"show file sizes and times"`
	Type      string `long:"type" description:"show only entries of this type (f=file, d=dir)"`
	URLs      bool   `long:"urls" description:"show URLs to build data files (implies -l)"`
}

var buildDataListCmd BuildDataListCmd

func (c *BuildDataListCmd) Execute(args []string) error {
	if c.URLs && c.Local {
		return fmt.Errorf("using --urls is incompatible with the build-data -l/--local option because local build data files do not have a URL")
	}
	if c.URLs {
		c.Long = true
	}
	dir := c.Args.Dir
	if dir == "" {
		dir = "."
	}

	bdfs, repoLabel, err := c.getFileSystem()
	if err != nil {
		return err
	}

	if GlobalOpt.Verbose {
		log.Printf("Listing build files for %s in dir %q", repoLabel, dir)
	}

	// Only used for constructing the URLs for remote build data.
	var repoRevSpec sourcegraph.RepoRevSpec
	if !c.Local {
		cl := NewAPIClientWithAuthIfPresent()
		rrepo, err := getRemoteRepo(cl)
		if err != nil {
			return err
		}
		repoRevSpec.RepoSpec = rrepo.RepoSpec()

		lrepo, err := openLocalRepo()
		if err != nil {
			return err
		}
		repoRevSpec.Rev = lrepo.CommitID
		repoRevSpec.CommitID = lrepo.CommitID
	}

	printFile := func(fi os.FileInfo) {
		if c.Type == "f" && !fi.Mode().IsRegular() {
			return
		}
		if c.Type == "d" && !fi.Mode().IsDir() {
			return
		}

		var suffix string
		if fi.IsDir() {
			suffix = "/"
		}

		var urlStr string
		if c.URLs {
			spec := sourcegraph.BuildDataFileSpec{RepoRev: repoRevSpec, Path: filepath.Join(dir, fi.Name())}

			// TODO(sqs): use sourcegraph.Router when it is merged to go-sourcegraph master
			u, err := router.NewAPIRouter(nil).Get(router.RepoBuildDataEntry).URLPath(router.MapToArray(spec.RouteVars())...)
			if err != nil {
				log.Fatal(err)
			}

			// Strip leading "/" so that the URL is relative to the
			// endpoint URL even if the endpoint URL contains a path.
			urlStr = getEndpointURL().ResolveReference(&url.URL{Path: u.Path[1:]}).String()
		}

		if c.Long {
			var timeStr string
			if !fi.ModTime().IsZero() {
				timeStr = fi.ModTime().Format("Jan _2 15:04")
			}
			fmt.Printf("% 7d %12s %s%s %s\n", fi.Size(), timeStr, fi.Name(), suffix, urlStr)
		} else {
			fmt.Println(fi.Name() + suffix)
		}
	}

	var fis []os.FileInfo
	if c.Recursive {
		w := fs.WalkFS(dir, rwvfs.Walkable(bdfs))
		for w.Step() {
			if err := w.Err(); err != nil {
				return err
			}
			printFile(treeFileInfo{w.Path(), w.Stat()})
		}
	} else {
		fis, err = bdfs.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, fi := range fis {
			printFile(fi)
		}
	}

	return nil
}

type treeFileInfo struct {
	path string
	os.FileInfo
}

func (fi treeFileInfo) Name() string { return fi.path }

type BuildDataCatCmd struct {
	buildDataSingleRepoCommonOpts

	Args struct {
		File string `name:"FILE" default:"." description:"file whose contents to print"`
	} `positional-args:"yes"`
}

var buildDataCatCmd BuildDataCatCmd

func (c *BuildDataCatCmd) Execute(args []string) error {
	file := c.Args.File
	if file == "" {
		return fmt.Errorf("no file specified")
	}

	bdfs, repoLabel, err := c.getFileSystem()
	if err != nil {
		return err
	}

	if GlobalOpt.Verbose {
		log.Printf("Displaying build file %q for %s", file, repoLabel)
	}

	f, err := bdfs.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(os.Stdout, f)
	return err
}

type BuildDataRemoveCmd struct {
	buildDataSingleRepoCommonOpts

	Recursive bool `short:"r" description:"recursively delete files and dir"`
	All       bool `long:"all" description:"remove all build data (local only)"`
	Args      struct {
		Files []string `name:"FILES" default:"." description:"file to remove"`
	} `positional-args:"yes"`
}

var buildDataRemoveCmd BuildDataRemoveCmd

func (c *BuildDataRemoveCmd) Execute(args []string) error {
	if len(c.Args.Files) == 0 && !c.All {
		return fmt.Errorf("no files specified")
	}

	if c.All {
		if !c.Local {
			return fmt.Errorf("--all and --local must be used together")
		}
		lrepo, err := openLocalRepo()
		if err != nil {
			return err
		}
		if err := os.RemoveAll(filepath.Join(lrepo.RootDir, store.SrclibStoreDir)); err != nil {
			return err
		}
		if err := os.RemoveAll(filepath.Join(lrepo.RootDir, buildstore.BuildDataDirName)); err != nil {
			return err
		}
		return nil
	}

	bdfs, repoLabel, err := c.getFileSystem()
	if err != nil {
		return err
	}

	if GlobalOpt.Verbose {
		log.Printf("Removing build files %v for %s", c.Args.Files, repoLabel)
	}

	vfs := removeLoggedFS{rwvfs.Walkable(bdfs)}

	for _, file := range c.Args.Files {
		if c.Recursive {
			if err := buildstore.RemoveAll(file, vfs); err != nil {
				return err
			}
		} else {
			if err := vfs.Remove(file); err != nil {
				return err
			}
		}
	}
	return nil
}

type removeLoggedFS struct{ rwvfs.WalkableFileSystem }

func (fs removeLoggedFS) Remove(path string) error {
	if err := fs.WalkableFileSystem.Remove(path); err != nil {
		return err
	}
	if GlobalOpt.Verbose {
		log.Printf("Removed %s", path)
	}
	return nil
}

type BuildDataFetchCmd struct {
	buildDataSingleCommitCommonOpts

	DryRun bool `short:"n" long:"dry-run" description:"don't do anything, just show what would be done"`
}

var buildDataFetchCmd BuildDataFetchCmd

func (c *BuildDataFetchCmd) Execute(args []string) error {
	localFS, localRepoLabel, err := c.getLocalFileSystem()
	if err != nil {
		return err
	}

	remoteFS, remoteRepoLabel, repoRevSpec, err := c.getRemoteFileSystem()
	if err != nil {
		return err
	}

	// Use uncached API client because the .srclib-cache already
	// caches it, and we want to be able to stream large files.
	//
	// TODO(sqs): this uncached client isn't authed because it doesn't
	// have the other API client's http.Client or http.RoundTripper
	cl := newAPIClientWithAuth(false)
	remoteFS, err = cl.BuildData.FileSystem(repoRevSpec)
	if err != nil {
		return err
	}

	if GlobalOpt.Verbose {
		log.Printf("Fetching remote build files for %s to %s...", remoteRepoLabel, localRepoLabel)
	}

	// TODO(sqs): check if file exists in local cache and don't fetch it if it does and if it is identical

	par := parallel.NewRun(8)
	w := fs.WalkFS(".", rwvfs.Walkable(remoteFS))
	for w.Step() {
		path := w.Path()
		if err := w.Err(); err != nil {
			if path == "." {
				log.Printf("# No build data to pull from %s", remoteRepoLabel)
				return nil
			}
			return fmt.Errorf("walking remote dir tree: %s", err)
		}
		fi := w.Stat()
		if fi == nil {
			continue
		}
		if !fi.Mode().IsRegular() {
			continue
		}
		par.Do(func() error {
			return fetchFile(remoteFS, localFS, path, fi, c.DryRun)
		})
	}
	if err := par.Wait(); err != nil {
		return fmt.Errorf("error fetching: %s", err)
	}
	return nil
}

func fetchFile(remote vfs.FileSystem, local rwvfs.FileSystem, path string, fi os.FileInfo, dryRun bool) error {
	kb := float64(fi.Size()) / 1024
	if GlobalOpt.Verbose || dryRun {
		log.Printf("Fetching %s (%.1fkb)", path, kb)
	}
	if dryRun {
		return nil
	}

	if err := rwvfs.MkdirAll(local, filepath.Dir(path)); err != nil {
		return err
	}

	rf, err := remote.Open(path)
	if err != nil {
		return fmt.Errorf("remote file: %s", err)
	}
	defer rf.Close()

	lf, err := local.Create(path)
	if err != nil {
		return fmt.Errorf("local file: %s", err)
	}
	defer lf.Close()

	if _, err := io.Copy(lf, rf); err != nil {
		return fmt.Errorf("copy from remote to local: %s", err)
	}

	if GlobalOpt.Verbose {
		log.Printf("Fetched %s (%.1fkb)", path, kb)
	}

	if err := lf.Close(); err != nil {
		return fmt.Errorf("local file: %s", err)
	}
	return nil
}

type BuildDataUploadCmd struct {
	buildDataSingleCommitCommonOpts

	DryRun bool `short:"n" long:"dry-run" description:"don't do anything, just show what would be done"`
}

var buildDataUploadCmd BuildDataUploadCmd

func (c *BuildDataUploadCmd) Execute(args []string) error {
	localFS, localRepoLabel, err := c.getLocalFileSystem()
	if err != nil {
		return err
	}

	remoteFS, remoteRepoLabel, _, err := c.getRemoteFileSystem()
	if err != nil {
		return err
	}

	if GlobalOpt.Verbose {
		log.Printf("Uploading build files from %s to %s...", localRepoLabel, remoteRepoLabel)
	}

	// TODO(sqs): check if file exists remotely and don't upload it if it does and if it is identical

	par := parallel.NewRun(8)
	w := fs.WalkFS(".", rwvfs.Walkable(localFS))
	for w.Step() {
		if err := w.Err(); err != nil {
			return err
		}
		fi := w.Stat()
		if fi == nil {
			continue
		}
		if !fi.Mode().IsRegular() {
			continue
		}
		path := w.Path()
		par.Do(func() error {
			return uploadFile(localFS, remoteFS, path, fi, c.DryRun)
		})
	}
	return par.Wait()
}

func uploadFile(local vfs.FileSystem, remote rwvfs.FileSystem, path string, fi os.FileInfo, dryRun bool) error {
	kb := float64(fi.Size()) / 1024
	if GlobalOpt.Verbose || dryRun {
		log.Printf("Uploading %s (%.1fkb)", path, kb)
	}
	if dryRun {
		return nil
	}

	lf, err := local.Open(path)
	if err != nil {
		return err
	}

	if err := rwvfs.MkdirAll(remote, filepath.Dir(path)); err != nil {
		return err
	}
	rf, err := remote.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := rf.Close(); err != nil {
			log.Println("Error closing after error:", err)
		}
	}()

	if _, err := io.Copy(rf, lf); err != nil {
		return err
	}

	if err := rf.Close(); err != nil {
		return err
	}

	if GlobalOpt.Verbose {
		log.Printf("Uploaded %s (%.1fkb)", path, kb)
	}
	return nil
}
