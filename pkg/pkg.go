package pkg

import vcs "github.com/shurcooL/go/vcs"

type Repo struct {
	// Root is the import path corresponding to the root of the repository.
	Root string

	VCS vcs.Vcs

	Local  Local
	Remote Remote
}

type Local struct {
	// RemoteURL is the remote URL, including scheme.
	RemoteURL string

	Revision string

	Status      string
	Stash       string
	LocalBranch string
}

type Remote struct {
	// RepoURL is the repository URL, including scheme, as determined from the import path.
	RepoURL string

	Revision    string
	IsContained bool // True if remote commit is contained in the default local branch.
}

type RepoFilter func(repo *Repo) (show bool)

type RepoStringer func(repo *Repo) string
