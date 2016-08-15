package status

import (
	"reflect"
	"testing"
)

func TestEqualRepoURLs(t *testing.T) {
	tests := []struct {
		rawurl0 string
		rawurl1 string
		want    bool
	}{
		{
			rawurl0: "https://github.com/user/repo",
			rawurl1: "https://github.com/user/repo",
			want:    true,
		},
		{
			// .git suffix does not have special treatment, and these are not considered as equal URLs.
			// See https://github.com/shurcooL/gostatus/issues/37#issuecomment-225336823 for rationale.
			rawurl0: "https://github.com/user/repo",
			rawurl1: "https://github.com/user/repo.git",
			want:    false,
		},
		{
			rawurl0: "https://github.com/user/repo",
			rawurl1: "https://github.com/user/wrongrepo",
			want:    false,
		},
		{
			rawurl0: "https://github.com/user/repo",
			rawurl1: "git@github.com:user/repo",
			want:    true,
		},
	}
	for _, test := range tests {
		if got, want := EqualRepoURLs(test.rawurl0, test.rawurl1), test.want; got != want {
			t.Errorf("EqualRepoURLs(%q, %q): got %v, want %v", test.rawurl0, test.rawurl1, got, want)
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
			in:   "https://github.com/user/repo",
			want: "https://github.com/user/repo",
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
