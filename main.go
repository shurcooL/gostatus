package main

import (
	"os"
	"runtime"
	"strings"

	. "gist.github.com/7480523.git"
	. "gist.github.com/7651991.git"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	reduceFunc := func(in string) interface{} {
		if x := SomethingFromImportPath(in); x != nil {
			Standard := x.Bpkg.Goroot && x.Bpkg.ImportPath != "" && !strings.Contains(x.Bpkg.ImportPath, ".")

			if !Standard {
				x.Update()
				return x.String()
			}
		}
		return nil
	}

	outChan := GoReduceLinesFromReader(os.Stdin, 8, reduceFunc)

	for out := range outChan {
		println(out.(string))
	}
}
