package sourcegraph

import (
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-vcs/vcs"
)

// SearchService communicates with the search-related endpoints in
// the Sourcegraph API.
type SearchService interface {
	// Search searches the full index.
	Search(opt *SearchOptions) (*SearchResults, Response, error)

	// Complete completes the token at the RawQuery's InsertionPoint.
	Complete(q RawQuery) (*Completions, Response, error)

	// Suggest suggests queries given an existing query. It can be
	// called with an empty query to get example queries that pertain
	// to the current user's repositories, orgs, etc.
	Suggest(q RawQuery) ([]*Suggestion, Response, error)
}

type SearchResults struct {
	Defs   []*Def                  `json:",omitempty"`
	People []*Person               `json:",omitempty"`
	Repos  []*Repo                 `json:",omitempty"`
	Tree   []*RepoTreeSearchResult `json:",omitempty"`

	// RawQuery is the raw query passed to search.
	RawQuery RawQuery

	// Tokens are the unresolved tokens.
	Tokens Tokens `json:",omitempty"`

	// Plan is the query plan used to fetch the results.
	Plan *Plan `json:",omitempty"`

	// ResolvedTokens holds the resolved tokens from the original query
	// string.
	ResolvedTokens Tokens

	ResolveErrors []TokenError `json:",omitempty"`

	// Tips are helpful tips for the user about their query. They are
	// not errors per se, but they use the TokenError type because it
	// allows us to associate a message with a particular token (and
	// JSON de/serialize that).
	Tips []TokenError `json:",omitempty"`

	// Canceled is true if the query was canceled. More information
	// about how to correct the issue can be found in the
	// ResolveErrors and Tips.
	Canceled bool
}

// Empty is whether there are no search results for any result type.
func (r *SearchResults) Empty() bool {
	return len(r.Defs) == 0 && len(r.People) == 0 && len(r.Repos) == 0 && len(r.Tree) == 0
}

// A RepoTreeSearchResult is a tree search result that includes the repo
// and rev it came from.
type RepoTreeSearchResult struct {
	vcs.SearchResult
	RepoRev RepoRevSpec
}

// searchService implements SearchService.
type searchService struct {
	client *Client
}

var _ SearchService = &searchService{}

type SearchOptions struct {
	Query string `url:"q" schema:"q"`

	Defs   bool
	Repos  bool
	People bool
	Tree   bool

	ListOptions
}

func (s *searchService) Search(opt *SearchOptions) (*SearchResults, Response, error) {
	url, err := s.client.URL(router.Search, nil, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var results *SearchResults
	resp, err := s.client.Do(req, &results)
	if err != nil {
		return nil, resp, err
	}

	return results, resp, nil
}

// Completions holds search query completions.
type Completions struct {
	// TokenCompletions are suggested completions for the token at the
	// raw query's InsertionPoint.
	TokenCompletions Tokens

	// ResolvedTokens is the resolution of the original query's tokens
	// used to produce the completions. It is useful for debugging.
	ResolvedTokens Tokens

	ResolveErrors   []TokenError `json:",omitempty"`
	ResolutionFatal bool         `json:",omitempty"`
}

func (s *searchService) Complete(q RawQuery) (*Completions, Response, error) {
	url, err := s.client.URL(router.SearchComplete, nil, q)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var comps *Completions
	resp, err := s.client.Do(req, &comps)
	if err != nil {
		return nil, resp, err
	}

	return comps, resp, nil
}

func (s *searchService) Suggest(q RawQuery) ([]*Suggestion, Response, error) {
	url, err := s.client.URL(router.SearchSuggestions, nil, q)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var suggs []*Suggestion
	resp, err := s.client.Do(req, &suggs)
	if err != nil {
		return nil, resp, err
	}

	return suggs, resp, nil
}
