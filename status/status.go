package status

import (
	// TODO: Make a note about these imports...
	//       Until then, see their godoc pages:
	. "gist.github.com/7480523.git" // http://godoc.org/gist.github.com/7480523.git
)

// PorcelainPresenter is a simple porcelain presenter of GoPackage to humans.
// It is currently the same as the PlumbingPresenter, but this may evolve.
//
// It currently is, and must remain read-only and safe for concurrent execution.
var PorcelainPresenter GoPackageStringer = PlumbingPresenter

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
