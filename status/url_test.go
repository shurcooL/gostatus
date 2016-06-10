package status

import (
	"reflect"
	"testing"
)

func TestEqualRepoURLs(t *testing.T) {
	tests := []struct {
		rawurl1 string
		rawurl2 string
		want    bool
	}{
		{
			rawurl1: "https://github.com/user/repo",
			rawurl2: "https://github.com/user/repo",
			want:    true,
		},
		{
			rawurl1: "https://github.com/user/repo",
			rawurl2: "https://github.com/user/repo.git",
			want:    true,
		},
		{
			rawurl1: "https://github.com/user/repo",
			rawurl2: "https://github.com/user/wrongrepo",
			want:    false,
		},
		{
			rawurl1: "https://github.com/user/repo",
			rawurl2: "git@github.com:user/repo",
			want:    true,
		},
		{
			rawurl1: "https://github.com/user/repo",
			rawurl2: "git@github.com:user/repo.git",
			want:    true,
		},
	}
	for _, test := range tests {
		if got, want := EqualRepoURLs(test.rawurl1, test.rawurl2), test.want; got != want {
			t.Errorf("EqualRepoURLs(%q, %q): got %v, want %v", test.rawurl1, test.rawurl2, got, want)
		}
	}
}

func TestParseURL(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr error
	}{
		{
			in:   "git@github.com:user/repo",
			want: "ssh://git@github.com/user/repo",
		},
		{
			in:   "git@github.com:user/repo.git",
			want: "ssh://git@github.com/user/repo.git",
		},
		{
			in:   "https://github.com/user/repo",
			want: "https://github.com/user/repo",
		},
		{
			in:   "https://github.com/user/repo.git",
			want: "https://github.com/user/repo.git",
		},
	}
	for _, test := range tests {
		u, err := parseURL(test.in)
		if got, want := err, test.wantErr; !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
		if test.wantErr != nil {
			continue
		}

		if got, want := u.String(), test.want; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}
