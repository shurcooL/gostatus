package main

import (
	"encoding/json"
	"fmt"
)

// RepoFilter is a repo filter.
type RepoFilter func(r *Repo) (show bool)

// RepoPresenter is a repo presenter.
// All implementations must be read-only and safe for concurrent execution.
type RepoPresenter func(r *Repo) string

// PorcelainPresenter is a simple porcelain repo presenter to humans.
var PorcelainPresenter RepoPresenter = func(r *Repo) string {
	if r.vcs == nil {
		// Go package not under VCS.
		return CompactPresenter(r) + "\n	? Not under (recognized) version control"
	}

	s := CompactPresenter(r)
	if r.Local.Branch != r.Remote.Branch {
		s += "\n	b Non-default branch checked out"
	}
	if r.Local.Status != "" {
		s += "\n	* Uncommited changes in working dir"
	}
	switch {
	case r.Remote.Revision == "":
		s += "\n	! No remote"
	case !*fFlag && !equalRepoURLs(r.Local.RemoteURL, r.Remote.RepoURL):
		s += fmt.Sprintf("\n	# Remote URL (%s) doesn't match repo URL inferred from import path (%s)", r.Local.RemoteURL, r.Remote.RepoURL)
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
	case !*fFlag && !equalRepoURLs(r.Local.RemoteURL, r.Remote.RepoURL):
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
