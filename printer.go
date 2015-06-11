package main

import (
	"io"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/sourcegraph/syntaxhighlight"
)

type ColorDefs map[syntaxhighlight.Kind]string

func (c ColorDefs) Set(k, v string) bool {
	ok := true
	switch k {
	case "String":
		c[syntaxhighlight.String] = v
	case "Keyword":
		c[syntaxhighlight.Keyword] = v
	case "Comment":
		c[syntaxhighlight.Comment] = v
	case "Type":
		c[syntaxhighlight.Type] = v
	case "Literal":
		c[syntaxhighlight.Literal] = v
	case "Punctuation":
		c[syntaxhighlight.Punctuation] = v
	case "Plaintext":
		c[syntaxhighlight.Plaintext] = v
	case "Tag":
		c[syntaxhighlight.Tag] = v
	case "HTMLTag":
		c[syntaxhighlight.HTMLTag] = v
	case "HTMLAttrName":
		c[syntaxhighlight.HTMLAttrName] = v
	case "HTMLAttrValue":
		c[syntaxhighlight.HTMLAttrValue] = v
	case "Decimal":
		c[syntaxhighlight.Decimal] = v
	default:
		ok = false
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
