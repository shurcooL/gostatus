package main

import "encoding/json"

// RepoFilter is a repo filter.
type RepoFilter func(r *Repo) (show bool)

// RepoPresenter is a repo presenter.
// All implementations must be read-only and safe for concurrent execution.
type RepoPresenter func(r *Repo) string

// PorcelainPresenter is a simple porcelain repo presenter to humans.
var PorcelainPresenter RepoPresenter = func(r *Repo) string {
	if r.vcs == nil {
		// Go package not under VCS.
		return "????" + " " + r.Root
	}

	var s string
	if r.Local.Branch != r.Remote.Branch {
		s += "b"
	} else {
		s += " "
	}
	if r.Local.Status != "" {
		s += "*"
	} else {
		s += " "
	}
	if r.Remote.Revision == "" {
		s += "!"
	} else if !*fFlag && (r.Local.RemoteURL != r.Remote.RepoURL) {
		s += "#"
	} else if r.Local.Revision != r.Remote.Revision {
		if !r.LocalContainsRemoteRevision {
			s += "+"
		} else {
			s += "-"
		}
	} else {
		s += " "
	}
	if r.Local.Stash != "" {
		s += "$"
	} else {
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
