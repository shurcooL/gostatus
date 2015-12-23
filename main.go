// gostatus is a command line tool that shows the status of (many) Go packages.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kisielk/gotool"
	"github.com/shurcooL/gostatus/pkg"
	"github.com/shurcooL/gostatus/status"
)

const numWorkers = 8

var (
	vFlag        = flag.Bool("v", false, "Verbose output: show all Go packages, not just ones with notable status.")
	stdinFlag    = flag.Bool("stdin", false, "Read the list of newline separated Go packages from stdin.")
	plumbingFlag = flag.Bool("plumbing", false, "Give the output in an easy-to-parse format for scripts.")
	debugFlag    = flag.Bool("debug", false, "Give the output with verbose debug information.")
)

func init() {
	flag.BoolVar(&status.FFlag, "f", false, "Force not to verify that each package has been checked out from the source control repository implied by its import path. This can be useful if the source is a local fork of the original.")
}

var wd = func() string {
	// Get current directory.
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("failed to get current directory:", err)
	}
	return wd
}()

func usage() {
	fmt.Fprint(os.Stderr, "Usage: gostatus [flags] [packages]\n")
	fmt.Fprint(os.Stderr, "       [newline separated packages] | gostatus --stdin [flags]\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, `
Examples:
  # Show status of package in current directory, if notable.
  gostatus .

  # Show status of all packages with notable status.
  gostatus all

  # Show status of all dependencies (recursive) of package in cur working dir.
  go list -f '{{join .Deps "\n"}}' . | gostatus --stdin -v

Legend:
  ???? - Not under (recognized) version control
  b - Non-master branch checked out
  * - Uncommited changes in working dir
  + - Update available
  - - Local revision is ahead of remote (need to push?)
  ! - No remote
  # - Remote path doesn't match import path
  $ - Stash exists
`)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	shouldShow := func(goPackage *pkg.Repo) bool {
		// Check for notable status.
		return status.PorcelainPresenter(goPackage)[:4] != "    "
	}
	if *vFlag == true {
		shouldShow = func(*pkg.Repo) bool { return true }
	}

	var presenter pkg.RepoStringer = status.PorcelainPresenter
	if *debugFlag == true {
		presenter = status.DebugPresenter
	} else if *plumbingFlag == true {
		//presenter = status.PlumbingPresenter
	}

	workspace := NewGoWorkspace(shouldShow, presenter)

	switch *stdinFlag {
	case false:
		go func() { // This needs to happen in the background because sending input will be blocked on processing and receiving output.
			importPathPatterns := flag.Args()
			importPaths := gotool.ImportPaths(importPathPatterns)
			for _, importPath := range importPaths {
				workspace.Add(importPath)
			}
			workspace.Done()
		}()
	case true:
		go func() { // This needs to happen in the background because sending input will be blocked on processing and receiving output.
			br := bufio.NewReader(os.Stdin)
			for line, err := br.ReadString('\n'); err == nil; line, err = br.ReadString('\n') {
				importPath := line[:len(line)-1] // Trim last newline.
				workspace.Add(importPath)
			}
			workspace.Done()
		}()
	}

	// Output results.
	for repoStatus := range workspace.Out {
		// DONE: Consider running presenter concurrently if it's slow and saving result, similar to repo.Presenter, etc.
		fmt.Println(repoStatus)
	}
}
