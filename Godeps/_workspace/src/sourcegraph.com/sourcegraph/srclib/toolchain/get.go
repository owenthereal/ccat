package toolchain

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib"
)

// Get downloads the toolchain named by the toolchain path (if it does not
// already exist in the SRCLIBPATH). If update is true, it uses the network to
// update the toolchain.
//
// Assumes that the clone URL is "https://" + path + ".git".
func Get(path string, update bool) (*Info, error) {
	path = filepath.Clean(path)
	if tc, err := Lookup(path); !os.IsNotExist(err) {
		return tc, err
	}

	dir := strings.SplitN(srclib.Path, ":", 2)[0]
	toolchainDir := filepath.Join(dir, path)

	if fi, err := os.Stat(toolchainDir); os.IsNotExist(err) {
		// older gits don't heed git https redirects, so manually substitute in
		// the github.com clone url for sourcegraph.com clone urls
		var substitutedPath string
		if strings.HasPrefix(path, "sourcegraph.com/") {
			substitutedPath = "github.com/" + strings.TrimPrefix(path, "sourcegraph.com/")
		} else {
			substitutedPath = path
		}
		cloneURL := "https://" + substitutedPath + ".git"
		cmd := exec.Command("git", "clone", cloneURL, toolchainDir)
		cmd.Stdout, cmd.Stderr = os.Stderr, os.Stderr
		if err := cmd.Run(); err != nil {
			return nil, err
		}
	} else if fi.Mode().IsDir() {
		cmd := exec.Command("git", "pull", "origin", "master")
		cmd.Dir = toolchainDir
		cmd.Stdout, cmd.Stderr = os.Stderr, os.Stderr
		if err := cmd.Run(); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	tc, err := Lookup(path)
	if err != nil {
		return nil, fmt.Errorf("get toolchain failed: %s (is %s a srclib toolchain repository?)", err, path)
	}
	return tc, nil
}
