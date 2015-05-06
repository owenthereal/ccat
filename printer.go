package main

import (
	"bytes"
	"io"

	"github.com/sourcegraph/syntaxhighlight"
)

type ColorDefs map[syntaxhighlight.Kind]string

var LightColorDefs = ColorDefs{
	syntaxhighlight.String:        "brown",
	syntaxhighlight.Keyword:       "darkblue",
	syntaxhighlight.Comment:       "lightgrey",
	syntaxhighlight.Type:          "teal",
	syntaxhighlight.Literal:       "teal",
	syntaxhighlight.Punctuation:   "darkred",
	syntaxhighlight.Plaintext:     "darkblue",
	syntaxhighlight.Tag:           "blue",
	syntaxhighlight.HTMLTag:       "lightgreen",
	syntaxhighlight.HTMLAttrName:  "blue",
	syntaxhighlight.HTMLAttrValue: "green",
	syntaxhighlight.Decimal:       "darkblue",
}

var DarkColorDefs = ColorDefs{
	syntaxhighlight.String:        "brown",
	syntaxhighlight.Keyword:       "blue",
	syntaxhighlight.Comment:       "darkgrey",
	syntaxhighlight.Type:          "turquoise",
	syntaxhighlight.Literal:       "turquoise",
	syntaxhighlight.Punctuation:   "red",
	syntaxhighlight.Plaintext:     "blue",
	syntaxhighlight.Tag:           "blue",
	syntaxhighlight.HTMLTag:       "lightgreen",
	syntaxhighlight.HTMLAttrName:  "blue",
	syntaxhighlight.HTMLAttrValue: "green",
	syntaxhighlight.Decimal:       "blue",
}

func AsCCat(src []byte, cdefs ColorDefs) ([]byte, error) {
	var buf bytes.Buffer
	err := syntaxhighlight.Print(syntaxhighlight.NewScanner(src), &buf, Printer{cdefs})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Printer struct {
	ColorDefs ColorDefs
}

func (p Printer) Print(w io.Writer, kind syntaxhighlight.Kind, tokText string) error {
	c := p.ColorDefs[kind]
	if c != "" {
		tokText = Colorize(c, tokText)
	}

	_, err := io.WriteString(w, tokText)
	if err != nil {
		return err
	}

	return nil
}
