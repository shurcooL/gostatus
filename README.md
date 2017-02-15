gostatus
========

[![Build Status](https://travis-ci.org/shurcooL/gostatus.svg?branch=master)](https://travis-ci.org/shurcooL/gostatus) [![GoDoc](https://godoc.org/github.com/shurcooL/gostatus?status.svg)](https://godoc.org/github.com/shurcooL/gostatus)

gostatus is a command line tool that shows the status of Go repositories.

Installation
------------

```bash
go get -u github.com/shurcooL/gostatus
```

Usage
-----

```bash
Usage: gostatus [flags] [packages]
       [newline separated packages] | gostatus -stdin [flags]
  -c	Compact output with inline notation.
  -debug
    	Cause the repository data to be printed in verbose debug format.
  -f	Force not to verify that each package has been checked out from the source control repository implied by its import path. This can be useful if the source is a local fork of the original.
  -stdin
    	Read the list of newline separated Go packages from stdin.
  -v	Verbose mode. Show all Go packages, not just ones with notable status.

Examples:
  # Show status of all packages.
  gostatus all

  # Show status of package in current directory.
  gostatus

  # Show status of all dependencies (recursive) of package in current dir.
  go list -f '{{join .Deps "\n"}}' . | gostatus -stdin -v

Legend:
  ? - Not under version control or unreachable remote
  b - Non-default branch checked out
  * - Uncommited changes in working dir
  + - Update available
  - - Local revision is ahead of remote
  ! - No remote
  / - Remote repository not found (was it deleted? made private?)
  # - Remote path doesn't match import path
  $ - Stash exists
```

Examples
--------

```bash
# Show status of all packages.
$ gostatus all

# Show status of package in current directory.
$ gostatus

# Show status of specified package.
$ gostatus import/path

# Show status of all dependencies (recursive) of package in current dir.
$ go list -f '{{join .Deps "\n"}}' . | gostatus -stdin -v

# Show status of all dependencies (recursive) of specified package.
$ go list -f '{{join .Deps "\n"}}' import/path | gostatus -stdin -v
```

Sample Output
-------------

```bash
$ gostatus all
  +  github.com/dchest/uniuri/...
	+ Update available
  +  github.com/syndtr/goleveldb/...
	+ Update available
b    github.com/shurcooL/go-goon/...
	b Non-default branch checked out
 *   github.com/shurcooL/Conception-go/...
	* Uncommited changes in working dir
  #  github.com/russross/blackfriday/...
	# Remote path doesn't match import path
   $ github.com/microcosm-cc/bluemonday/...
	$ Stash exists
  /  github.com/go-forks/go-pkg-xmlx/...
	/ Remote repository not found (was it deleted? made private?):
		remote repository not found:
		exit status 128: remote: Repository not found.
		fatal: repository 'https://github.com/go-forks/go-pkg-xmlx/' not found
```

There are a few observations that can be made from that sample output.

-	`uniuri` and `goleveldb` repos are ***out of date***, I should update them via `go get -u`.
-	`go-goon` repo has a ***non-default*** branch checked out, I should be aware of that.
-	`Conception-go` repo has ***uncommited changes***. I should remember to commit or discard the changes.
-	`blackfriday` repo has a ***remote that doesn't match its import path***. It's likely my fork in place of the original repo for temporary development purposes.
-	`bluemonday` repo has a ***stash***. Perhaps I have some unfinished and uncommited work that I should take care of.
-	All other repos are ***up to date*** and looking good (they're not displayed unless `-v` is used).

Directories
-----------

| Path                                                            | Synopsis                                                                                          |
|-----------------------------------------------------------------|---------------------------------------------------------------------------------------------------|
| [status](https://godoc.org/github.com/shurcooL/gostatus/status) | Package status provides a func to check if two repo URLs are equal in the context of Go packages. |

License
-------

-	[MIT License](LICENSE)
