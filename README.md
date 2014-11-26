gostatus
========

A command line tool that shows the status of (many) Go package repositories.

Installation
------------

```bash
go get -u github.com/shurcooL/gostatus
```

Usage
-----

```bash
Usage: gostatus [flags] [packages]
       [newline separated packages] | gostatus --stdin [flags]
  -debug=false: Give the output with verbose debug information.
  -plumbing=false: Give the output in an easy-to-parse format for scripts.
  -stdin=false: Read the list of newline separated Go packages from stdin.
  -v=false: Verbose output: show all Go packages, not just ones with notable status.

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
```

Examples
--------

```bash
# Show status of all your packages
$ gostatus all

# Show status of package in current directory
$ gostatus .

# Show status of specified package
$ gostatus some/import/path

# Show status of all dependencies (recursive) of package in current working dir
$ go list -f '{{join .Deps "\n"}}' . | gostatus --stdin -v

# Show status of all dependencies (recursive) of specified package
$ go list -f '{{join .Deps "\n"}}' some/import/path | gostatus --stdin -v

# Run go get -u on all reported updates
$ gostatus -plumbing all | sed '/^@  +. /!d; s///' | xargs go get -u
```

Sample Output
-------------

```bash
$ gostatus all
  +  github.com/dchest/uniuri/...
  +  github.com/syndtr/goleveldb/...
b    github.com/shurcooL/go-goon/...
 *   github.com/shurcooL/Conception-go/...
  #  github.com/russross/blackfriday/...
   $ github.com/microcosm-cc/bluemonday/...
```

There are a few observations that can be made from that sample output.

- `uniuri` and `goleveldb` repos are ***out of date***, I should update them via `go get -u`.
- `go-goon` repo has a ***non-default*** branch checked out, I should be aware of that.
- `Conception-go` repo has ***uncommited changes***. I should remember to commit or discard the changes.
- `blackfriday` repo has a ***remote that doesn't match its import path***. It's likely my fork in place of the original repo for temporary development purposes.
- `bluemonday` repo has a ***stash***. Perhaps I have some unfinished and uncommited work that I should take care of.
- All other repos are ***up to date*** and looking good (they're not displayed unless `-v` is used).
