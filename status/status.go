package status

import (
	"fmt"

	"github.com/shurcooL/gostatus/pkg"
)

// PorcelainPresenter is a simple porcelain presenter of GoPackage to humans.
//
// It currently is, and must remain read-only and safe for concurrent execution.
var PorcelainPresenter pkg.RepoStringer = PlumbingPresenterV3

// Force not to verify that each package has been checked out from the source control repository implied by its import path. This can be useful if the source is a local fork of the original.
var FFlag bool

// This format should remain stable across versions and regardless of user configuration.
//
// It currently is, and must remain read-only and safe for concurrent execution.
var PlumbingPresenterV3 pkg.RepoStringer = func(repo *pkg.Repo) string {
	out := ""

	if repo != nil {
		// TODO: Take care of symlinks?
		//repoImportPath := gist7480523.GetRepoImportPath(repo.Vcs.RootPath(), goPackage.Bpkg.SrcRoot)
		repoImportPath := repo.Root

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
		} else if !FFlag && (false /*repo.RepoRoot == nil || repo.RepoRoot.Repo != repo.VcsLocal.Remote*/) {
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

		out += " " + repoImportPath + "/..."
	} else {
		out += "????"

		out += " " + "<goPackage.Bpkg.ImportPath>" // TODO.
	}

	return out
}

// DebugPresenter gives debug output.
//
// It currently is, and must remain read-only and safe for concurrent execution.
var DebugPresenter pkg.RepoStringer = func(repo *pkg.Repo) string {
	out := ""

	//out += fmt.Sprintf("\tRootPath=%q", repo.Vcs.RootPath())
	out += fmt.Sprintf("\tRoot=%q", repo.Root)
	out += fmt.Sprintf("\tLocalBranch=%q", repo.Local.LocalBranch)
	out += fmt.Sprintf("\tDefaultBranch=%q", repo.VCS.GetDefaultBranch())
	out += fmt.Sprintf("\tStatus=%q", repo.Local.Status)
	out += fmt.Sprintf("\tStash=%q", repo.Local.Stash)
	/*if repo.RepoRoot == nil {
		out += fmt.Sprintf("\tRepoRoot=<nil>")
	} else {
		out += fmt.Sprintf("\tRepoRoot.Repo=%q", repo.RepoRoot.Repo)
	}*/
	out += fmt.Sprintf("\tRemoteURL=%q", repo.RemoteURL)
	out += fmt.Sprintf("\tLocal.Revision=%q", repo.Local.Revision)
	out += fmt.Sprintf("\tRemote.Revision=%q", repo.Remote.Revision)
	out += fmt.Sprintf("\tIsContained=%v", repo.Remote.IsContained)

	return out
}
