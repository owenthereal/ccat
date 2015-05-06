package sourcegraph

type MockSearchService struct {
	Search_   func(opt *SearchOptions) (*SearchResults, Response, error)
	Complete_ func(q RawQuery) (*Completions, Response, error)
	Suggest_  func(q RawQuery) ([]*Suggestion, Response, error)
}

var _ SearchService = MockSearchService{}

func (s MockSearchService) Search(opt *SearchOptions) (*SearchResults, Response, error) {
	return s.Search_(opt)
}

func (s MockSearchService) Complete(q RawQuery) (*Completions, Response, error) {
	return s.Complete_(q)
}

func (s MockSearchService) Suggest(q RawQuery) ([]*Suggestion, Response, error) {
	return s.Suggest_(q)
}
