package router

import (
	"net/http"
	"strings"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/sourcegraph/mux"
)

// RepoSpecPattern is the path pattern for encoding RepoSpec.
var RepoSpecPathPattern = `{RepoSpec:(?:(?:[^/.@][^/@]*/)+(?:[^/.@][^/@]*))|(?:R\$\d+)}`

// RepoRevSpecPattern is the path pattern for encoding RepoRevSpec.
var RepoRevSpecPattern = RepoSpecPathPattern + `{Rev:(?:@` + PathComponentNoLeadingDot + `)?}`

// PathComponentNoLeadingDot is a pattern that matches any string that doesn't contain "/.".
var PathComponentNoLeadingDot = `(?:[^/]*(?:/` + noDotDotOrSlash + `)*)`

// noDotDotOrSlash matches a single path component and does not permit
// "..".
const noDotDotOrSlash = `(?:[^/.]+[^/]*)+`

// FixRepoRevSpecVars is a mux.PostMatchFunc that cleans and normalizes the
// RepoRevSpecPattern vars.
func FixRepoRevSpecVars(req *http.Request, match *mux.RouteMatch, r *mux.Route) {
	if rev, present := match.Vars["Rev"]; present {
		if rev == "" {
			delete(match.Vars, "Rev")
		} else {
			match.Vars["Rev"] = strings.TrimPrefix(rev, "@")
		}
	}
}

// PrepareRepoRevSpecRouteVars is a mux.BuildVarsFunc that converts
// from a RepoRevSpec's route vars to components used to generate
// routes.
func PrepareRepoRevSpecRouteVars(vars map[string]string) map[string]string {
	if rev, present := vars["Rev"]; !present {
		vars["Rev"] = ""
	} else if rev != "" {
		vars["Rev"] = "@" + rev
	}
	return vars
}
