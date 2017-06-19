package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/shurcooL/go/indentwriter"
	"github.com/shurcooL/gostatus/status"
)

// RepoFilter is a repo filter.
type RepoFilter func(r *Repo) (show bool)

// RepoPresenter is a repo presenter.
// All implementations must be read-only and safe for concurrent execution.
type RepoPresenter func(r *Repo) string

// PorcelainPresenter is a simple porcelain repo presenter to humans.
var PorcelainPresenter RepoPresenter = func(r *Repo) string {
	if r.vcsError != nil {
		return CompactPresenter(r) + "\n	? Unsupported version control: " + r.vcsError.Error()
	}
	if r.vcs == nil {
		// Go package not under VCS.
		return CompactPresenter(r) + "\n	? Not under version control"
	}

	s := CompactPresenter(r)
	if r.Local.Branch != r.Remote.Branch {
		s += "\n	b Non-default branch checked out"
	}
	if r.Local.Status != "" {
		s += "\n	* Uncommited changes in working dir"
	}
	switch {
	case r.Local.RemoteURL == "":
		s += "\n	! No remote"
	case r.Remote.NotFound != nil:
		s += "\n	/ Remote repository not found (was it deleted? made private?):" +
			"\n" + indent(r.Remote.NotFound.Error())
	case r.Remote.Revision == "":
		s += "\n	? Unreachable remote (check your connection)"
	case !*fFlag && !status.EqualRepoURLs(r.Local.RemoteURL, r.Remote.RepoURL):
		s += "\n	# Remote URL doesn't match repo URL inferred from import path:" +
			fmt.Sprintf("\n		  (actual) %s", r.Local.RemoteURL) +
			fmt.Sprintf("\n		(expected) %s", status.FormatRepoURL(r.Local.RemoteURL, r.Remote.RepoURL))
	case r.Local.Revision != r.Remote.Revision:
		switch {
		case !r.Local.ContainsRemoteRevision && r.Remote.ContainsLocalRevision:
			s += "\n	+ Update available"
		case r.Local.ContainsRemoteRevision && !r.Remote.ContainsLocalRevision:
			s += "\n	- Local revision is ahead of remote revision"
		case !r.Local.ContainsRemoteRevision && !r.Remote.ContainsLocalRevision:
			s += "\n	± Update available; local revision is ahead of remote revision"
		default:
			panic(fmt.Errorf("internal error: both r.Local.ContainsRemoteRevision and r.Remote.ContainsLocalRevision are true, yet r.Local.Revision != r.Remote.Revision; this shouldn't be possible, please report if it happens"))
		}
	}
	if r.Local.Stash != "" {
		s += "\n	$ Stash exists"
	}
	return s
}

// indent indents s by 2 tabs.
func indent(s string) string {
	var buf bytes.Buffer
	w := indentwriter.New(&buf, 2)
	_, err := io.WriteString(w, s)
	if err != nil {
		panic(fmt.Errorf("indent writes to bytes.Buffer, and should never fail, yet it did: %v", err))
	}
	return buf.String()
}

// CompactPresenter is a simple porcelain repo presenter to humans in compact form.
var CompactPresenter RepoPresenter = func(r *Repo) string {
	if r.vcsError != nil {
		return "???? " + r.Root + "/..."
	}
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
	case r.Local.RemoteURL == "":
		s += "!"
	case r.Remote.NotFound != nil:
		s += "/"
	case r.Remote.Revision == "":
		s += "?"
	case !*fFlag && !status.EqualRepoURLs(r.Local.RemoteURL, r.Remote.RepoURL):
		s += "#"
	case r.Local.Revision != r.Remote.Revision:
		switch {
		case !r.Local.ContainsRemoteRevision && r.Remote.ContainsLocalRevision:
			s += "+"
		case r.Local.ContainsRemoteRevision && !r.Remote.ContainsLocalRevision:
			s += "-"
		case !r.Local.ContainsRemoteRevision && !r.Remote.ContainsLocalRevision:
			s += "±"
		default:
			panic(fmt.Errorf("internal error: both r.Local.ContainsRemoteRevision and r.Remote.ContainsLocalRevision are true, yet r.Local.Revision != r.Remote.Revision; this shouldn't be possible, please report if it happens"))
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
