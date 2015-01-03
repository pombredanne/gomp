package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/gyuho/filex"
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

	if err := filex.WriteLines(*outputPathPtr, slice); err != nil {
		log.Fatal(err)
	}
}
