// Package status provides a func to check if two repo URLs are equal
// in the context of Go packages.
package status

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// EqualRepoURLs reports whether two URLs are equal, ignoring scheme and userinfo.
// It parses URLs with support for SCP-like syntax, like the cmd/go tool.
// If there are any errors parsing the URLs, it resorts to doing a string comparison.
func EqualRepoURLs(rawurl0, rawurl1 string) bool {
	u, _, err := parseURL(rawurl0)
	if err != nil {
		return rawurl0 == rawurl1
	}
	v, _, err := parseURL(rawurl1)
	if err != nil {
		return rawurl0 == rawurl1
	}
	u.Scheme, v.Scheme = "", "" // Ignore scheme.
	u.User, v.User = nil, nil   // Ignore username and password information.
	// Ignore the .git extension, which GitHub ignores for the git User-Agent.
	u.Path, v.Path = strings.TrimSuffix(u.Path, ".git"), strings.TrimSuffix(v.Path, ".git")
	return strings.ToLower(u.String()) == strings.ToLower(v.String())
}

// FormatRepoURL tries to rewrite rawurl to follow the same format as layout URL.
// If either of two URLs has parsing errors, then rawurl is returned unmodified.
func FormatRepoURL(layout, rawurl string) string {
	u, _, err := parseURL(rawurl)
	if err != nil {
		return rawurl
	}
	l, scpSyntax, err := parseURL(layout)
	if err != nil {
		return rawurl
	}
	u.Scheme = l.Scheme // Take scheme from layout.
	u.User = l.User     // Take username and password information from layout.
	if scpSyntax {
		return fmt.Sprintf("%s@%s:%s", u.User.Username(), u.Host, strings.TrimPrefix(u.Path, "/"))
	}
	return u.String()
}

// scpSyntaxRE matches the SCP-like addresses used by Git to access repositories by SSH.
var scpSyntaxRE = regexp.MustCompile(`^([a-zA-Z0-9_]+)@([a-zA-Z0-9._-]+):(.*)$`)

// parseURL is like url.Parse but with support for SCP-like syntax.
func parseURL(rawurl string) (_ *url.URL, scpSyntax bool, _ error) {
	// Match SCP-like syntax and convert it to a URL.
	if m := scpSyntaxRE.FindStringSubmatch(rawurl); m != nil {
		// E.g., "git@github.com:user/repo" becomes "ssh://git@github.com/user/repo".
		return &url.URL{
			Scheme: "ssh",
			User:   url.User(m[1]),
			Host:   m[2],
			Path:   "/" + m[3],
		}, true, nil
	}

	u, err := url.Parse(rawurl)
	return u, false, err
}
