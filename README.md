gostatus
========

A command line tool that shows the status of Go packages.

Legend:
- `@` - Git repo
- `b` - Non-master branch checked out
- `*` - Uncommited changes in working dir
- `+` - Update available (latest remote revision doesn't match local revision)
- `/` - Command (package main)

Caveat: It currently prints remote version information only for git repositories. Mercurial support to be done...

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
