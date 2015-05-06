package store

import (
	"reflect"
	"testing"
)

func TestAncestorDirsExceptRoot(t *testing.T) {
	tests := map[string][]string{
		".":     nil,
		"/":     nil,
		"":      nil,
		"a":     nil,
		"a/":    []string{"a"}, // maybe we don't want this behavior
		"a/b":   []string{"a"},
		"a/b/c": []string{"a", "a/b"},
	}
	for p, want := range tests {
		dirs := ancestorDirsExceptRoot(p)
		if !reflect.DeepEqual(dirs, want) {
			t.Errorf("%v: got %v, want %v", p, dirs, want)
		}
	}
}
