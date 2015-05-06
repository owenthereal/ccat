package router

import "github.com/jingweno/ccat/Godeps/_workspace/src/github.com/sourcegraph/mux"

const (
	Build            = "build"
	BuildDequeueNext = "build.dequeue-next"
	BuildUpdate      = "build.update"
	BuildLog         = "build.log"
	Builds           = "builds"
	BuildTasks       = "build.tasks"
	BuildTaskUpdate  = "build.task"
	BuildTasksCreate = "build.tasks.create"
	BuildTaskLog     = "build.task.log"

	Org               = "org"
	OrgMembers        = "org.members"
	OrgSettings       = "org.settings"
	OrgSettingsUpdate = "org.settings.update"

	Users                 = "users"
	User                  = "user"
	UserOrgs              = "user.orgs"
	UserAuthors           = "user.authors"
	UserClients           = "user.clients"
	UserEmails            = "user.emails"
	UserFromGitHub        = "user.from-github"
	UserRepoContributions = "user.repo-contributions"
	UserRepoDependencies  = "user.repo-dependencies"
	UserRepoDependents    = "user.repo-dependents"
	UserRefreshProfile    = "user.refresh-profile"
	UserSettings          = "user.settings"
	UserSettingsUpdate    = "user.settings.update"
	UserComputeStats      = "user.compute-stats"

	Person = "person"

	RepoPullRequests              = "repo.pull-requests"
	RepoPullRequest               = "repo.pull-request"
	RepoPullRequestMerge          = "repo.pull-request.merge"
	RepoPullRequestComments       = "repo.pull-request.comments"
	RepoPullRequestCommentsCreate = "repo.pull-request.comments.create"
	RepoPullRequestCommentsEdit   = "repo.pull-request.comments.edit"
	RepoPullRequestCommentsDelete = "repo.pull-request.comments.delete"

	RepoIssues              = "repo.issues"
	RepoIssue               = "repo.issue"
	RepoIssueComments       = "repo.issue.comments"
	RepoIssueCommentsCreate = "repo.issue.comments.create"
	RepoIssueCommentsEdit   = "repo.issue.comments.edit"
	RepoIssueCommentsDelete = "repo.issue.comments.delete"

	Repos              = "repos"
	ReposCreate        = "repos.create"
	ReposGetOrCreate   = "repos.get-or-create"
	Repo               = "repo"
	RepoAuthors        = "repo.authors"
	RepoClients        = "repo.clients"
	RepoDependents     = "repo.dependents"
	RepoDependencies   = "repo.dependencies"
	RepoBadge          = "repo.badge"
	RepoBadges         = "repo.badges"
	RepoCounter        = "repo.counter"
	RepoCounters       = "repo.counters"
	RepoReadme         = "repo.readme"
	RepoBuildsCreate   = "repo.builds.create"
	RepoBuildDataEntry = "repo.build-data.entry"
	RepoTreeEntry      = "repo.tree.entry"
	RepoTreeSearch     = "repo.tree.search"
	RepoRefreshProfile = "repo.refresh-profile"
	RepoRefreshVCSData = "repo.refresh-vcs-data"
	RepoComputeStats   = "repo.compute-stats"

	RepoSettings       = "repo.settings"
	RepoSettingsUpdate = "repo.settings.update"

	RepoStats          = "repo.stats"
	RepoCombinedStatus = "repo.combined-status"
	RepoStatusCreate   = "repo.status.create"

	RepoBuild = "repo.build"

	RepoCommits        = "repo.commits"
	RepoCommit         = "repo.commit"
	RepoCompareCommits = "repo.compare-commits"
	RepoTags           = "repo.tags"
	RepoBranches       = "repo.branches"

	Search         = "search"
	SearchComplete = "search.complete"

	SearchSuggestions = "search.suggestions"

	Snippet = "snippet"

	Defs          = "defs"
	Def           = "def"
	DefRefs       = "def.refs"
	DefExamples   = "def.examples"
	DefAuthors    = "def.authors"
	DefClients    = "def.clients"
	DefDependents = "def.dependents"
	DefVersions   = "def.versions"

	Delta                   = "delta"
	DeltaUnits              = "delta.units"
	DeltaDefs               = "delta.defs"
	DeltaDependencies       = "delta.dependencies"
	DeltaFiles              = "delta.files"
	DeltaAffectedAuthors    = "delta.affected-authors"
	DeltaAffectedClients    = "delta.affected-clients"
	DeltaAffectedDependents = "delta.affected-dependents"
	DeltaReviewers          = "delta.reviewers"
	DeltasIncoming          = "deltas.incoming"

	Unit  = "unit"
	Units = "units"

	Markdown = "markdown"

	ExtGitHubReceiveWebhook = "ext.github.receive-webhook"

	// Redirects for old routes.
	RedirectOldRepoBadgesAndCounters = "repo.redirect-old-badges-and-counters"
)

// NewAPIRouter creates a new API router with route URL pattern definitions but
// no handlers attached to the routes.
//
// It is in a separate package from app so that other packages may use it to
// generate URLs without resulting in Go import cycles (and so we can release
// the router as open-source to support our client library).
func NewAPIRouter(base *mux.Router) *mux.Router {
	if base == nil {
		base = mux.NewRouter()
	}

	base.StrictSlash(true)

	base.Path("/builds").Methods("GET").Name(Builds)
	builds := base.PathPrefix("/builds").Subrouter()
	builds.Path("/next").Methods("POST").Name(BuildDequeueNext)
	buildPath := "/{BID}"
	builds.Path(buildPath).Methods("GET").Name(Build)
	builds.Path(buildPath).Methods("PUT").Name(BuildUpdate)
	build := builds.PathPrefix(buildPath).Subrouter()
	build.Path("/log").Methods("GET").Name(BuildLog)
	build.Path("/tasks").Methods("GET").Name(BuildTasks)
	build.Path("/tasks").Methods("POST").Name(BuildTasksCreate)
	build.Path("/tasks/{TaskID}").Methods("PUT").Name(BuildTaskUpdate)
	build.Path("/tasks/{TaskID}/log").Methods("GET").Name(BuildTaskLog)

	base.Path("/repos").Methods("GET").Name(Repos)
	base.Path("/repos").Methods("POST").Name(ReposCreate)

	base.Path("/repos/github.com/{owner:[^/]+}/{repo:[^/]+}/{what:(?:badges|counters)}/{which}.{Format}").Methods("GET").Name(RedirectOldRepoBadgesAndCounters)

	repoRev := base.PathPrefix(`/repos/` + RepoRevSpecPattern).PostMatchFunc(FixRepoRevSpecVars).BuildVarsFunc(PrepareRepoRevSpecRouteVars).Subrouter()
	repoRev.Path("/.stats").Methods("PUT").Name(RepoComputeStats)
	repoRev.Path("/.stats").Methods("GET").Name(RepoStats)
	repoRev.Path("/.status").Methods("GET").Name(RepoCombinedStatus)
	repoRev.Path("/.status").Methods("POST").Name(RepoStatusCreate)
	repoRev.Path("/.authors").Methods("GET").Name(RepoAuthors)
	repoRev.Path("/.readme").Methods("GET").Name(RepoReadme)
	repoRev.Path("/.build").Methods("GET").Name(RepoBuild)
	repoRev.Path("/.builds").Methods("POST").Name(RepoBuildsCreate)
	repoRev.Path("/.dependencies").Methods("GET").Name(RepoDependencies)
	repoRev.PathPrefix("/.build-data"+TreeEntryPathPattern).PostMatchFunc(FixTreeEntryVars).BuildVarsFunc(PrepareTreeEntryRouteVars).Methods("GET", "HEAD", "PUT", "DELETE").Name(RepoBuildDataEntry)
	repoRev.Path("/.badges/{Badge}.{Format}").Methods("GET").Name(RepoBadge)

	// repo contains routes that are NOT specific to a revision. In these routes, the URL may not contain a revspec after the repo (that is, no "github.com/foo/bar@myrevspec").
	repoPath := `/repos/` + RepoSpecPathPattern
	base.Path(repoPath).Methods("GET").Name(Repo)
	base.Path(repoPath).Methods("PUT").Name(ReposGetOrCreate)
	repo := base.PathPrefix(repoPath).Subrouter()
	repo.Path("/.clients").Methods("GET").Name(RepoClients)
	repo.Path("/.dependents").Methods("GET").Name(RepoDependents)
	repo.Path("/.external-profile").Methods("PUT").Name(RepoRefreshProfile)
	repo.Path("/.vcs-data").Methods("PUT").Name(RepoRefreshVCSData)
	repo.Path("/.settings").Methods("GET").Name(RepoSettings)
	repo.Path("/.settings").Methods("PUT").Name(RepoSettingsUpdate)
	repo.Path("/.commits").Methods("GET").Name(RepoCommits)
	repo.Path("/.commits/{Rev:" + PathComponentNoLeadingDot + "}/.compare").Methods("GET").Name(RepoCompareCommits)
	repo.Path("/.commits/{Rev:" + PathComponentNoLeadingDot + "}").Methods("GET").Name(RepoCommit)
	repo.Path("/.branches").Methods("GET").Name(RepoBranches)
	repo.Path("/.tags").Methods("GET").Name(RepoTags)
	repo.Path("/.badges").Methods("GET").Name(RepoBadges)
	repo.Path("/.counters").Methods("GET").Name(RepoCounters)
	repo.Path("/.counters/{Counter}.{Format}").Methods("GET").Name(RepoCounter)

	repo.Path("/.pulls").Methods("GET").Name(RepoPullRequests)
	pullPath := "/.pulls/{Pull}"
	repo.Path(pullPath).Methods("GET").Name(RepoPullRequest)
	pull := repo.PathPrefix(pullPath).Subrouter()
	pull.Path("/merge").Methods("PUT").Name(RepoPullRequestMerge)
	pull.Path("/comments").Methods("GET").Name(RepoPullRequestComments)
	pull.Path("/comments").Methods("POST").Name(RepoPullRequestCommentsCreate)
	pull.Path("/comments/{CommentID}").Methods("PATCH", "PUT").Name(RepoPullRequestCommentsEdit)
	pull.Path("/comments/{CommentID}").Methods("DELETE").Name(RepoPullRequestCommentsDelete)

	repo.Path("/.issues").Methods("GET").Name(RepoIssues)
	issuePath := "/.issues/{Issue}"
	repo.Path(issuePath).Methods("GET").Name(RepoIssue)
	issue := repo.PathPrefix(issuePath).Subrouter()
	issue.Path("/comments").Methods("GET").Name(RepoIssueComments)
	issue.Path("/comments").Methods("POST").Name(RepoIssueCommentsCreate)
	issue.Path("/comments/{CommentID}").Methods("PATCH", "PUT").Name(RepoIssueCommentsEdit)
	issue.Path("/comments/{CommentID}").Methods("DELETE").Name(RepoIssueCommentsDelete)

	deltaPath := "/.deltas/{Rev:.+}..{DeltaHeadRev:" + PathComponentNoLeadingDot + "}"
	repo.Path(deltaPath).Methods("GET").Name(Delta)
	deltas := repo.PathPrefix(deltaPath).Subrouter()
	deltas.Path("/.units").Methods("GET").Name(DeltaUnits)
	deltas.Path("/.defs").Methods("GET").Name(DeltaDefs)
	deltas.Path("/.dependencies").Methods("GET").Name(DeltaDependencies)
	deltas.Path("/.files").Methods("GET").Name(DeltaFiles)
	deltas.Path("/.affected-authors").Methods("GET").Name(DeltaAffectedAuthors)
	deltas.Path("/.affected-clients").Methods("GET").Name(DeltaAffectedClients)
	deltas.Path("/.affected-dependents").Methods("GET").Name(DeltaAffectedDependents)
	deltas.Path("/.reviewers").Methods("GET").Name(DeltaReviewers)

	repo.Path("/.deltas-incoming").Methods("GET").Name(DeltasIncoming)

	// See router_util/tree_route.go for an explanation of how we match tree
	// entry routes.
	repoRev.Path("/.tree" + TreeEntryPathPattern).PostMatchFunc(FixTreeEntryVars).BuildVarsFunc(PrepareTreeEntryRouteVars).Methods("GET").Name(RepoTreeEntry)

	repoRev.Path("/.tree-search").Methods("GET").Name(RepoTreeSearch)

	base.Path(`/people/` + PersonSpecPattern).Methods("GET").Name(Person)

	base.Path("/users").Methods("GET").Name(Users)
	userPath := `/users/` + UserSpecPattern
	base.Path(userPath).Methods("GET").Name(User)
	user := base.PathPrefix(userPath).Subrouter()
	user.Path("/orgs").Methods("GET").Name(UserOrgs)
	user.Path("/clients").Methods("GET").Name(UserClients)
	user.Path("/authors").Methods("GET").Name(UserAuthors)
	user.Path("/emails").Methods("GET").Name(UserEmails)
	user.Path("/repo-contributions").Methods("GET").Name(UserRepoContributions)
	user.Path("/repo-dependencies").Methods("GET").Name(UserRepoDependencies)
	user.Path("/repo-dependents").Methods("GET").Name(UserRepoDependents)
	user.Path("/external-profile").Methods("PUT").Name(UserRefreshProfile)
	user.Path("/stats").Methods("PUT").Name(UserComputeStats)
	user.Path("/settings").Methods("GET").Name(UserSettings)
	user.Path("/settings").Methods("PUT").Name(UserSettingsUpdate)
	base.Path("/external-users/github/{GitHubUserSpec}").Methods("GET").Name(UserFromGitHub)

	orgPath := "/orgs/{OrgSpec}"
	base.Path(orgPath).Methods("GET").Name(Org)
	org := base.PathPrefix(orgPath).Subrouter()
	org.Path("/settings").Methods("GET").Name(OrgSettings)
	org.Path("/settings").Methods("PUT").Name(OrgSettingsUpdate)
	org.Path("/members").Methods("GET").Name(OrgMembers)

	base.Path("/search").Methods("GET").Name(Search)
	base.Path("/search/complete").Methods("GET").Name(SearchComplete)
	base.Path("/search/suggestions").Methods("GET").Name(SearchSuggestions)

	base.Path("/snippet").Methods("GET", "POST", "ORIGIN").Name(Snippet)

	base.Path("/.defs").Methods("GET").Name(Defs)

	// See router_util/def_route.go for an explanation of how we match def
	// routes.
	defPath := `/.defs/` + DefPathPattern
	repoRev.Path(defPath).Methods("GET").PostMatchFunc(FixDefUnitVars).BuildVarsFunc(PrepareDefRouteVars).Name(Def)
	def := repoRev.PathPrefix(defPath).PostMatchFunc(FixDefUnitVars).BuildVarsFunc(PrepareDefRouteVars).Subrouter()
	def.Path("/.refs").Methods("GET").Name(DefRefs)
	def.Path("/.examples").Methods("GET").Name(DefExamples)
	def.Path("/.authors").Methods("GET").Name(DefAuthors)
	def.Path("/.clients").Methods("GET").Name(DefClients)
	def.Path("/.dependents").Methods("GET").Name(DefDependents)
	def.Path("/.versions").Methods("GET").Name(DefVersions)

	base.Path("/.units").Methods("GET").Name(Units)
	unitPath := `/.units/{UnitType}/{Unit:.*}`
	repoRev.Path(unitPath).Methods("GET").Name(Unit)

	base.Path("/markdown").Methods("POST").Name(Markdown)

	base.Path("/ext/github/webhook").Methods("POST").Name(ExtGitHubReceiveWebhook)

	if ExtraConfig != nil {
		ExtraConfig(base, user)
	}

	return base
}

// ExtraConfig, if non-nil, is called by NewAPIRouter with the
// *mux.Router after setting up routes in this package and before
// returning it. It can be used by external packages that use this API
// router and want to add additional routes to it.
var ExtraConfig func(base, user *mux.Router)
