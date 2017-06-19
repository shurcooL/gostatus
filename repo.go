package main

import "github.com/shurcooL/vcsstate"

// Repo represents a repository that contains Go packages and its state when VCS is non-nil.
// It represents a Go package that is not under a VCS when VCS is nil.
type Repo struct {
	// Path is the local filesystem path to the repository or Go package.
	Path string

	// Root is the import path corresponding to the root of the repository or Go package.
	Root string

	// vcs allows getting the state of the VCS. It's nil if there's no VCS.
	vcs      vcsstate.VCS
	vcsError error

	Local struct {
		// RemoteURL is the remote URL, including scheme.
		RemoteURL string

		Status   string
		Branch   string // Checked out branch.
		Revision string
		Stash    string

		ContainsRemoteRevision bool // Computed if Remote.Revision != "".
	}
	Remote struct {
		// RepoURL is the repository URL, including scheme, as determined dynamically from the import path.
		RepoURL string

		NotFound error  // Whether remote repository was not found.
		Branch   string // Default branch, as determined from remote.
		Revision string

		ContainsLocalRevision bool // Computed if Local.Revision != "".
	}
}
