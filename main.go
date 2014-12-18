// A command line tool that shows the status of (many) Go packages.
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/kisielk/gotool"
	"github.com/shurcooL/gostatus/status"

	"github.com/shurcooL/go/gists/gist7480523"
	"github.com/shurcooL/go/gists/gist7651991"
)

const numWorkers = 8

var vFlag = flag.Bool("v", false, "Verbose output: show all Go packages, not just ones with notable status.")
var stdinFlag = flag.Bool("stdin", false, "Read the list of newline separated Go packages from stdin.")
var plumbingFlag = flag.Bool("plumbing", false, "Give the output in an easy-to-parse format for scripts.")
var debugFlag = flag.Bool("debug", false, "Give the output with verbose debug information.")

func init() {
	flag.BoolVar(&status.FFlag, "f", false, "Force not to verify that each package has been checked out from the source control repository implied by its import path. This can be useful if the source is a local fork of the original.")
}

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
	os.Exit(2)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Usage = usage
	flag.Parse()

	// Get current directory.
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	shouldShow := func(goPackage *gist7480523.GoPackage) bool {
		// Check for notable status.
		return status.PorcelainPresenter(goPackage)[:4] != "    "
	}
	if *vFlag == true {
		shouldShow = func(_ *gist7480523.GoPackage) bool { return true }
	}

	var presenter gist7480523.GoPackageStringer = status.PorcelainPresenter
	if *debugFlag == true {
		presenter = status.DebugPresenter
	} else if *plumbingFlag == true {
		presenter = status.PlumbingPresenter
	}

	// A map of repos that have been checked, to avoid doing same repo more than once.
	var lock sync.Mutex
	checkedRepos := map[string]bool{}

	// Input: Go package Import Path
	// Output: If a valid Go package and not inside GOROOT, output a status string, else nil.
	reduceFunc := func(in string) interface{} {
		goPackage, err := gist7480523.GoPackageFromPath(in, wd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't load package: %s\n", err)
			return nil
		}
		if goPackage == nil {
			return nil
		}
		if goPackage.Bpkg.Goroot {
			return nil
		}

		goPackage.UpdateVcs()
		// Check that the same repo hasn't already been done.
		if goPackage.Dir.Repo != nil {
			rootPath := goPackage.Dir.Repo.Vcs.RootPath()
			lock.Lock()
			if !checkedRepos[rootPath] {
				checkedRepos[rootPath] = true
				lock.Unlock()
			} else {
				lock.Unlock()
				// Skip repos that were done.
				return nil
			}
		}

		goPackage.UpdateVcsFields()

		if shouldShow(goPackage) == false {
			return nil
		}
		return presenter(goPackage)
	}

	// Run reduceFunc on all import paths in parallel.
	var outChan <-chan interface{}
	switch *stdinFlag {
	case false:
		importPathPatterns := flag.Args()
		if len(importPathPatterns) == 0 {
			importPathPatterns = []string{"."}
		}
		importPaths := gotool.ImportPaths(importPathPatterns)
		outChan = gist7651991.GoReduceLinesFromSlice(importPaths, numWorkers, reduceFunc)
	case true:
		outChan = gist7651991.GoReduceLinesFromReader(os.Stdin, numWorkers, reduceFunc)
	}

	// Output results.
	for out := range outChan {
		fmt.Println(out.(string))
	}
}
