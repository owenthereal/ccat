// Package syntaxhighlight provides syntax highlighting for code. It currently
// uses a language-independent lexer and performs decently on JavaScript, Java,
// Ruby, Python, Go, and C.
package syntaxhighlight

import (
	"bytes"
	"html"
	"io"
	"strings"
	"text/scanner"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/sourcegraph/annotate"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/vcsstore/vcsclient"
)

// Kind represents a syntax highlighting kind (class) which will be assigned to tokens.
// A syntax highlighting scheme (style) maps text style properties to each token kind.
type Kind uint8

const (
	Whitespace Kind = iota
	String
	Keyword
	Comment
	Type
	Literal
	Punctuation
	Plaintext
	Tag
	HTMLTag
	HTMLAttrName
	HTMLAttrValue
	Decimal
)

//go:generate gostringer -type=Kind

type Printer interface {
	Print(w io.Writer, kind Kind, tokText string) error
}

type HTMLConfig struct {
	String        string
	Keyword       string
	Comment       string
	Type          string
	Literal       string
	Punctuation   string
	Plaintext     string
	Tag           string
	HTMLTag       string
	HTMLAttrName  string
	HTMLAttrValue string
	Decimal       string
}

type HTMLPrinter HTMLConfig

func (c HTMLConfig) class(kind Kind) string {
	switch kind {
	case String:
		return c.String
	case Keyword:
		return c.Keyword
	case Comment:
		return c.Comment
	case Type:
		return c.Type
	case Literal:
		return c.Literal
	case Punctuation:
		return c.Punctuation
	case Plaintext:
		return c.Plaintext
	case Tag:
		return c.Tag
	case HTMLTag:
		return c.HTMLTag
	case HTMLAttrName:
		return c.HTMLAttrName
	case HTMLAttrValue:
		return c.HTMLAttrValue
	case Decimal:
		return c.Decimal
	}
	return ""
}

func (p HTMLPrinter) Print(w io.Writer, kind Kind, tokText string) error {
	class := ((HTMLConfig)(p)).class(kind)
	if class != "" {
		_, err := w.Write([]byte(`<span class="`))
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, class)
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(`">`))
		if err != nil {
			return err
		}
	}
	template.HTMLEscape(w, []byte(tokText))
	if class != "" {
		_, err := w.Write([]byte(`</span>`))
		if err != nil {
			return err
		}
	}
	return nil
}

type Annotator interface {
	Annotate(start int, kind Kind, tokText string) (*annotate.Annotation, error)
}

// NilAnnotator is a special kind of annotator that always returns nil, but stores
// within itself the snippet of source code that is passed through it as tokens.
//
// This functionality is useful when one wishes to obtain the tokenized source as a data
// structure, as opposed to an annotated string, allowing full control over rendering and
// displaying it.
type NilAnnotator struct {
	Config     HTMLConfig
	Code       *sourcegraph.SourceCode
	byteOffset int
}

func NewNilAnnotator(e *vcsclient.FileWithRange) *NilAnnotator {
	ann := NilAnnotator{
		Config: DefaultHTMLConfig,
		Code: &sourcegraph.SourceCode{
			Lines: make([]*sourcegraph.SourceCodeLine, 0, bytes.Count(e.Contents, []byte("\n"))),
		},
		byteOffset: int(e.StartByte),
	}
	ann.addLine(ann.byteOffset)
	return &ann
}

func (a *NilAnnotator) addToken(t interface{}) {
	line := a.Code.Lines[len(a.Code.Lines)-1]
	if line.Tokens == nil {
		line.Tokens = make([]interface{}, 0, 1)
	}
	// If this token and the previous one are both strings, merge them.
	n := len(line.Tokens)
	if t1, ok := t.(string); ok && n > 0 {
		if t2, ok := (line.Tokens[n-1]).(string); ok {
			line.Tokens[n-1] = string(t1 + t2)
			return
		}
	}
	line.Tokens = append(line.Tokens, t)
}

func (a *NilAnnotator) addLine(startByte int) {
	a.Code.Lines = append(a.Code.Lines, &sourcegraph.SourceCodeLine{StartByte: startByte})
	if len(a.Code.Lines) > 1 {
		lastLine := a.Code.Lines[len(a.Code.Lines)-2]
		lastLine.EndByte = startByte - 1
	}
}

func (a *NilAnnotator) addMultilineToken(startByte int, unsafeHTML string, class string) {
	lines := strings.Split(unsafeHTML, "\n")
	for n, unsafeHTML := range lines {
		if len(unsafeHTML) > 0 {
			a.addToken(&sourcegraph.SourceCodeToken{
				StartByte: startByte,
				EndByte:   startByte + len(unsafeHTML),
				Class:     class,
				Label:     html.EscapeString(unsafeHTML),
			})
			startByte += len(unsafeHTML)
		}
		if n < len(lines)-1 {
			a.addLine(startByte)
		}
	}
}

func (a *NilAnnotator) Annotate(start int, kind Kind, tokText string) (*annotate.Annotation, error) {
	class := ((HTMLConfig)(a.Config)).class(kind)
	txt := html.EscapeString(tokText)
	start += a.byteOffset

	switch {
	// New line
	case tokText == "\n":
		a.addLine(start + 1)

	// Whitespace token
	case class == "":
		a.addToken(txt)

	// Multiline token (ie. block comments, string literals)
	case strings.Contains(tokText, "\n"):
		// Here we pass the unescaped string so we can calculate line lenghts correctly.
		// This method is expected to take responsibility of escaping any token text.
		a.addMultilineToken(start+1, tokText, class)

	// Token
	default:
		a.addToken(&sourcegraph.SourceCodeToken{
			StartByte: start,
			EndByte:   start + len(tokText),
			Class:     class,
			Label:     txt,
		})
	}

	return nil, nil
}

type HTMLAnnotator HTMLConfig

func (a HTMLAnnotator) Annotate(start int, kind Kind, tokText string) (*annotate.Annotation, error) {
	class := ((HTMLConfig)(a)).class(kind)
	if class != "" {
		left := []byte(`<span class="`)
		left = append(left, []byte(class)...)
		left = append(left, []byte(`">`)...)
		return &annotate.Annotation{
			Start: start, End: start + len(tokText),
			Left: left, Right: []byte("</span>"),
		}, nil
	}
	return nil, nil
}

// DefaultHTMLConfig's class names match those of google-code-prettify
// (https://code.google.com/p/google-code-prettify/).
var DefaultHTMLConfig = HTMLConfig{
	String:        "str",
	Keyword:       "kwd",
	Comment:       "com",
	Type:          "typ",
	Literal:       "lit",
	Punctuation:   "pun",
	Plaintext:     "pln",
	Tag:           "tag",
	HTMLTag:       "htm",
	HTMLAttrName:  "atn",
	HTMLAttrValue: "atv",
	Decimal:       "dec",
}

func Print(s *scanner.Scanner, w io.Writer, p Printer) error {
	tok := s.Scan()
	for tok != scanner.EOF {
		tokText := s.TokenText()
		err := p.Print(w, tokenKind(tok, tokText), tokText)
		if err != nil {
			return err
		}

		tok = s.Scan()
	}

	return nil
}

func Annotate(src []byte, a Annotator) (annotate.Annotations, error) {
	s := NewScanner(src)

	var anns annotate.Annotations
	read := 0

	tok := s.Scan()
	for tok != scanner.EOF {
		tokText := s.TokenText()

		ann, err := a.Annotate(read, tokenKind(tok, tokText), tokText)
		if err != nil {
			return nil, err
		}
		read += len(tokText)
		if ann != nil {
			anns = append(anns, ann)
		}

		tok = s.Scan()
	}

	return anns, nil
}

func AsHTML(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := Print(NewScanner(src), &buf, HTMLPrinter(DefaultHTMLConfig))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// NewScanner is a helper that takes a []byte src, wraps it in a reader and creates a Scanner.
func NewScanner(src []byte) *scanner.Scanner {
	return NewScannerReader(bytes.NewReader(src))
}

// NewScannerReader takes a reader src and creates a Scanner.
func NewScannerReader(src io.Reader) *scanner.Scanner {
	var s scanner.Scanner
	s.Init(src)
	s.Error = func(_ *scanner.Scanner, _ string) {}
	s.Whitespace = 0
	s.Mode = s.Mode ^ scanner.SkipComments
	return &s
}

func tokenKind(tok rune, tokText string) Kind {
	switch tok {
	case scanner.Ident:
		if _, isKW := keywords[tokText]; isKW {
			return Keyword
		}
		if r, _ := utf8.DecodeRuneInString(tokText); unicode.IsUpper(r) {
			return Type
		}
		return Plaintext
	case scanner.Float, scanner.Int:
		return Decimal
	case scanner.Char, scanner.String, scanner.RawString:
		return String
	case scanner.Comment:
		return Comment
	}
	if unicode.IsSpace(tok) {
		return Whitespace
	}
	return Punctuation
}
