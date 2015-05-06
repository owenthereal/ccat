package toolchain

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib"
)

// TempDirName is directory under SRCLIBPATH where to store temp directories for toolchains.
const TempDirName = ".tmp"

// TempDir returns toolchains temp directory. Directory is created it doesn't
// exist.
func TempDir(toolchainPath string) (string, error) {
	tc, err := Lookup(toolchainPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf(
				"get toolchain failed: %s (is %s a srclib toolchain repository?)",
				err,
				toolchainPath,
			)
		}
		return "", err
	}

	srclibpathEntry := strings.SplitN(srclib.Path, ":", 2)[0]
	tmpDir := filepath.Join(srclibpathEntry, TempDirName, tc.Path)

	if err := os.MkdirAll(tmpDir, 0700); err != nil {
		return "", err
	}

	return tmpDir, nil
}
