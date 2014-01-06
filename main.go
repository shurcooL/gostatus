// A command line tool that shows the status of (many) Go packages.
package main

import (
	"flag"
	"fmt"
	"github.com/shurcooL/gostatus/status"
	"os"
	"runtime"
	"strings"
	"sync"
	// TODO: Make a note about these imports...
	//       Until then, see their godoc pages:
	. "gist.github.com/7480523.git" // http://godoc.org/gist.github.com/7480523.git
	. "gist.github.com/7651991.git" // http://godoc.org/gist.github.com/7651991.git
)

func usage() {
	const legend = `
Examples:
  # Show status of all your packages
  go list all | gostatus

  # Show status of all dependencies (recursive) of package in cur working dir
  go list -f '{{join .Deps "\n"}}' . | gostatus

Legend:
  @ - Vcs repo
  b - Non-master branch checked out
  * - Uncommited changes in working dir
  + - Update available (latest remote revision doesn't match local revision)
  / - Command (package main)
`

	fmt.Fprint(os.Stderr, "Usage: [newline separated packages] | gostatus\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, legend)
	os.Exit(2)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var short = flag.Bool("short", false, "Only show modified or branch (short) packages.")
	flag.Usage = usage
	flag.Parse()

	var presenter GoPackageStringer = status.Presenter

	// A map of repos that have been checked, to avoid doing same repo more than once
	var lock sync.Mutex
	checkedRepos := map[string]bool{}

	// Input: Go package Import Path
	// Output: If a valid Go package and not part of standard library, output a status string, else nil
	reduceFunc := func(in string) interface{} {
		if goPackage := GoPackageFromImportPath(in); goPackage != nil {
			if !goPackage.Standard {
				// HACK: Check that the same repo hasn't already been done
				if goPackage.Vcs != nil {
					rootPath := goPackage.Vcs.RootPath()
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

				goPackage.UpdateVcsFields()
				return presenter(goPackage)
			}
		}
		return nil
	}

	// Run reduceFunc on all lines from stdin in parallel (max 8 goroutines)
	outChan := GoReduceLinesFromReader(os.Stdin, 8, reduceFunc)

	// Output results
	for out := range outChan {
		outLine := out.(string)
		strArry := strings.Split(outLine, "")
		if *short {
			if strArry[1] == "b" || strArry[2] == "*" || strArry[3] == "+" {
				fmt.Println(outLine)
			}
		} else {
			fmt.Println(outLine)
		}
	}
}
