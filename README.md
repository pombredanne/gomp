gomp [![Build Status](https://travis-ci.org/gyuho/gomp.svg?branch=master)](https://travis-ci.org/gyuho/gomp) [![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/gyuho/gomp)
==========


`go get -v github.com/gyuho/gomp`


```
gomp can list all non-standard packages in your projects.
This can be useful for checking all the external dependencies.

Usage:
  gomp [flags]

Examples:
'gomp -o imports.txt .' lists all the external dependencies in imports.txt file.

Flags:
  -g, --goroot="/usr/local/go": goroot is your GOROOT path. By default, it uses your runtime.GOROOT().
  -o, --output="": output is the path to store the results. By default, it prints out to standard output.

```
