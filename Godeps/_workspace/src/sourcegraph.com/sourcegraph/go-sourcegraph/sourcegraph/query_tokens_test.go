package sourcegraph

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
	"time"
)

func TestTokens_JSON(t *testing.T) {
	tokens := Tokens{
		AnyToken("a"),
		Term("b"),
		Term(""),
		RepoToken{URI: "r", Repo: &Repo{RID: math.MaxInt32 - 1}},
		RevToken{Rev: "v"},
		FileToken{Path: "p"},
		UserToken{Login: "u"},
	}

	b, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	wantJSON := `[
  {
    "String": "a",
    "Type": "AnyToken"
  },
  {
    "String": "b",
    "Type": "Term"
  },
  {
    "String": "",
    "Type": "Term"
  },
  {
    "URI": "r",
    "Repo": {
      "RID": 2147483646,
      "URI": "",
      "URIAlias": null,
      "Name": "",
      "OwnerUserID": 0,
      "VCS": "",
      "HTTPCloneURL": "",
      "SSHCloneURL": null,
      "HomepageURL": null,
      "DefaultBranch": "",
      "Language": "",
      "GitHubStars": 0,
      "Deprecated": false,
      "Fork": false,
      "Mirror": false,
      "Private": false,
      "CreatedAt": "0001-01-01T00:00:00Z",
      "UpdatedAt": "0001-01-01T00:00:00Z",
      "PushedAt": "0001-01-01T00:00:00Z",
      "Permissions": {
        "Read": false,
        "Write": false,
        "Admin": false
      }
    },
    "Type": "RepoToken"
  },
  {
    "Rev": "v",
    "Type": "RevToken"
  },
  {
    "Path": "p",
    "Entry": null,
    "Type": "FileToken"
  },
  {
    "Login": "u",
    "Type": "UserToken"
  }
]`
	if string(b) != wantJSON {
		t.Errorf("got JSON\n%s\n\nwant JSON\n%s", b, wantJSON)
	}

	var tokens2 Tokens
	if err := json.Unmarshal(b, &tokens2); err != nil {
		t.Fatal(err)
	}

	// Normalize tokens for comparison.
	normTok := func(toks ...Token) {
		for _, tok := range toks {
			if tok, ok := tok.(RepoToken); ok {
				tok.Repo.CreatedAt = time.Time{}
				tok.Repo.UpdatedAt = time.Time{}
				tok.Repo.PushedAt = time.Time{}
			}
		}
	}
	normTok(tokens2...)
	normTok(tokens...)

	if !reflect.DeepEqual(tokens2, tokens) {
		t.Errorf("got tokens\n%+v\n\nwant tokens\n%+v", tokens2, tokens)
	}
}

func TestTokens_nil(t *testing.T) {
	tokens := Tokens(nil)

	b, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	wantJSON := `[]`
	if string(b) != wantJSON {
		t.Errorf("got JSON\n%s\n\nwant JSON\n%s", b, wantJSON)
	}
}
