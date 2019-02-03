// gostatus is a command line tool that shows the status of Go repositories.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kisielk/gotool"
)

// parallelism for workers.
const parallelism = 8

var (
	debugFlag   = flag.Bool("debug", false, "Cause the repository data to be printed in verbose debug format.")
	fFlag       = flag.Bool("f", false, "Force not to verify that each package has been checked out from the source control repository implied by its import path. This can be useful if the source is a local fork of the original.")
	stdinFlag   = flag.Bool("stdin", false, "Read the list of newline separated Go packages from stdin.")
	vFlag       = flag.Bool("v", false, "Verbose mode. Show all Go packages, not just ones with notable status.")
	compactFlag = flag.Bool("c", false, "Compact output with inline notation.")
)

func usage() {
	fmt.Fprint(os.Stderr, "Usage: gostatus [flags] [packages]\n")
	fmt.Fprint(os.Stderr, "       [newline separated packages] | gostatus -stdin [flags]\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, `
Examples:
  # Show status of all packages.
  gostatus all

  # Show status of package in current directory.
  gostatus

  # Show status of all dependencies (recursive) of package in current dir.
  go list -deps | gostatus -stdin -v

Legend:
  ? - Not under version control or unreachable remote
  b - Non-default branch checked out
  * - Uncommited changes in working dir
  + - Update available
  - - Local revision is ahead of remote revision
  ± - Update available; local revision is ahead of remote revision
  ! - No remote
  / - Remote repository not found (was it deleted? made private?)
  # - Remote path doesn't match import path
  $ - Stash exists
`)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var shouldShow RepoFilter
	switch {
	default:
		shouldShow = func(r *Repo) bool {
			// Check for notable status.
			return CompactPresenter(r)[:4] != "    "
		}
	case *vFlag:
		shouldShow = func(*Repo) bool { return true }
	}

	var presenter RepoPresenter
	switch {
	case *debugFlag:
		presenter = DebugPresenter
	case *compactFlag:
		presenter = CompactPresenter
	default:
		presenter = PorcelainPresenter
	}

	workspace := NewWorkspace(shouldShow, presenter)

	// Feed input into workspace processing pipeline.
	switch *stdinFlag {
	case false:
		go func() { // This needs to happen in the background because sending input will be blocked on processing and receiving output.
			importPaths := gotool.ImportPaths(flag.Args())
			for _, importPath := range importPaths {
				workspace.ImportPaths <- importPath
			}
			close(workspace.ImportPaths)
		}()
	case true:
		go func() { // This needs to happen in the background because sending input will be blocked on processing and receiving output.
			br := bufio.NewReader(os.Stdin)
			for line, err := br.ReadString('\n'); err == nil; line, err = br.ReadString('\n') {
				importPath := line[:len(line)-1] // Trim last newline.
				workspace.ImportPaths <- importPath
			}
			close(workspace.ImportPaths)
		}()
	}

	// Output results.
	for workspace.Statuses != nil || workspace.Errors != nil {
		select {
		case status, ok := <-workspace.Statuses:
			if !ok {
				workspace.Statuses = nil
				continue
			}
			fmt.Println(status)
		case error, ok := <-workspace.Errors:
			if !ok {
				workspace.Errors = nil
				continue
			}
			fmt.Fprintln(os.Stderr, error)
		}
	}
}

var wd = func() string {
	// Get current directory.
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("failed to get current directory:", err)
	}
	return wd
}()
