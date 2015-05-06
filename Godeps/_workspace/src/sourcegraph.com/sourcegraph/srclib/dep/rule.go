package dep

import (
	"fmt"
	"path/filepath"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/buildstore"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/config"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/plan"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/toolchain"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/unit"
	"sourcegraph.com/sourcegraph/makex"
)

const depresolveOp = "depresolve"

func init() {
	plan.RegisterRuleMaker(depresolveOp, makeDepRules)
	buildstore.RegisterDataType("depresolve", []*ResolvedDep{})
}

func makeDepRules(c *config.Tree, dataDir string, existing []makex.Rule, opt plan.Options) ([]makex.Rule, error) {
	const op = depresolveOp
	var rules []makex.Rule
	for _, u := range c.SourceUnits {
		toolRef := u.Ops[op]
		if toolRef == nil {
			choice, err := toolchain.ChooseTool(depresolveOp, u.Type)
			if err != nil {
				return nil, err
			}
			toolRef = choice
		}

		rules = append(rules, &ResolveDepsRule{dataDir, u, toolRef, opt})
	}
	return rules, nil
}

type ResolveDepsRule struct {
	dataDir string
	Unit    *unit.SourceUnit
	Tool    *srclib.ToolRef
	opt     plan.Options
}

func (r *ResolveDepsRule) Target() string {
	return filepath.Join(r.dataDir, plan.SourceUnitDataFilename([]*ResolvedDep{}, r.Unit))
}

func (r *ResolveDepsRule) Prereqs() []string {
	return []string{filepath.Join(r.dataDir, plan.SourceUnitDataFilename(unit.SourceUnit{}, r.Unit))}
}

func (r *ResolveDepsRule) Recipes() []string {
	return []string{
		fmt.Sprintf("src tool %s %q %q < $^ 1> $@", r.opt.ToolchainExecOpt, r.Tool.Toolchain, r.Tool.Subcmd),
	}
}

func (r *ResolveDepsRule) SourceUnit() *unit.SourceUnit { return r.Unit }
