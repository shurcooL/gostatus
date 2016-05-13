package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RepoFilter is a repo filter.
type RepoFilter func(r *Repo) (show bool)

// RepoPresenter is a repo presenter.
// All implementations must be read-only and safe for concurrent execution.
type RepoPresenter func(r *Repo) string

func trimPrefix(s, prefix string) (string, bool) {
	n := len(s)
	s = strings.TrimPrefix(s, prefix)
	return s, n != len(s)
}

func trimSuffix(s, suffix string) (string, bool) {
	n := len(s)
	s = strings.TrimSuffix(s, suffix)
	return s, n != len(s)
}

// Given url in the form "https://github.com/foo/bar", returns
// ("github.com", "foo/bar").
// If url doesn't match that format, returns ("", "")
func httpsToCanonical(url string) (string, string) {
	url, ok := trimPrefix(url, "https://")
	if !ok {
		return "", ""
	}
	parts := strings.Split(url, "/")
	if len(parts) != 3 {
		return "", ""
	}
	return parts[0], parts[1] + "/" + parts[2]
}

// Given url in the form "git@github.com:foo/bar.git", returns
// ("github.com", "foo/bar").
// If url doesn't match that format, returns ("", "")
func gitToCanonical(url string) (string, string) {
	url, ok := trimPrefix(url, "git@")
	url, ok2 := trimSuffix(url, ".git")
	if !ok || !ok2 {
		return "", ""
	}
	parts := strings.Split(url, ":")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func toCanonical(url string) (string, string) {
	host, repo := httpsToCanonical(url)
	if host != "" {
		return host, repo
	}
	return gitToCanonical(url)
}

// Heuristic to check if 2 urls represent the same repo.
// Currently only smart enough to consider https://github.com/foo/bar
// and git@github.com:foo/bar.git to be the same
func sameRepoURL(url1, url2 string) bool {
	if url1 == url2 {
		return true
	}
	host1, repo1 := toCanonical(url1)
	if host1 == "" {
		return false
	}
	host2, repo2 := toCanonical(url2)
	return host1 == host2 && repo1 == repo2
}

// PorcelainPresenter is a simple porcelain repo presenter to humans.
var PorcelainPresenter RepoPresenter = func(r *Repo) string {
	if r.vcs == nil {
		// Go package not under VCS.
		return r.Root + "\n	? Not under (recognized) version control"
	}

	s := r.Root + "/..."
	if r.Local.Branch != r.Remote.Branch {
		s += "\n	b Non-default branch checked out"
	}
	if r.Local.Status != "" {
		s += "\n	* Uncommited changes in working dir"
	}
	switch {
	case r.Remote.Revision == "":
		s += "\n	! No remote"
	case !*fFlag && !sameRepoURL(r.Local.RemoteURL, r.Remote.RepoURL):
		s += fmt.Sprintf("\n	# Remote path (%s) doesn't match import path (%s)", r.Remote.RepoURL, r.Local.RemoteURL)
	case r.Local.Revision != r.Remote.Revision:
		if !r.LocalContainsRemoteRevision {
			s += "\n	+ Update available"
		} else {
			s += "\n	- Local revision is ahead of remote"
		}
	}
	if r.Local.Stash != "" {
		s += "\n	$ Stash exists"
	}
	return s
}

// CompactPresenter is a simple porcelain repo presenter to humans in compact form.
var CompactPresenter RepoPresenter = func(r *Repo) string {
	if r.vcs == nil {
		// Go package not under VCS.
		return "???? " + r.Root
	}

	var s string
	switch {
	case r.Local.Branch != r.Remote.Branch:
		s += "b"
	default:
		s += " "
	}
	switch {
	case r.Local.Status != "":
		s += "*"
	default:
		s += " "
	}
	switch {
	case r.Remote.Revision == "":
		s += "!"
	case !*fFlag && !sameRepoURL(r.Local.RemoteURL, r.Remote.RepoURL):
		s += "#"
	case r.Local.Revision != r.Remote.Revision:
		if !r.LocalContainsRemoteRevision {
			s += "+"
		} else {
			s += "-"
		}
	default:
		s += " "
	}
	switch {
	case r.Local.Stash != "":
		s += "$"
	default:
		s += " "
	}
	s += " " + r.Root + "/..."
	return s
}

// DebugPresenter produces verbose debug output.
var DebugPresenter RepoPresenter = func(r *Repo) string {
	b, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		// json.Marshal should never fail to marshal the given struct. If it does, it's a bug
		// in the program and should be fixed.
		panic(err)
	}
	return string(b)
}
