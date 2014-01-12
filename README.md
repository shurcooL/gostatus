gostatus
========

A command line tool that shows the status of (many) Go packages.

Installation
------------

```bash
$ mkdir /tmp/gostatus/ && GOPATH=/tmp/gostatus/ go get github.com/shurcooL/gostatus
```

Copy `/tmp/gostatus/bin/gostatus` to somewhere in your `$PATH`. Feel free to remove `/tmp/gostatus/`.

Usage
-----

```bash
Usage: [newline separated packages] | gostatus [--all] [--plumbing]
  -all=false: Show all Go packages, not just ones with notable status.
  -plumbing=false: Give the output in an easy-to-parse format for scripts.

Examples:
  # Show status of packages with notable status
  go list all | gostatus

  # Show status of all dependencies (recursive) of package in cur working dir
  go list -f '{{join .Deps "\n"}}' . | gostatus --all

Legend:
  @ - Vcs repo
  b - Non-master branch checked out
  * - Uncommited changes in working dir
  + - Update available (latest remote revision doesn't match local revision)
  / - Command (package main)
```

Examples
--------

```bash
# Show status of all your packages
$ go list all | gostatus

# Show status of package in current working dir
$ go list . | gostatus

# Show status of specified package
$ echo some/import/path | gostatus

# Show status of all dependencies (recursive) of package in current working dir
$ go list -f '{{join .Deps "\n"}}' . | gostatus

# Show status of all dependencies (recursive) of specified package
$ go list -f '{{join .Deps "\n"}}' some/import/path | gostatus
```

Sample Output
-------------

```bash
$ go list all | gostatus --all
@     code.google.com/p/go-uuid/uuid
@     code.google.com/p/snappy-go/snappy
@     code.google.com/p/goprotobuf/proto
@     github.com/jmhodges/levigo
@     github.com/bradfitz/gomemcache/memcache
@  +  github.com/dchest/uniuri
@     github.com/bmizerany/assert
@  +  github.com/syndtr/goleveldb/leveldb
@     github.com/vmihailenco/bufio
@     tux21b.org/v1/gocql
@     github.com/Ysgard/GoGLutils
@   / github.com/chsc/gogl
@     github.com/Jragonmiris/mathgl
@     github.com/ftrvxmtrx/tga
@     github.com/davecheney/profile
@     github.com/go-gl/glfw3
@     github.com/howeyc/fsnotify
@     github.com/pkg/math
@     github.com/russross/blackfriday
@     github.com/davecgh/go-spew/spew
@     github.com/sergi/go-diff/diffmatchpatch
@b    github.com/shurcooL/go-goon
@   / github.com/shurcooL/goe
@ * / github.com/shurcooL/Conception-go
@     github.com/shurcooL/goglu/glu21
@   / github.com/shurcooL/gostatus
@     honnef.co/go/importer
```

There are a few observations that can be made from that sample output.

- `uniuri` and `leveldb` packages are ***out of date***, I should update them via `go get -u`.
- `go-goon` package has a branch other than ***master*** checked out, I should be aware of that.
- `Conception-go` package has ***uncommited changes***. I should remember to commit or discard the changes.
- All other packages are ***up to date*** and looking good.

Caveats
-------

- It currently lists one Go package per repo (even if there are many), in order to avoid polling same repo more than once. A proper solution will be to cache and reuse the results, to be done.
