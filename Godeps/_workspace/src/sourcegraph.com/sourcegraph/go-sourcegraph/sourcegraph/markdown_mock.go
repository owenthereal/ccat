package sourcegraph

type MockMarkdownService struct {
	Render_ func(markdown []byte, opt MarkdownOpt) (*MarkdownData, Response, error)
}

func (s MockMarkdownService) Render(markdown []byte, opt MarkdownOpt) (*MarkdownData, Response, error) {
	return s.Render_(markdown, opt)
}
