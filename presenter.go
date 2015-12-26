package main

import "fmt"

// RepoFilter is a repo filter.
type RepoFilter func(repo *Repo) (show bool)

// RepoPresenter is a repo presenter.
type RepoPresenter func(repo *Repo) string

// PorcelainPresenter is a simple porcelain repo presenter to humans.
//
// It currently is, and must remain read-only and safe for concurrent execution.
var PorcelainPresenter RepoPresenter = func(repo *Repo) string {
	out := ""

	if repo != nil {
		if repo.Local.LocalBranch != repo.VCS.GetDefaultBranch() {
			out += "b"
		} else {
			out += " "
		}
		if repo.Local.Status != "" {
			out += "*"
		} else {
			out += " "
		}
		if repo.Remote.Revision == "" {
			out += "!"
		} else if !*fFlag && (repo.Local.RemoteURL != repo.Remote.RepoURL) {
			out += "#"
		} else if repo.Local.Revision != repo.Remote.Revision {
			if !repo.Remote.IsContained {
				out += "+"
			} else {
				out += "-"
			}
		} else {
			out += " "
		}
		if repo.Local.Stash != "" {
			out += "$"
		} else {
			out += " "
		}

		out += " " + repo.Root + "/..."
	} else {
		out += "????"

		out += " " + "<goPackage.Bpkg.ImportPath>" // TODO.
	}

	return out
}

// DebugPresenter produces debug output.
//
// It currently is, and must remain read-only and safe for concurrent execution.
var DebugPresenter RepoPresenter = func(repo *Repo) string {
	out := ""

	out += fmt.Sprintf("\tRootPath=%q", repo.VCS.RootPath())
	out += fmt.Sprintf("\tRoot=%q", repo.Root)
	out += fmt.Sprintf("\tLocalBranch=%q", repo.Local.LocalBranch)
	out += fmt.Sprintf("\tDefaultBranch=%q", repo.VCS.GetDefaultBranch())
	out += fmt.Sprintf("\tStatus=%q", repo.Local.Status)
	out += fmt.Sprintf("\tStash=%q", repo.Local.Stash)
	out += fmt.Sprintf("\tRemote.RepoURL=%q", repo.Remote.RepoURL)
	out += fmt.Sprintf("\tLocal.RemoteURL=%q", repo.Local.RemoteURL)
	out += fmt.Sprintf("\tLocal.Revision=%q", repo.Local.Revision)
	out += fmt.Sprintf("\tRemoteURL=%q", repo.Local.RemoteURL)
	out += fmt.Sprintf("\tIsContained=%v", repo.Remote.IsContained)

	return out
}
