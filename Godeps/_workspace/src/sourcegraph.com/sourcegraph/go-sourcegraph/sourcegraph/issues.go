package sourcegraph

import (
	"fmt"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/sourcegraph/go-github/github"
	"strconv"
	"strings"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
)

// IssuesService communicates with the issue-related endpoints in the
// Sourcegraph API.
type IssuesService interface {
	// Get fetches a issue.
	Get(issue IssueSpec, opt *IssueGetOptions) (*Issue, Response, error)

	// List issues for a repository.
	ListByRepo(repo RepoSpec, opt *IssueListOptions) ([]*Issue, Response, error)

	// ListComments lists comments on a issue.
	ListComments(issue IssueSpec, opt *IssueListCommentsOptions) ([]*IssueComment, Response, error)

	// CreateComment creates a comment on an issue.
	CreateComment(issue IssueSpec, comment *IssueComment) (*IssueComment, Response, error)

	// EditComment updates a comment on an issue.
	EditComment(issue IssueSpec, comment *IssueComment) (*IssueComment, Response, error)

	// DeleteComment deletes a comment on an issue.
	DeleteComment(issue IssueSpec, commentID int) (Response, error)
}

// issuesService implements IssuesService.
type issuesService struct {
	client *Client
}

var _ IssuesService = &issuesService{}

// IssueSpec specifies a issue.
type IssueSpec struct {
	Repo RepoSpec

	Number int // Sequence number of the issue
}

func (s IssueSpec) RouteVars() map[string]string {
	return map[string]string{"RepoSpec": s.Repo.URI, "Issue": strconv.Itoa(s.Number)}
}

func UnmarshalIssueSpec(routeVars map[string]string) (IssueSpec, error) {
	issueNumber, err := strconv.Atoi(routeVars["Issue"])
	if err != nil {
		return IssueSpec{}, err
	}
	repoURI := routeVars["RepoSpec"]
	if repoURI == "" {
		return IssueSpec{}, fmt.Errorf("RepoSpec was empty")
	}
	return IssueSpec{
		Repo:   RepoSpec{URI: repoURI},
		Number: issueNumber,
	}, nil
}

// Issue is a issue returned by the Sourcegraph API.
type Issue struct {
	github.Issue
}

// Spec returns the IssueSpec that specifies r.
func (r *Issue) Spec() IssueSpec {
	// Extract the URI from the HTMLURL field.
	uri := strings.Join(strings.Split(strings.TrimPrefix(*r.HTMLURL, "https://"), "/")[0:3], "/")
	return IssueSpec{
		Repo:   RepoSpec{URI: uri},
		Number: *r.Number,
	}
}

type IssueGetOptions struct{}

func (s *issuesService) Get(issue IssueSpec, opt *IssueGetOptions) (*Issue, Response, error) {
	url, err := s.client.URL(router.RepoIssue, issue.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var issue_ *Issue
	resp, err := s.client.Do(req, &issue_)
	if err != nil {
		return nil, resp, err
	}

	return issue_, resp, nil
}

type IssueListOptions struct {
	State string `url:",omitempty"` // "open", "closed", or "all"
	ListOptions
}

func (s *issuesService) ListByRepo(repo RepoSpec, opt *IssueListOptions) ([]*Issue, Response, error) {
	url, err := s.client.URL(router.RepoIssues, repo.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var issues []*Issue
	resp, err := s.client.Do(req, &issues)
	if err != nil {
		return nil, resp, err
	}

	return issues, resp, nil
}

type IssueListCommentsOptions struct {
	ListOptions
}

type IssueComment struct {
	RenderedBody string // the body rendered to HTML (from raw markdown)
	Checklist    *Checklist
	github.IssueComment
}

type IssueCommentSpec struct {
	Issue   IssueSpec
	Comment int
}

func (s IssueCommentSpec) RouteVars() map[string]string {
	rv := s.Issue.RouteVars()
	rv["CommentID"] = strconv.Itoa(s.Comment)
	return rv
}

func (s *issuesService) ListComments(issue IssueSpec, opt *IssueListCommentsOptions) ([]*IssueComment, Response, error) {
	url, err := s.client.URL(router.RepoIssueComments, issue.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var comments []*IssueComment
	resp, err := s.client.Do(req, &comments)
	if err != nil {
		return nil, resp, err
	}

	return comments, resp, nil
}

func (s *issuesService) CreateComment(issue IssueSpec, comment *IssueComment) (*IssueComment, Response, error) {
	url, err := s.client.URL(router.RepoIssueCommentsCreate, issue.RouteVars(), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", url.String(), comment)
	if err != nil {
		return nil, nil, err
	}

	var createdComment IssueComment
	resp, err := s.client.Do(req, &createdComment)
	if err != nil {
		return nil, nil, err
	}

	return &createdComment, resp, nil
}

func (s *issuesService) EditComment(issue IssueSpec, comment *IssueComment) (*IssueComment, Response, error) {
	if comment.ID == nil {
		return nil, nil, fmt.Errorf("comment ID not specified")
	}

	url, err := s.client.URL(router.RepoIssueCommentsEdit, IssueCommentSpec{Issue: issue, Comment: *comment.ID}.RouteVars(), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("PATCH", url.String(), comment)
	if err != nil {
		return nil, nil, err
	}

	var updatedComment IssueComment
	resp, err := s.client.Do(req, &updatedComment)
	if err != nil {
		return nil, nil, err
	}

	return &updatedComment, resp, nil
}

func (s *issuesService) DeleteComment(issue IssueSpec, commentID int) (Response, error) {
	url, err := s.client.URL(router.RepoIssueCommentsDelete, IssueCommentSpec{Issue: issue, Comment: commentID}.RouteVars(), nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

var _ IssuesService = &MockIssuesService{}
