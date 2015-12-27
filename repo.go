package main

import "github.com/shurcooL/vcsstate"

// Repo represents a repository that contains Go packages, and its state.
type Repo struct {
	// Path is the local filesystem path to the repository.
	Path string

	// Root is the import path corresponding to the root of the repository.
	Root string

	VCS vcsstate.VCS

	Local struct {
		// RemoteURL is the remote URL, including scheme.
		RemoteURL string

		Status   string
		Branch   string // Checked out branch.
		Revision string
		Stash    string
	}
	Remote struct {
		// RepoURL is the repository URL, including scheme, as determined dynamically from the import path.
		RepoURL string

		Revision string
	}
	LocalContainsRemoteRevision bool
}
