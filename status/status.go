package status

import (
	"fmt"
	"strings"

	// TODO: Make a note about these imports...
	//       Until then, see their godoc pages:
	. "gist.github.com/7480523.git" // http://godoc.org/gist.github.com/7480523.git
)

// PorcelainPresenter is a simple porcelain presenter of GoPackage to humans.
//
// It currently is, and must remain read-only and safe for concurrent execution.
var PorcelainPresenter GoPackageStringer = func(goPackage *GoPackage) string {
	out := ""

	if repo := goPackage.Dir.Repo; repo != nil {
		repoRootImportPath := strings.TrimPrefix(repo.Vcs.RootPath(), goPackage.Bpkg.SrcRoot+"/")

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
		if (strings.HasPrefix(repoRootImportPath, "github.com/") &&
			repo.VcsLocal.Remote != "https://"+repoRootImportPath &&
			repo.VcsLocal.Remote != "https://"+repoRootImportPath+".git") ||
			(strings.HasPrefix(repoRootImportPath, "code.google.com/") &&
				repo.VcsLocal.Remote != "https://"+repoRootImportPath) {
			out += "#"
		} else if repo.VcsLocal.LocalRev != repo.VcsRemote.RemoteRev {
			if repo.VcsRemote.RemoteRev != "" {
				out += "+"
			} else {
				out += "!"
			}
		} else {
			out += " "
		}

		out += " " + repoRootImportPath + "/..."
	} else {
		out += "???"

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
		out += fmt.Sprintf("\tRemote=%q", repo.VcsLocal.Remote)
		out += fmt.Sprintf("\tLocalRev=%q", repo.VcsLocal.LocalRev)
		out += fmt.Sprintf("\tRemoteRev=%q", repo.VcsRemote.RemoteRev)
	}

	return out
}
