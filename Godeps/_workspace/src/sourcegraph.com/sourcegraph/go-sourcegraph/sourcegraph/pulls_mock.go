package sourcegraph

type MockPullRequestsService struct {
	Get_           func(pull PullRequestSpec, opt *PullRequestGetOptions) (*PullRequest, Response, error)
	ListByRepo_    func(repo RepoSpec, opt *PullRequestListOptions) ([]*PullRequest, Response, error)
	ListComments_  func(pull PullRequestSpec, opt *PullRequestListCommentsOptions) ([]*PullRequestComment, Response, error)
	CreateComment_ func(pull PullRequestSpec, comment *PullRequestComment) (*PullRequestComment, Response, error)
	EditComment_   func(pull PullRequestSpec, comment *PullRequestComment) (*PullRequestComment, Response, error)
	DeleteComment_ func(pull PullRequestSpec, commentID int) (Response, error)
	Merge_         func(pull PullRequestSpec, mergeRequest *PullRequestMergeRequest) (*PullRequestMergeResult, Response, error)
}

func (s MockPullRequestsService) Get(pull PullRequestSpec, opt *PullRequestGetOptions) (*PullRequest, Response, error) {
	return s.Get_(pull, opt)
}

func (s MockPullRequestsService) ListByRepo(repo RepoSpec, opt *PullRequestListOptions) ([]*PullRequest, Response, error) {
	return s.ListByRepo_(repo, opt)
}

func (s MockPullRequestsService) ListComments(pull PullRequestSpec, opt *PullRequestListCommentsOptions) ([]*PullRequestComment, Response, error) {
	return s.ListComments_(pull, opt)
}

func (s MockPullRequestsService) CreateComment(pull PullRequestSpec, comment *PullRequestComment) (*PullRequestComment, Response, error) {
	return s.CreateComment_(pull, comment)
}

func (s MockPullRequestsService) EditComment(pull PullRequestSpec, comment *PullRequestComment) (*PullRequestComment, Response, error) {
	return s.EditComment_(pull, comment)
}

func (s MockPullRequestsService) DeleteComment(pull PullRequestSpec, commentID int) (Response, error) {
	return s.DeleteComment_(pull, commentID)
}

func (s MockPullRequestsService) Merge(pull PullRequestSpec, mergeRequest *PullRequestMergeRequest) (*PullRequestMergeResult, Response, error) {
	return s.Merge_(pull, mergeRequest)
}
