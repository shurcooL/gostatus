// A command line tool that shows the status of (many) Go packages.
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	. "gist.github.com/7480523.git"
	. "gist.github.com/7651991.git"
)

func usage() {
	const legend = `
Legend:
  @ - Git repo
  b - Non-master branch checked out
  * - Uncommited changes in working dir
  + - Update available (latest remote revision doesn't match local revision)
  / - Command (package main)
`

	fmt.Fprint(os.Stderr, "usage: [newline separated packages] | gostatus\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, legend)
	os.Exit(2)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Usage = usage
	flag.Parse()

	var lock sync.Mutex
	checkedGitRepos := map[string]bool{}

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
						return "@---- " + x.Bpkg.ImportPath
					}
				}

				x.Update()
				return x.String()
			}
		}
		return nil
	}

	outChan := GoReduceLinesFromReader(os.Stdin, 8, reduceFunc)

	for out := range outChan {
		// TODO: Instead of skipping git repos that were done, cache their state and reuse it
		if strings.HasPrefix(out.(string), "@---- ") {
			continue
		}

		fmt.Println(out.(string))
	}
}
