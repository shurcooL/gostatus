package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	. "gist.github.com/7480523.git"
	. "gist.github.com/7651991.git"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

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
