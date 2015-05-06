package router

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/sourcegraph/mux"
)

// TreeEntryPathPattern is the path pattern for tree entries.
var TreeEntryPathPattern = `{Path:(?:/.*)*}`

// FixTreeEntryVars is a mux.PostMatchFunc that cleans and normalizes the path to a tree entry.
func FixTreeEntryVars(req *http.Request, match *mux.RouteMatch, r *mux.Route) {
	path := filepath.Clean(strings.TrimPrefix(match.Vars["Path"], "/"))
	if path == "" || path == "." {
		match.Vars["Path"] = "."
	} else {
		match.Vars["Path"] = path
	}
}

// PrepareTreeEntryRouteVars is a mux.BuildVarsFunc that converts from a cleaned
// and normalized Path to a Path that we use to generate tree entry URLs.
func PrepareTreeEntryRouteVars(vars map[string]string) map[string]string {
	if path := vars["Path"]; path == "." {
		vars["Path"] = ""
	} else {
		vars["Path"] = "/" + path
	}
	return vars
}
