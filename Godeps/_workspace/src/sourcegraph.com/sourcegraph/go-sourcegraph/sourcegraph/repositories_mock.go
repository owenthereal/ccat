package sourcegraph

import (
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-vcs/vcs"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/vcsstore/vcsclient"
)

type MockReposService struct {
	Get_               func(repo RepoSpec, opt *RepoGetOptions) (*Repo, Response, error)
	GetStats_          func(repo RepoRevSpec) (RepoStats, Response, error)
	GetCombinedStatus_ func(spec RepoRevSpec) (*CombinedStatus, Response, error)
	CreateStatus_      func(spec RepoRevSpec, st RepoStatus) (*RepoStatus, Response, error)
	GetOrCreate_       func(repo RepoSpec, opt *RepoGetOptions) (*Repo, Response, error)
	GetSettings_       func(repo RepoSpec) (*RepoSettings, Response, error)
	UpdateSettings_    func(repo RepoSpec, settings RepoSettings) (Response, error)
	RefreshProfile_    func(repo RepoSpec) (Response, error)
	RefreshVCSData_    func(repo RepoSpec) (Response, error)
	ComputeStats_      func(repo RepoRevSpec) (Response, error)
	GetBuild_          func(repo RepoRevSpec, opt *RepoGetBuildOptions) (*RepoBuildInfo, Response, error)
	Create_            func(newRepoSpec NewRepoSpec) (*Repo, Response, error)
	GetReadme_         func(repo RepoRevSpec) (*vcsclient.TreeEntry, Response, error)
	List_              func(opt *RepoListOptions) ([]*Repo, Response, error)
	ListCommits_       func(repo RepoSpec, opt *RepoListCommitsOptions) ([]*Commit, Response, error)
	GetCommit_         func(rev RepoRevSpec, opt *RepoGetCommitOptions) (*Commit, Response, error)
	ListBranches_      func(repo RepoSpec, opt *RepoListBranchesOptions) ([]*vcs.Branch, Response, error)
	ListTags_          func(repo RepoSpec, opt *RepoListTagsOptions) ([]*vcs.Tag, Response, error)
	ListBadges_        func(repo RepoSpec) ([]*Badge, Response, error)
	ListCounters_      func(repo RepoSpec) ([]*Counter, Response, error)
	ListAuthors_       func(repo RepoRevSpec, opt *RepoListAuthorsOptions) ([]*AugmentedRepoAuthor, Response, error)
	ListClients_       func(repo RepoSpec, opt *RepoListClientsOptions) ([]*AugmentedRepoClient, Response, error)
	ListDependencies_  func(repo RepoRevSpec, opt *RepoListDependenciesOptions) ([]*AugmentedRepoDependency, Response, error)
	ListDependents_    func(repo RepoSpec, opt *RepoListDependentsOptions) ([]*AugmentedRepoDependent, Response, error)
	ListByContributor_ func(user UserSpec, opt *RepoListByContributorOptions) ([]*AugmentedRepoContribution, Response, error)
	ListByClient_      func(user UserSpec, opt *RepoListByClientOptions) ([]*AugmentedRepoUsageByClient, Response, error)
	ListByRefdAuthor_  func(user UserSpec, opt *RepoListByRefdAuthorOptions) ([]*AugmentedRepoUsageOfAuthor, Response, error)
}

func (s MockReposService) Get(repo RepoSpec, opt *RepoGetOptions) (*Repo, Response, error) {
	return s.Get_(repo, opt)
}

func (s MockReposService) GetStats(repo RepoRevSpec) (RepoStats, Response, error) {
	return s.GetStats_(repo)
}

func (s MockReposService) GetCombinedStatus(spec RepoRevSpec) (*CombinedStatus, Response, error) {
	return s.GetCombinedStatus_(spec)
}

func (s MockReposService) CreateStatus(spec RepoRevSpec, st RepoStatus) (*RepoStatus, Response, error) {
	return s.CreateStatus_(spec, st)
}

func (s MockReposService) GetOrCreate(repo RepoSpec, opt *RepoGetOptions) (*Repo, Response, error) {
	return s.GetOrCreate_(repo, opt)
}

func (s MockReposService) GetSettings(repo RepoSpec) (*RepoSettings, Response, error) {
	return s.GetSettings_(repo)
}

func (s MockReposService) UpdateSettings(repo RepoSpec, settings RepoSettings) (Response, error) {
	return s.UpdateSettings_(repo, settings)
}

func (s MockReposService) RefreshProfile(repo RepoSpec) (Response, error) {
	return s.RefreshProfile_(repo)
}

func (s MockReposService) RefreshVCSData(repo RepoSpec) (Response, error) {
	return s.RefreshVCSData_(repo)
}

func (s MockReposService) ComputeStats(repo RepoRevSpec) (Response, error) {
	return s.ComputeStats_(repo)
}

func (s MockReposService) GetBuild(repo RepoRevSpec, opt *RepoGetBuildOptions) (*RepoBuildInfo, Response, error) {
	return s.GetBuild_(repo, opt)
}

func (s MockReposService) Create(newRepoSpec NewRepoSpec) (*Repo, Response, error) {
	return s.Create_(newRepoSpec)
}

func (s MockReposService) GetReadme(repo RepoRevSpec) (*vcsclient.TreeEntry, Response, error) {
	return s.GetReadme_(repo)
}

func (s MockReposService) List(opt *RepoListOptions) ([]*Repo, Response, error) { return s.List_(opt) }

func (s MockReposService) ListCommits(repo RepoSpec, opt *RepoListCommitsOptions) ([]*Commit, Response, error) {
	return s.ListCommits_(repo, opt)
}

func (s MockReposService) GetCommit(rev RepoRevSpec, opt *RepoGetCommitOptions) (*Commit, Response, error) {
	return s.GetCommit_(rev, opt)
}

func (s MockReposService) ListBranches(repo RepoSpec, opt *RepoListBranchesOptions) ([]*vcs.Branch, Response, error) {
	return s.ListBranches_(repo, opt)
}

func (s MockReposService) ListTags(repo RepoSpec, opt *RepoListTagsOptions) ([]*vcs.Tag, Response, error) {
	return s.ListTags_(repo, opt)
}

func (s MockReposService) ListBadges(repo RepoSpec) ([]*Badge, Response, error) {
	return s.ListBadges_(repo)
}

func (s MockReposService) ListCounters(repo RepoSpec) ([]*Counter, Response, error) {
	return s.ListCounters_(repo)
}

func (s MockReposService) ListAuthors(repo RepoRevSpec, opt *RepoListAuthorsOptions) ([]*AugmentedRepoAuthor, Response, error) {
	return s.ListAuthors_(repo, opt)
}

func (s MockReposService) ListClients(repo RepoSpec, opt *RepoListClientsOptions) ([]*AugmentedRepoClient, Response, error) {
	return s.ListClients_(repo, opt)
}

func (s MockReposService) ListDependencies(repo RepoRevSpec, opt *RepoListDependenciesOptions) ([]*AugmentedRepoDependency, Response, error) {
	return s.ListDependencies_(repo, opt)
}

func (s MockReposService) ListDependents(repo RepoSpec, opt *RepoListDependentsOptions) ([]*AugmentedRepoDependent, Response, error) {
	return s.ListDependents_(repo, opt)
}

func (s MockReposService) ListByContributor(user UserSpec, opt *RepoListByContributorOptions) ([]*AugmentedRepoContribution, Response, error) {
	return s.ListByContributor_(user, opt)
}

func (s MockReposService) ListByClient(user UserSpec, opt *RepoListByClientOptions) ([]*AugmentedRepoUsageByClient, Response, error) {
	return s.ListByClient_(user, opt)
}

func (s MockReposService) ListByRefdAuthor(user UserSpec, opt *RepoListByRefdAuthorOptions) ([]*AugmentedRepoUsageOfAuthor, Response, error) {
	return s.ListByRefdAuthor_(user, opt)
}
