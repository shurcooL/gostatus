gostatus
========

A command line tool that shows the status of Go packages.

Legend:

- `@` - Git repo
- `b` - Non-master branch checked out
- `*` - Uncommited changes in working dir
- `+` - Update available (latest remote revision doesn't match local revision)
- `/` - Command (package main)

Caveats:

- It currently prints remote version information only for git repositories. Mercurial support to be done...
- It currently lists one Go package per git repo (even if there are many), in order to avoid polling same git repo more than once. A proper solution will be to cache and reuse the results, to be done.

Installation
------------

```bash
$ mkdir /tmp/new-dl-dir && GOPATH=/tmp/new-dl-dir go get github.com/shurcooL/gostatus
```

Copy `/tmp/new-dl-dir/bin/gostatus` to somewhere in your `$PATH`. Feel free to delete `/tmp/new-dl-dir`.

Usage
-----

```bash
$ [newline separated packages] | gostatus

# TODO: Consider implementing this
#$ gostatus [packages]
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
$ go list all | gostatus
      code.google.com/p/go-uuid/uuid
      code.google.com/p/snappy-go/snappy
      code.google.com/p/goprotobuf/proto
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

- My `github.com/dchest/uniuri` and `github.com/syndtr/goleveldb/leveldb` packages are out of date, I should update them via `go get -u`.
- My `github.com/shurcooL/go-goon` package has a branch other than "main" checked out, I should be aware of that.
- My `github.com/shurcooL/Conception-go` package has a dirty working tree. I should remember to commit or discard the changes.
- All other packages are up to date and looking good.
