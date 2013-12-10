outdated
========

A command line tool that lists Go packages with newer versions.

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
$ mkdir /tmp/new-dl-dir && GOPATH=/tmp/new-dl-dir go get github.com/shurcooL/outdated
```

Copy `/tmp/new-dl-dir/bin/outdated` to somewhere in your `$PATH`. Feel free to delete `/tmp/new-dl-dir`.

Usage
-----

```bash
$ [newline separated packages] | outdated

# TODO: Implement this
#$ outdated [packages]
```

Examples
--------

```bash
# Run outdated on all your packages
$ go list all | outdated

# Run outdated on specified package
$ go list github.com/import/path | outdated

# Run outdated on package in current working dir
$ go list . | outdated

# Run on all dependencies (recursive) of specified package
$ go list -f '{{join .Deps "\n"}}' github.com/import/path | outdated

# Run on all dependencies (recursive) of package in current working dir
$ go list -f '{{join .Deps "\n"}}' . | outdated
```
