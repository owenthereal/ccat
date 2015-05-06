package src

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"sourcegraph.com/sourcegraph/makex"

	"strings"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/buildstore"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/config"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/flagutil"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/plan"
)

func init() {
	c, err := CLI.AddCommand("make",
		"plans and executes plan",
		`Generates a plan (in Makefile form, in memory) for analyzing the tree and executes the plan. `,
		&makeCmd,
	)
	if err != nil {
		log.Fatal(err)
	}

	setDefaultRepoURIOpt(c)
	setDefaultRepoSubdirOpt(c)
}

type MakeCmd struct {
	config.Options

	ToolchainExecOpt `group:"execution"`
	BuildCacheOpt    `group:"build cache"`

	Quiet  bool `short:"q" long:"quiet" description:"silence all output"`
	DryRun bool `short:"n" long:"dry-run" description:"print what would be done and exit"`

	Dir Directory `short:"C" long:"directory" description:"change to DIR before doing anything" value-name:"DIR"`

	Args struct {
		Goals []string `name:"GOALS..." description:"Makefile targets to build (default: all)"`
	} `positional-args:"yes"`
}

var makeCmd MakeCmd

func (c *MakeCmd) Execute(args []string) error {
	if c.Dir != "" {
		if err := os.Chdir(c.Dir.String()); err != nil {
			return err
		}
	}

	mf, err := CreateMakefile(c.ToolchainExecOpt, c.BuildCacheOpt)
	if err != nil {
		return err
	}

	goals := c.Args.Goals
	if len(goals) == 0 {
		if defaultRule := mf.DefaultRule(); defaultRule != nil {
			goals = []string{defaultRule.Target()}
		}
	}

	mkConf := &makex.Default
	mk := mkConf.NewMaker(mf, goals...)

	if c.Quiet {
		mk.RuleOutput = func(r makex.Rule) (out io.WriteCloser, err io.WriteCloser, logger *log.Logger) {
			return nopWriteCloser{}, nopWriteCloser{},
				log.New(nopWriteCloser{}, "", 0)
		}
	}

	if c.DryRun {
		return mk.DryRun(os.Stdout)
	}
	return mk.Run()
}

// CreateMakefile creates a Makefile to build a tree. The cwd should
// be the root of the tree you want to make (due to some probably
// unnecessary assumptions that CreateMaker makes).
func CreateMakefile(execOpt ToolchainExecOpt, cacheOpt BuildCacheOpt) (*makex.Makefile, error) {
	localRepo, err := OpenRepo(".")
	if err != nil {
		return nil, err
	}
	buildStore, err := buildstore.LocalRepo(localRepo.RootDir)
	if err != nil {
		return nil, err
	}

	treeConfig, err := config.ReadCached(buildStore.Commit(localRepo.CommitID))
	if err != nil {
		return nil, err
	}
	if len(treeConfig.SourceUnits) == 0 {
		log.Println("No source unit files found. Did you mean to run `src config`? (This is not an error; it just means that src didn't find anything to build or analyze here.)")
	}

	toolchainExecOptArgs, err := flagutil.MarshalArgs(&execOpt)
	if err != nil {
		return nil, err
	}

	// TODO(sqs): buildDataDir is hardcoded.
	buildDataDir := filepath.Join(buildstore.BuildDataDirName, localRepo.CommitID)
	mf, err := plan.CreateMakefile(buildDataDir, buildStore, localRepo.VCSType, treeConfig, plan.Options{
		ToolchainExecOpt: strings.Join(toolchainExecOptArgs, " "),
		NoCache:          cacheOpt.NoCacheWrite,
	})
	if err != nil {
		return nil, err
	}
	return mf, nil
}
