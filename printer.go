package main

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/sourcegraph/syntaxhighlight"
)

var (
	stringKind        = kind{"String", syntaxhighlight.String}
	keywordKind       = kind{"Keyword", syntaxhighlight.Keyword}
	commentKind       = kind{"Comment", syntaxhighlight.Comment}
	typeKind          = kind{"Type", syntaxhighlight.Type}
	literalKind       = kind{"Literal", syntaxhighlight.Literal}
	punctuationKind   = kind{"Punctuation", syntaxhighlight.Punctuation}
	plaintextKind     = kind{"Plaintext", syntaxhighlight.Plaintext}
	tagKind           = kind{"Tag", syntaxhighlight.Tag}
	htmlTagKind       = kind{"HTMLTag", syntaxhighlight.HTMLTag}
	htmlAttrNameKind  = kind{"HTMLAttrName", syntaxhighlight.HTMLAttrName}
	htmlAttrValueKind = kind{"HTMLAttrValue", syntaxhighlight.HTMLAttrValue}
	decimalKind       = kind{"Decimal", syntaxhighlight.Decimal}

	kinds = []kind{
		stringKind,
		keywordKind,
		commentKind,
		typeKind,
		literalKind,
		punctuationKind,
		plaintextKind,
		tagKind,
		htmlTagKind,
		htmlAttrNameKind,
		htmlAttrValueKind,
		decimalKind,
	}

	LightColorPalettes = ColorPalettes{
		stringKind:        "brown",
		keywordKind:       "darkblue",
		commentKind:       "lightgrey",
		typeKind:          "teal",
		literalKind:       "teal",
		punctuationKind:   "darkred",
		plaintextKind:     "darkblue",
		tagKind:           "blue",
		htmlTagKind:       "lightgreen",
		htmlAttrNameKind:  "blue",
		htmlAttrValueKind: "green",
		decimalKind:       "darkblue",
	}

	DarkColorPalettes = ColorPalettes{
		stringKind:        "brown",
		keywordKind:       "blue",
		commentKind:       "darkgrey",
		typeKind:          "turquoise",
		literalKind:       "turquoise",
		punctuationKind:   "red",
		plaintextKind:     "blue",
		tagKind:           "blue",
		htmlTagKind:       "lightgreen",
		htmlAttrNameKind:  "blue",
		htmlAttrValueKind: "green",
		decimalKind:       "blue",
	}

	// cache kind name and syntax highlight kind
	// for faster lookup
	kindsByName map[string]kind
	kindsByKind map[syntaxhighlight.Kind]kind
)

func init() {
	kindsByName = make(map[string]kind)
	for _, k := range kinds {
		kindsByName[k.Name] = k
	}

	kindsByKind = make(map[syntaxhighlight.Kind]kind)
	for _, k := range kinds {
		kindsByKind[k.Kind] = k
	}
}

type kind struct {
	Name string
	Kind syntaxhighlight.Kind
}

type ColorPalettes map[kind]string

func (c ColorPalettes) Set(k, v string) bool {
	kind, ok := kindsByName[k]
	if ok {
		c[kind] = v
	}

	return ok
}

func (c ColorPalettes) Get(k syntaxhighlight.Kind) string {
	// ignore whitespace kind
	if k == syntaxhighlight.Whitespace {
		return ""
	}

	kind, ok := kindsByKind[k]
	if !ok {
		panic(fmt.Sprintf("Unknown syntax highlight kind %d\n", k))
	}

	return c[kind]
}

func (c ColorPalettes) String() string {
	var s []string
	for _, k := range kinds {
		color := c[k]
		s = append(s, fmt.Sprintf("%13s\t%s", k.Name, Colorize(color, color)))
	}

	return strings.Join(s, "\n")
}

func CPrint(r io.Reader, w io.Writer, palettes ColorPalettes) error {
	return syntaxhighlight.Print(syntaxhighlight.NewScannerReader(r), w, Printer{palettes})
}

type Printer struct {
	ColorPalettes ColorPalettes
}

func (p Printer) Print(w io.Writer, kind syntaxhighlight.Kind, tokText string) error {
	c := p.ColorPalettes.Get(kind)
	if len(c) > 0 {
		tokText = Colorize(c, tokText)
	}

	_, err := io.WriteString(w, tokText)

	return err
}

func HtmlPrint(r io.Reader, w io.Writer, palettes ColorPalettes) error {
	keys := []string{}
	for k := range htmlCodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	w.Write([]byte("<style>\n"))
	for _, s := range keys {
		if s == "" {
			continue
		}
		w.Write([]byte(fmt.Sprintf(".%s { color: %s; }\n", s, s)))
	}
	w.Write([]byte("</style>\n"))
	w.Write([]byte("<pre>\n"))
	err := syntaxhighlight.Print(syntaxhighlight.NewScannerReader(r), w, HtmlCodePrinter{palettes})
	w.Write([]byte("\n</pre>\n"))
	return err
}

type HtmlCodePrinter struct {
	ColorPalettes ColorPalettes
}

func (p HtmlCodePrinter) Print(w io.Writer, kind syntaxhighlight.Kind, tokText string) error {
	c := p.ColorPalettes.Get(kind)
	if len(c) > 0 {
		tokText = Htmlize(c, tokText)
	}

	_, err := io.WriteString(w, tokText)

	return err
}
