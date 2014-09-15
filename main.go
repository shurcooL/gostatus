// A command line tool that shows the status of (many) Go packages.
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/shurcooL/gostatus/status"

	. "github.com/shurcooL/go/gists/gist7480523"
	. "github.com/shurcooL/go/gists/gist7651991"
	"github.com/shurcooL/go/gists/gist7802150"
)

var allFlag = flag.Bool("all", false, "Show all Go packages, not just ones with notable status.")
var plumbingFlag = flag.Bool("plumbing", false, "Give the output in an easy-to-parse format for scripts.")
var debugFlag = flag.Bool("debug", false, "Give the output with verbose debug information.")

func usage() {
	fmt.Fprint(os.Stderr, "Usage: [newline separated packages] | gostatus [flags]\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, `
Examples:
  # Show status of packages with notable status.
  go list all | gostatus

  # Show status of all dependencies (recursive) of package in cur working dir.
  go list -f '{{join .Deps "\n"}}' . | gostatus --all

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
	os.Exit(2)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Usage = usage
	flag.Parse()

	shouldShow := func(goPackage *GoPackage) bool {
		// Check for notable status.
		return status.PorcelainPresenter(goPackage)[:4] != "    " && status.PorcelainPresenter(goPackage)[:4] != "????"
	}

	if *allFlag == true {
		shouldShow = func(_ *GoPackage) bool { return true }
	}

	var presenter GoPackageStringer = status.PorcelainPresenter

	if *debugFlag == true {
		presenter = status.DebugPresenter
	} else if *plumbingFlag == true {
		presenter = status.PlumbingPresenter
	}

	// A map of repos that have been checked, to avoid doing same repo more than once
	var lock sync.Mutex
	checkedRepos := map[string]bool{}

	// Input: Go package Import Path
	// Output: If a valid Go package and not part of standard library, output a status string, else nil
	reduceFunc := func(in string) interface{} {
		goPackage := GoPackageFromImportPath(in)
		if goPackage == nil {
			return nil
		}
		if goPackage.Standard {
			return nil
		}

		goPackage.UpdateVcs()
		// HACK: Check that the same repo hasn't already been done
		if goPackage.Dir.Repo != nil {
			rootPath := goPackage.Dir.Repo.Vcs.RootPath()
			lock.Lock()
			if !checkedRepos[rootPath] {
				checkedRepos[rootPath] = true
				lock.Unlock()
			} else {
				lock.Unlock()
				// TODO: Instead of skipping repos that were done, cache their state and reuse it
				return nil
			}
		}

		if goPackage.Dir.Repo != nil {
			gist7802150.MakeUpdated(goPackage.Dir.Repo.VcsLocal)
			gist7802150.MakeUpdated(goPackage.Dir.Repo.VcsRemote)
		}

		if shouldShow(goPackage) == false {
			return nil
		}
		return presenter(goPackage)
	}

	// Run reduceFunc on all lines from stdin in parallel (max 8 goroutines)
	outChan := GoReduceLinesFromReader(os.Stdin, 8, reduceFunc)

	// Output results
	for out := range outChan {
		fmt.Println(out.(string))
	}
}
