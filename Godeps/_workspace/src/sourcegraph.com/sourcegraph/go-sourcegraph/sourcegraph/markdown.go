package sourcegraph

import "github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"

type MarkdownService interface {
	Render(markdown []byte, opt MarkdownOpt) (*MarkdownData, Response, error)
}

type markdownService struct {
	client *Client
}

type MarkdownRequestBody struct {
	Markdown []byte
	MarkdownOpt
}

type MarkdownOpt struct {
	EnableCheckboxes bool
}

type MarkdownData struct {
	Rendered  []byte
	Checklist *Checklist
}

func (m *markdownService) Render(markdown []byte, opt MarkdownOpt) (*MarkdownData, Response, error) {
	url, err := m.client.URL(router.Markdown, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := m.client.NewRequest("POST", url.String(), &MarkdownRequestBody{
		Markdown:    markdown,
		MarkdownOpt: opt,
	})
	if err != nil {
		return nil, nil, err
	}

	var out MarkdownData
	resp, err := m.client.Do(req, &out)
	if err != nil {
		return nil, resp, err
	}

	return &out, resp, nil
}
