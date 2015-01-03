package gomp

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"go/parser"
	"go/token"

	"github.com/gyuho/filex"
)

// GetStdPkg returns all lists of Go standard packages.
// Usually pass "/usr/local/go".
// There is an alternative way: https://github.com/golang/tools/blob/master/imports/mkstdlib.go
func GetStdPkg(goRootPath string) (map[string]bool, error) {
	if goRootPath == "" {
		goRootPath = os.Getenv("GOROOT")
		if goRootPath == "" {
			return nil, errors.New("can't find GOROOT: try to set it to /usr/local/go")
		}
	}
	stdpkgPath := filepath.Join(goRootPath, "src")
	rmap, err := filex.WalkDir(stdpkgPath)
	if err != nil {
		goRootPath = os.Getenv("GOROOT")
		stdpkgPath = filepath.Join(goRootPath, "src")
		rmap, err = filex.WalkDir(stdpkgPath)
		if err != nil {
			return nil, err
		}
	}
	smap := make(map[string]bool)
	for _, val := range rmap {
		stdName := strings.Replace(val, stdpkgPath, "", -1)
		stdName = filepath.Clean(stdName)
		if strings.HasPrefix(stdName, "/") {
			stdName = stdName[1:]
		}
		if strings.HasPrefix(stdName, "cmd") {
			continue
		}
		if strings.Contains(stdName, "testdata") {
			continue
		}
		if strings.Contains(stdName, "internal") {
			continue
		}
		if len(stdName) < 2 {
			continue
		}
		if _, ok := smap[stdName]; !ok {
			smap[stdName] = true
		}
	}
	return smap, nil
}

// GetImports gets all import paths from Go source code.
// It returns the map from import path string to the file paths.
func GetImports(targetDir string) (map[string][]string, error) {
	rmap, err := filex.WalkExt(targetDir, ".go")
	if err != nil {
		return nil, err
	}
	fmap := make(map[string][]string)
	var mutex sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(rmap))
	for _, fpath := range rmap {
		go func(targetcode string) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, targetcode, nil, parser.ImportsOnly)
			if err != nil {
				log.Fatal(err)
			}
			for _, elem := range f.Imports {
				mutex.Lock()
				pathValue := elem.Path.Value
				pathValue = strings.Replace(pathValue, `"`, "", -1)
				pathValue = strings.TrimSpace(pathValue)
				if _, ok := fmap[pathValue]; !ok {
					fmap[pathValue] = []string{targetcode}
				} else {
					fmap[pathValue] = append(fmap[pathValue], targetcode)
				}
				mutex.Unlock()
			}
			wg.Done()
		}(fpath)
	}
	wg.Wait()
	return fmap, nil
}

// GetNonStdImports get all import paths that are not Go standard package.
func GetNonStdImports(goRootPath, targetDir string) (map[string][]string, error) {
	rmap, err := GetImports(targetDir)
	if err != nil {
		return nil, err
	}
	stdMap, err := GetStdPkg(goRootPath)
	if err != nil {
		return nil, err
	}
	for k := range rmap {
		if _, ok := stdMap[k]; ok {
			delete(rmap, k)
		}
	}
	return rmap, nil
}
