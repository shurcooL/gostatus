package status

import (
	"fmt"
	"strings"

	. "github.com/shurcooL/go/gists/gist7480523"
)

// PorcelainPresenter is a simple porcelain presenter of GoPackage to humans.
//
// It currently is, and must remain read-only and safe for concurrent execution.
var PorcelainPresenter GoPackageStringer = PlumbingPresenterV2

// This format should remain stable across versions and regardless of user configuration.
//
// It currently is, and must remain read-only and safe for concurrent execution.
var PlumbingPresenterV2 GoPackageStringer = func(goPackage *GoPackage) string {
	out := ""

	if repo := goPackage.Dir.Repo; repo != nil {
		repoImportPath := GetRepoImportPath(repo.Vcs.RootPath(), goPackage.Bpkg.SrcRoot)

		if repo.VcsLocal.LocalBranch != repo.Vcs.GetDefaultBranch() {
			out += "b"
		} else {
			out += " "
		}
		if repo.VcsLocal.Status != "" {
			out += "*"
		} else {
			out += " "
		}
		if (strings.HasPrefix(repoImportPath, "github.com/") &&
			repo.VcsLocal.Remote != "https://"+repoImportPath &&
			repo.VcsLocal.Remote != "https://"+repoImportPath+".git") ||
			(strings.HasPrefix(repoImportPath, "code.google.com/") &&
				repo.VcsLocal.Remote != "https://"+repoImportPath) {
			out += "#"
		} else if repo.VcsLocal.LocalRev != repo.VcsRemote.RemoteRev {
			if repo.VcsRemote.RemoteRev != "" {
				if !repo.VcsRemote.IsContained {
					out += "+"
				} else {
					out += "-"
				}
			} else {
				out += "!"
			}
		} else {
			out += " "
		}
		if repo.VcsLocal.Stash != "" {
			out += "$"
		} else {
			out += " "
		}

		out += " " + repoImportPath + "/..."
	} else {
		out += "????"

		out += " " + goPackage.Bpkg.ImportPath
	}

	return out
}

// PlumbingPresenter gives the output in an easy-to-parse format for scripts.
// This format should remain stable across versions and regardless of user configuration.
//
// It currently is, and must remain read-only and safe for concurrent execution.
var PlumbingPresenter GoPackageStringer = func(goPackage *GoPackage) string {
	out := ""

	if repo := goPackage.Dir.Repo; repo != nil {
		out += "@"
		if repo.VcsLocal.LocalBranch != repo.Vcs.GetDefaultBranch() {
			out += "b"
		} else {
			out += " "
		}
		if repo.VcsLocal.Status != "" {
			out += "*"
		} else {
			out += " "
		}
		if repo.VcsLocal.LocalRev != repo.VcsRemote.RemoteRev {
			out += "+"
		} else {
			out += " "
		}
	} else {
		out += "    "
	}
	if goPackage.Bpkg.IsCommand() {
		out += "/"
	} else {
		out += " "
	}

	out += " " + goPackage.Bpkg.ImportPath

	return out
}

// DebugPresenter gives debug output.
//
// It currently is, and must remain read-only and safe for concurrent execution.
var DebugPresenter GoPackageStringer = func(goPackage *GoPackage) string {
	out := goPackage.Bpkg.ImportPath

	out += fmt.Sprintf("\tgoPackage.Dir.Repo=%p", goPackage.Dir.Repo)
	out += fmt.Sprintf("\tgoPackage.Bpkg.SrcRoot=%q", goPackage.Bpkg.SrcRoot)

	if repo := goPackage.Dir.Repo; repo != nil {
		out += fmt.Sprintf("\tRootPath=%q", repo.Vcs.RootPath())
		out += fmt.Sprintf("\tLocalBranch=%q", repo.VcsLocal.LocalBranch)
		out += fmt.Sprintf("\tDefaultBranch=%q", repo.Vcs.GetDefaultBranch())
		out += fmt.Sprintf("\tStatus=%q", repo.VcsLocal.Status)
		out += fmt.Sprintf("\tStash=%q", repo.VcsLocal.Stash)
		out += fmt.Sprintf("\tRemote=%q", repo.VcsLocal.Remote)
		out += fmt.Sprintf("\tLocalRev=%q", repo.VcsLocal.LocalRev)
		out += fmt.Sprintf("\tRemoteRev=%q", repo.VcsRemote.RemoteRev)
		out += fmt.Sprintf("\tIsContained=%v", repo.VcsRemote.IsContained)
	}

	return out
}
