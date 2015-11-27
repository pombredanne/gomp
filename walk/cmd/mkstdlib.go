package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	pathpkg "path"
	"runtime"
	"sort"

	"github.com/gyuho/gomp/walk"
)

func main() {
	goroot := pathpkg.Clean(runtime.GOROOT())
	sm, err := walk.StdPkg(goroot)
	if err != nil {
		panic(err)
	}
	nm, err := walk.StdPkgNaive(goroot)
	if err != nil {
		panic(err)
	}

	fmap := make(map[string]struct{})
	for k := range sm {
		fmap[k] = struct{}{}
	}
	for k := range nm {
		fmap[k] = struct{}{}
	}
	fs := []string{}
	for k := range fmap {
		fs = append(fs, k)
	}
	sort.Strings(fs)

	var buf bytes.Buffer
	outf := func(format string, args ...interface{}) {
		fmt.Fprintf(&buf, format, args...)
	}
	outf("// AUTO-GENERATED BY mkstdlib.go\n\n")
	outf("package walk\n\n")
	outf("var StandardPackageMap = map[string]struct{}{\n")
	for _, v := range fs {
		outf("\t%q: struct{}{},\n", v)
	}
	outf("}\n")
	fmtbuf, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	f, err := openToOverwrite("../pkg.go")
	if err != nil {
		panic(err)
	}
	f.Write(fmtbuf)
}

// openToOverwrite creates or opens a file for overwriting.
// Make sure to close the file.
func openToOverwrite(fpath string) (*os.File, error) {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		// OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)
		f, err = os.Create(fpath)
		if err != nil {
			return f, err
		}
	}
	return f, nil
}
