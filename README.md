outdated
========

A command line tool that lists Go packages with newer versions.

Installation
------------

```bash
$ mkdir /tmp/temp_dir_for_download
$ GOPATH=/tmp/temp_dir_for_download go get github.com/shurcooL/outdated
# Copy /tmp/temp_dir_for_download/bin/outdated to somewhere in your $PATH
# Feel free to delete /tmp/temp_dir_for_download
```

Usage
-----

```bash
# TODO: Implement this
$ outdated [packages]

$ [packages] | outdated
```

Examples
--------

```bash
# Run outdated on all your packages
$ go list all | outdated

# Run on all dependencies (recursive) of package in current working directory
$ go list -f '{{join .Deps "\n"}}' . | outdated

# Run on all dependencies (recursive) of specified package
$ go list -f '{{join .Deps "\n"}}' github.com/some/import/path | outdated
```
