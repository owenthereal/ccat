package sourcegraph

type MockIssuesService struct {
	Get_           func(issue IssueSpec, opt *IssueGetOptions) (*Issue, Response, error)
	ListByRepo_    func(repo RepoSpec, opt *IssueListOptions) ([]*Issue, Response, error)
	ListComments_  func(issue IssueSpec, opt *IssueListCommentsOptions) ([]*IssueComment, Response, error)
	CreateComment_ func(issue IssueSpec, comment *IssueComment) (*IssueComment, Response, error)
	EditComment_   func(issue IssueSpec, comment *IssueComment) (*IssueComment, Response, error)
	DeleteComment_ func(issue IssueSpec, commentID int) (Response, error)
}

func (s MockIssuesService) Get(issue IssueSpec, opt *IssueGetOptions) (*Issue, Response, error) {
	return s.Get_(issue, opt)
}

func (s MockIssuesService) ListByRepo(repo RepoSpec, opt *IssueListOptions) ([]*Issue, Response, error) {
	return s.ListByRepo_(repo, opt)
}

func (s MockIssuesService) ListComments(issue IssueSpec, opt *IssueListCommentsOptions) ([]*IssueComment, Response, error) {
	return s.ListComments_(issue, opt)
}

func (s MockIssuesService) CreateComment(issue IssueSpec, comment *IssueComment) (*IssueComment, Response, error) {
	return s.CreateComment_(issue, comment)
}

func (s MockIssuesService) EditComment(issue IssueSpec, comment *IssueComment) (*IssueComment, Response, error) {
	return s.EditComment_(issue, comment)
}

func (s MockIssuesService) DeleteComment(issue IssueSpec, commentID int) (Response, error) {
	return s.DeleteComment_(issue, commentID)
}
