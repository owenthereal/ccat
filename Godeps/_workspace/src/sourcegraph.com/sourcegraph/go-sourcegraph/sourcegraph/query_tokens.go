package sourcegraph

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/unit"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/vcsstore/vcsclient"
)

// A Token is the smallest indivisible component of a query, either a
// term or a "field:val" specifier (e.g., "repo:example.com/myrepo").
type Token interface {
	// String returns the string representation of the term.
	String() string
}

// A Term is a query term token. It is either a word or an arbitrary
// string (if quoted in the raw query).
type Term string

func (t Term) String() string {
	if strings.Contains(string(t), " ") {
		return `"` + string(t) + `"`
	}
	return string(t)
}

func (t Term) UnquotedString() string { return string(t) }

// An AnyToken is a token that has not yet been resolved into another
// token type. It resolves to Term if it can't be resolved to another
// token type.
type AnyToken string

func (u AnyToken) String() string { return string(u) }

// A RepoToken represents a repository, although it does not
// necessarily uniquely identify the repository. It consists of any
// number of slash-separated path components, such as "a/b" or
// "github.com/foo/bar".
type RepoToken struct {
	URI string

	Repo *Repo `json:",omitempty"`
}

func (t RepoToken) String() string { return t.URI }

func (t RepoToken) Spec() RepoSpec {
	var rid int
	if t.Repo != nil {
		rid = t.Repo.RID
	}
	return RepoSpec{URI: t.URI, RID: rid}
}

// A RevToken represents a specific revision (either a revspec or a
// commit ID) of a repository (which must be specified by a previous
// RepoToken in the query).
type RevToken struct {
	Rev string // Rev is either a revspec or commit ID

	Commit *Commit `json:",omitempty"`
}

func (t RevToken) String() string { return ":" + t.Rev }

// A UnitToken represents a source unit in a repository.
type UnitToken struct {
	// UnitType is the type of the source unit (e.g., GoPackage).
	UnitType string

	// Name is the name of the source unit (e.g., mypkg).
	Name string

	// Unit is the source unit object.
	Unit *unit.RepoSourceUnit
}

func (t UnitToken) String() string {
	s := "~" + t.Name
	if t.UnitType != "" {
		s += "@" + t.UnitType
	}
	return s
}

type FileToken struct {
	Path string

	Entry *vcsclient.TreeEntry
}

func (t FileToken) String() string { return "/" + filepath.Clean(t.Path) }

// A UserToken represents a user or org, although it does not
// necessarily uniquely identify one. It consists of the string "@"
// followed by a full or partial user/org login.
type UserToken struct {
	Login string

	User *User `json:",omitempty"`
}

func (t UserToken) String() string { return "@" + t.Login }

// Tokens wraps a list of tokens and adds some helper methods. It also
// serializes to JSON with "Type" fields added to each token and
// deserializes that same JSON back into a typed list of tokens.
type Tokens []Token

func (d Tokens) MarshalJSON() ([]byte, error) {
	jtoks := make([]jsonToken, len(d))
	for i, t := range d {
		jtoks[i] = jsonToken{t}
	}
	return json.Marshal(jtoks)
}

func (d *Tokens) UnmarshalJSON(b []byte) error {
	var jtoks []jsonToken
	if err := json.Unmarshal(b, &jtoks); err != nil {
		return err
	}
	if jtoks == nil {
		*d = nil
	} else {
		*d = make(Tokens, len(jtoks))
		for i, jtok := range jtoks {
			(*d)[i] = jtok.Token
		}
	}
	return nil
}

func (d Tokens) RawQueryString() string { return Join(d).String }

type jsonToken struct {
	Token `json:",omitempty"`
}

func (t jsonToken) MarshalJSON() ([]byte, error) {
	if t.Token == nil {
		return []byte("null"), nil
	}
	tokType := TokenType(t.Token)
	b, err := json.Marshal(t.Token)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return b, nil
	}
	switch b[0] {
	case '"':
		b = []byte(fmt.Sprintf(`{"String":%s,"Type":%q}`, b, tokType))
	case '{':
		b[len(b)-1] = ','
		b = append(b, []byte(fmt.Sprintf(`"Type":%q}`, tokType))...)
	}
	return b, nil
}

func (t *jsonToken) UnmarshalJSON(b []byte) error {
	var typ struct{ Type string }
	if err := json.Unmarshal(b, &typ); err != nil {
		return err
	}
	if typ.Type == "" {
		return nil
	}

	*t = jsonToken{}
	switch typ.Type {
	case "":
		return nil
	case "Term", "AnyToken":
		var tmp struct{ String string }
		if err := json.Unmarshal(b, &tmp); err != nil {
			return err
		}
		switch typ.Type {
		case "Term":
			t.Token = Term(tmp.String)
		case "AnyToken":
			t.Token = AnyToken(tmp.String)
		}
		return nil
	case "RepoToken":
		t.Token = &RepoToken{}
	case "RevToken":
		t.Token = &RevToken{}
	case "UnitToken":
		t.Token = &UnitToken{}
	case "FileToken":
		t.Token = &FileToken{}
	case "UserToken":
		t.Token = &UserToken{}
	default:
		return fmt.Errorf("unmarshal Tokens: unrecognized Type %q", typ.Type)
	}
	if err := json.Unmarshal(b, t.Token); err != nil {
		return err
	}
	t.Token = reflect.ValueOf(t.Token).Elem().Interface().(Token)
	return nil
}

func TokenType(tok Token) string {
	return strings.Replace(strings.Replace(reflect.ValueOf(tok).Type().String(), "*", "", -1), "sourcegraph.", "", -1)
}
