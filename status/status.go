package status

import (
	// TODO: Make a note about these imports...
	//       Until then, see their godoc pages:
	. "gist.github.com/7480523.git" // http://godoc.org/gist.github.com/7480523.git
)

// Presenter is a simple porcelain presenter of Something to humans.
// It may change, so don't parse it; another plumbing presenter should be used for that.
//
// It currently is, and must remain read-only and safe for concurrent execution.
var Presenter SomethingStringer = func(w *Something) string {
	out := ""

	if w.IsGitRepo {
		out += "@"
		if w.LocalBranch != "master" {
			out += "b"
		} else {
			out += " "
		}
		if w.Status != "" {
			out += "*"
		} else {
			out += " "
		}
		if w.Remote != w.Local {
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
