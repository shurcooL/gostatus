package status

import (
	"fmt"

	// TODO: Make a note about these imports...
	//       Until then, see their godoc pages:
	. "gist.github.com/7480523.git" // http://godoc.org/gist.github.com/7480523.git
)

// PorcelainPresenter is a simple porcelain presenter of GoPackage to humans.
//
// It currently is, and must remain read-only and safe for concurrent execution.
var PorcelainPresenter GoPackageStringer = func(w *GoPackage) string {
	out := ""

	if w := w.Dir.Repo; w != nil {
		out += " "
		if w.VcsLocal.LocalBranch != w.Vcs.GetDefaultBranch() {
			out += "b"
		} else {
			out += " "
		}
		if w.VcsLocal.Status != "" {
			out += "*"
		} else {
			out += " "
		}
		if w.VcsLocal.LocalRev != w.VcsRemote.RemoteRev {
			out += "+"
		} else {
			out += " "
		}
	} else {
		out += "?   "
	}
	if w.Bpkg.IsCommand() {
		out += "/"
	} else {
		out += " "
	}

	out += " " + w.Bpkg.ImportPath

	return out
}

// PlumbingPresenter gives the output in an easy-to-parse format for scripts.
// This format should remain stable across versions and regardless of user configuration.
//
// It currently is, and must remain read-only and safe for concurrent execution.
var PlumbingPresenter GoPackageStringer = func(w *GoPackage) string {
	out := ""

	if w := w.Dir.Repo; w != nil {
		out += "@"
		if w.VcsLocal.LocalBranch != w.Vcs.GetDefaultBranch() {
			out += "b"
		} else {
			out += " "
		}
		if w.VcsLocal.Status != "" {
			out += "*"
		} else {
			out += " "
		}
		if w.VcsLocal.LocalRev != w.VcsRemote.RemoteRev {
			out += "+"
		} else {
			out += " "
		}
	} else {
		out += "    "
	}
	if w.Bpkg.IsCommand() {
		out += "/"
	} else {
		out += " "
	}

	out += " " + w.Bpkg.ImportPath

	return out
}

// DebugPresenter gives debug output.
//
// It currently is, and must remain read-only and safe for concurrent execution.
var DebugPresenter GoPackageStringer = func(w *GoPackage) string {
	out := w.Bpkg.ImportPath

	out += fmt.Sprintf("\tgoPackage.Dir.Repo=%p", w.Dir.Repo)

	if w := w.Dir.Repo; w != nil {
		out += fmt.Sprintf("\tRootPath=%q", w.Vcs.RootPath())
		out += fmt.Sprintf("\tLocalBranch=%q", w.VcsLocal.LocalBranch)
		out += fmt.Sprintf("\tDefaultBranch=%q", w.Vcs.GetDefaultBranch())
		out += fmt.Sprintf("\tStatus=%q", w.VcsLocal.Status)
		out += fmt.Sprintf("\tLocalRev=%q", w.VcsLocal.LocalRev)
		out += fmt.Sprintf("\tRemoteRev=%q", w.VcsRemote.RemoteRev)
	}

	return out
}
