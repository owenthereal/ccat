// +build lgtest

package src_test

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/buildstore"
)

var testdataPath string
var srclibPath string

func init() {
	d, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if !strings.HasSuffix(d, filepath.Join("srclib", "src")) {
		log.Fatalf("directory %s must end with \"srclib"+string(os.PathSeparator)+"src\"", d)
	}
	testdataPath = filepath.Join(d, "..", "testdata")
	srclibPath = filepath.Join(testdataPath, "srclibpath")
}

func cleanup(dir string) {
	cleanCmd := exec.Command("rm", "-rf", buildstore.BuildDataDirName)
	cleanCmd.Dir = filepath.Join(testdataPath, dir)
	cleanCmd.Run()
}

func TestDoAll_cached_sample(t *testing.T) {
	cleanup("go-cache")

	var gitCmd *exec.Cmd
	var srcCmd *exec.Cmd
	gitCmd = exec.Command("git", "checkout", "071610bf3a597bc41aae05e27c5407444b7ea0d1")
	gitCmd.Dir = filepath.Join(testdataPath, "go-cached")
	if o, err := gitCmd.CombinedOutput(); err != nil {
		t.Fatal(string(o), err)
	}
	srcCmd = exec.Command("src", "do-all")
	srcCmd.Dir = filepath.Join(testdataPath, "go-cached")
	srcCmd.Env = append([]string{"SRCLIBPATH=" + srclibPath}, os.Environ()...)
	if o, err := srcCmd.CombinedOutput(); err != nil {
		t.Fatal(string(o), err)
	}
	gitCmd = exec.Command("git", "checkout", "34dd0f240fe12cdd8c9c6e24620cc0013518a55e")
	gitCmd.Dir = filepath.Join(testdataPath, "go-cached")
	if o, err := gitCmd.CombinedOutput(); err != nil {
		t.Fatal(string(o), err)
	}
	srcCmd = exec.Command("src", "do-all")
	srcCmd.Dir = filepath.Join(testdataPath, "go-cached")
	srcCmd.Env = append([]string{"SRCLIBPATH=" + srclibPath}, os.Environ()...)
	if o, err := srcCmd.CombinedOutput(); err != nil {
		t.Fatal(string(o), err)
	}
	firstOne, err := ioutil.ReadFile(filepath.Join(testdataPath, "go-cached", buildstore.BuildDataDirName, "071610bf3a597bc41aae05e27c5407444b7ea0d1", "one", "sample.graph.json"))
	if err != nil {
		t.Fatal(err)
	}
	secondOne, err := ioutil.ReadFile(filepath.Join(testdataPath, "go-cached", buildstore.BuildDataDirName, "34dd0f240fe12cdd8c9c6e24620cc0013518a55e", "one", "sample.graph.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(firstOne) != string(secondOne) {
		t.Error("Source unit \"one\" should have been cached: string(firstOne) != string(secondOne)")
	}
	firstTwo, err := ioutil.ReadFile(filepath.Join(testdataPath, "go-cached", buildstore.BuildDataDirName, "071610bf3a597bc41aae05e27c5407444b7ea0d1", "two", "sample.graph.json"))
	if err != nil {
		t.Fatal(err)
	}
	secondTwo, err := ioutil.ReadFile(filepath.Join(testdataPath, "go-cached", buildstore.BuildDataDirName, "34dd0f240fe12cdd8c9c6e24620cc0013518a55e", "two", "sample.graph.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(firstTwo) == string(secondTwo) {
		t.Error("Source unit \"two\" should not be cached: string(firstTwo) == string(secondTwo)")
	}
	cleanup("go-cache")
}
