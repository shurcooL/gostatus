package pkg

import vcs2 "github.com/shurcooL/go/vcs"

type Repo struct {
	// Root is the import path corresponding to the root of the repository.
	// TODO: Consider. Overlaps with RR.
	Root string

	// RemoteURL is the remote URL, including scheme.
	// TODO: Consider. Overlaps with RR.
	RemoteURL string

	// TODO: Consider. Overlaps with RR.
	VCS vcs2.Vcs

	Local  Local
	Remote Remote
}

type Local struct {
	Revision string

	Status string
	Stash  string
	//RemoteURL      string
	LocalBranch string
}

type Remote struct {
	Revision    string
	IsContained bool // True if remote commit is contained in the default local branch.
}

// RepoImportPath returns what would be the import path of the root folder of the repository. It may or may not
// be an actual Go package. E.g.,
//
// 	"github.com/owner/repo"
func (r Repo) RepoImportPath() string {
	return r.Root
}

// ImportPathPattern returns an import path pattern that matches all of the Go packages in this repo.
// E.g.,
//
// 	"github.com/owner/repo/..."
func (r Repo) ImportPathPattern() string {
	return r.Root + "/..."
}

type RepoStringer func(repo *Repo) string
