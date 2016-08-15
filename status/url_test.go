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

func TestFormatRepoURL(t *testing.T) {
	tests := []struct {
		layout string
		rawurl string
		want   string
	}{
		{
			layout: "https://github.com/user/repo",
			rawurl: "https://github.com/user/repo",
			want:   "https://github.com/user/repo",
		},
		{
			// .git suffix does not have special treatment, and these are not considered as equal URLs.
			// See https://github.com/shurcooL/gostatus/issues/37#issuecomment-225336823 for rationale.
			layout: "git@github.com:user/repo.git",
			rawurl: "https://github.com/user/repo",
			want:   "git@github.com:user/repo",
		},
		{
			layout: "https://github.com/user/wrongrepo",
			rawurl: "https://github.com/user/repo",
			want:   "https://github.com/user/repo",
		},
		{
			layout: "git@github.com:user/repo",
			rawurl: "https://github.com/user/repo",
			want:   "git@github.com:user/repo",
		},
	}
	for _, test := range tests {
		if got, want := FormatRepoURL(test.layout, test.rawurl), test.want; got != want {
			t.Errorf("FormatRepoURL(%q, %q): got %q, want %q", test.layout, test.rawurl, got, want)
		}
	}
}

func TestParseURL(t *testing.T) {
	tests := []struct {
		in            string
		want          string
		wantSCPSyntax bool
		wantErr       error
	}{
		{
			in:            "git@github.com:user/repo",
			want:          "ssh://git@github.com/user/repo",
			wantSCPSyntax: true,
		},
		{
			in:            "https://github.com/user/repo",
			want:          "https://github.com/user/repo",
			wantSCPSyntax: false,
		},
	}
	for _, test := range tests {
		u, scpSyntax, err := parseURL(test.in)
		if got, want := err, test.wantErr; !reflect.DeepEqual(got, want) {
			t.Errorf("err: got %q, want %q", got, want)
		}
		if test.wantErr != nil {
			continue
		}

		if got, want := u.String(), test.want; got != want {
			t.Errorf("url: got %q, want %q", got, want)
		}
		if got, want := scpSyntax, test.wantSCPSyntax; got != want {
			t.Errorf("scpSyntax: got %v, want %v", got, want)
		}
	}
}
