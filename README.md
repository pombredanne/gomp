gomp [![Build Status](https://travis-ci.org/gyuho/gomp.svg?branch=master)](https://travis-ci.org/gyuho/gomp) [![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/gyuho/gomp)
==========


`go get -v github.com/gyuho/gomp`


```
gomp lists Go dependencies parsing import paths.

Usage:
  gomp [flags]

Examples:
'gomp -o imports.txt .' lists all dependencies in the imports.txt file.

Flags:
  -g, --goroot="/usr/local/go": goroot is your GOROOT path. By default, it uses your runtime.GOROOT().
  -o, --output="": output is the path to store the results. By default, it prints out to standard output.

```
