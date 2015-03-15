gomp [![Build Status](https://travis-ci.org/gyuho/gomp.svg?branch=master)](https://travis-ci.org/gyuho/gomp) [![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/gyuho/gomp)
==========

gomp is a command line tool for listing imported (non-standard) packages, like pip-freeze in Python.




## Install

`go get -v github.com/gyuho/gomp`

`gomp -h`:

```
Usage of gomp:
  -goroot="/usr/local/go": Specify your GOROOT path. Usually set as /usr/local/go; Default value is set as runtime.GOROOT()
  -output="/home/ubuntu/imports.txt": Specify the output file path. Default value is set as imports.txt at os.Getwd()
  -target="/home/ubuntu": Specify the target path you want to extract from. Default value is set as os.Getwd()
```


`gomp -target=./go/src/github.com/username/project` will extracts the list of all external packages in the project directory excluding Go standard packages.






<i>README.md Updated at 2015-03-14 21:22:09</i>
