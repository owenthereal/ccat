package main

import (
	"io"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/sourcegraph/syntaxhighlight"
)

var colorDefsMap = map[string]syntaxhighlight.Kind{
	"String":        syntaxhighlight.String,
	"Keyword":       syntaxhighlight.Keyword,
	"Comment":       syntaxhighlight.Comment,
	"Type":          syntaxhighlight.Type,
	"Literal":       syntaxhighlight.Literal,
	"Punctuation":   syntaxhighlight.Punctuation,
	"Plaintext":     syntaxhighlight.Plaintext,
	"Tag":           syntaxhighlight.Tag,
	"HTMLTag":       syntaxhighlight.HTMLTag,
	"HTMLAttrName":  syntaxhighlight.HTMLAttrName,
	"HTMLAttrValue": syntaxhighlight.HTMLAttrValue,
	"Decimal":       syntaxhighlight.Decimal,
}

type ColorDefs map[syntaxhighlight.Kind]string

func (c ColorDefs) Set(k, v string) bool {
	kind, ok := colorDefsMap[k]
	if ok {
		c[kind] = v
	}

	return ok
}

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

func CPrint(r io.Reader, w io.Writer, cdefs ColorDefs) error {
	return syntaxhighlight.Print(syntaxhighlight.NewScannerReader(r), w, Printer{cdefs})
}

type Printer struct {
	ColorDefs ColorDefs
}

func (p Printer) Print(w io.Writer, kind syntaxhighlight.Kind, tokText string) error {
	c, exists := p.ColorDefs[kind]
	if exists {
		tokText = Colorize(c, tokText)
	}

	_, err := io.WriteString(w, tokText)

	return err
}
