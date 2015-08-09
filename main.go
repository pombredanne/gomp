// gomp is a command line tool for listing imported (non-standard) packages, like pip-freeze in Python.
//
// Example:
//	go get -v github.com/gyuho/gomp // to install
//	gomp -h // to see the manual page
//	gomp -target=./go/src/github.com/username/project
//	// This will extracts the list of all external packages in the project directory excluding Go standard packages.
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/gyuho/gomp/walk"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	goRootPathPtr := flag.String("goroot", runtime.GOROOT(), "Specify your GOROOT path. Usually set as /usr/local/go; Default value is set as runtime.GOROOT()")
	targetPathPtr := flag.String("target", pwd, "Specify the target path you want to extract from. Default value is set as os.Getwd()")
	outputPathPtr := flag.String("output", filepath.Join(pwd, "imports.txt"), "Specify the output file path. Default value is set as imports.txt at os.Getwd()")
	flag.Parse()

	rmap, err := walk.NonStdImports(*goRootPathPtr, *targetPathPtr)
	if err != nil {
		log.Fatal(err)
	}

	imap := make(map[string]bool)
	for key := range rmap {
		if _, ok := imap[key]; !ok {
			imap[key] = true
		}
	}

	slice := []string{}
	for imp := range imap {
		slice = append(slice, imp)
	}

	sort.Strings(slice)

	if err := fromLines(slice, *outputPathPtr); err != nil {
		log.Fatal(err)
	}
}

func fromLines(lines []string, fpath string) error {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		f, err = os.Create(fpath)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	// func NewWriter(w io.Writer) *Writer
	wr := bufio.NewWriter(f)

	for _, line := range lines {
		// func Fprintln(w io.Writer, a ...interface{}) (n int, err error)
		fmt.Fprintln(wr, line)
	}

	if err := wr.Flush(); err != nil {
		return err
	}
	return nil
}
