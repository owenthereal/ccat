package graph

import "testing"

func TestMakeURI(t *testing.T) {
	tests := []struct {
		cloneURL string
		want     string
	}{
		{"https://github.com/user/repo", "github.com/user/repo"},
		{"git://github.com/user/repo", "github.com/user/repo"},
		{"http://bitbucket.org/user/repo", "bitbucket.org/user/repo"},
		{"https://bitbucket.org/user/repo", "bitbucket.org/user/repo"},
		{"bitbucket.org/user/repo", "bitbucket.org/user/repo"},
	}

	for _, test := range tests {
		got := MakeURI(test.cloneURL)
		if test.want != got {
			t.Errorf("%s: want URI %s, got %s", test.cloneURL, test.want, got)
		}
	}
}
