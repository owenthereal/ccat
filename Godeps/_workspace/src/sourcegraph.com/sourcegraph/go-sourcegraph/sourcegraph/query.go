package sourcegraph

import (
	"bytes"
	"unicode/utf8"
)

// A RawQuery is a raw search query. To obtain the results for the
// query, it must be tokenized, parsed, resolved, planned, etc.
type RawQuery struct {
	// String is the raw query string from the client.
	String string

	// InsertionPoint is the 0-indexed character offset of the text
	// insertion cursor on the client.
	InsertionPoint int
}

// A ResolvedQuery is a query that has been parsed and resolved so
// that each token is given an unambiguous meaning.
type ResolvedQuery struct {
	// Tokens are resolved tokens, each of whose meaning is
	// unambiguous.
	Tokens []Token
}

// Join joins tokens to reconstruct a query string (that, when
// tokenized, would yield the same tokens). It returns the query and
// the insertion point (which is the position of the active token's
// last character, or the position after the last token's last
// character if there is no active token).
func Join(tokens []Token) RawQuery {
	ip := -1
	var buf bytes.Buffer
	for i, tok := range tokens {
		if i != 0 {
			buf.Write([]byte{' '})
		}
		buf.WriteString(tok.String())
	}
	if ip == -1 && len(tokens) > 0 {
		ip = utf8.RuneCount(buf.Bytes()) + 1
	}
	if ip == -1 {
		ip = 0
	}
	return RawQuery{String: buf.String(), InsertionPoint: ip}
}
