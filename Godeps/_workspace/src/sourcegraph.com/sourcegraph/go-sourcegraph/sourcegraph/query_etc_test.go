package sourcegraph

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestTokenError_JSON(t *testing.T) {
	rerr := []TokenError{
		TokenError{Message: "a"},
		TokenError{Token: Term("t"), Message: "a"},
		TokenError{Index: 1, Token: Term(""), Message: "a"},
		TokenError{Index: 2, Token: RepoToken{URI: "r"}, Message: "b"},
	}

	rerrJSON, err := json.MarshalIndent(rerr, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	want := `[
  {
    "Token": null,
    "Message": "a"
  },
  {
    "Token": {
      "String": "t",
      "Type": "Term"
    },
    "Message": "a"
  },
  {
    "Index": 1,
    "Token": {
      "String": "",
      "Type": "Term"
    },
    "Message": "a"
  },
  {
    "Index": 2,
    "Token": {
      "URI": "r",
      "Type": "RepoToken"
    },
    "Message": "b"
  }
]`
	if string(rerrJSON) != want {
		t.Errorf("got JSON\n%s\n\nwant JSON\n%s", rerrJSON, want)
	}

	var rerr2 []TokenError
	if err := json.Unmarshal(rerrJSON, &rerr2); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(rerr2, rerr) {
		t.Errorf("got\n%#v\n\nwant\n%#v", rerr2, rerr)
	}
}
