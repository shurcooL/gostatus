package main

import (
	"net/url"
	"regexp"
)

// equalRepoURLs reports whether two urls are equal, ignoring scheme and userinfo.
// If there are any errors parsing the urls, it resorts to doing a string comparison.
func equalRepoURLs(rawurl1, rawurl2 string) bool {
	u, err := parseURL(rawurl1)
	if err != nil {
		return rawurl1 == rawurl2
	}
	v, err := parseURL(rawurl2)
	if err != nil {
		return rawurl1 == rawurl2
	}
	u.Scheme, v.Scheme = "", "" // Ignore scheme.
	u.User, v.User = nil, nil   // Ignore username and password information.
	return u.String() == v.String()
}

// scpSyntaxRe matches the SCP-like addresses used by Git to access repositories by SSH.
var scpSyntaxRe = regexp.MustCompile(`^([a-zA-Z0-9_]+)@([a-zA-Z0-9._-]+):(.*)$`)

// parseURL is like url.Parse but with support for SCP-like syntax.
func parseURL(rawurl string) (*url.URL, error) {
	// Match SCP-like syntax and convert it to a URL.
	if m := scpSyntaxRe.FindStringSubmatch(rawurl); m != nil {
		// E.g., "git@github.com:user/repo" becomes "ssh://git@github.com/user/repo".
		return &url.URL{
			Scheme: "ssh",
			User:   url.User(m[1]),
			Host:   m[2],
			Path:   m[3],
		}, nil
	}

	return url.Parse(rawurl)
}
