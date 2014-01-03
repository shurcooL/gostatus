// A command line tool that shows the status of (many) Go packages.
package main

import (
	"flag"
	"fmt"
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
  @ - Git repo
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

	flag.Usage = usage
	flag.Parse()

	// A map of git repos that have been checked, to avoid doing same git repo more than once
	var lock sync.Mutex
	checkedGitRepos := map[string]bool{}

	// Input: Go package Import Path
	// Output: If a valid Go package and not part of standard library, output a status string, else nil
	reduceFunc := func(in string) interface{} {
		if x := SomethingFromImportPath(in); x != nil {
			Standard := x.Bpkg.Goroot && x.Bpkg.ImportPath != "" && !strings.Contains(x.Bpkg.ImportPath, ".")

			if !Standard {
				// HACK: Check that the same git repo hasn't already been done
				if isGitRepo, rootPath := GetGitRepoRoot(x.Path); isGitRepo {
					lock.Lock()
					if !checkedGitRepos[rootPath] {
						checkedGitRepos[rootPath] = true
						lock.Unlock()
					} else {
						lock.Unlock()
						// TODO: Instead of skipping git repos that were done, cache their state and reuse it
						return nil
					}
				}

				x.Update()
				return x.String()
			}
		}
		return nil
	}

	// Run reduceFunc on all lines from stdin in parallel (max 8 goroutines)
	outChan := GoReduceLinesFromReader(os.Stdin, 8, reduceFunc)

	// Output results
	for out := range outChan {
		fmt.Println(out.(string))
	}
}
