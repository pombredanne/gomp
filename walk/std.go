package walk

import (
	"bufio"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// Imports gets all import paths from Go source code.
// It returns the map from import path string to the file paths.
func Imports(targetDir string) (map[string][]string, error) {
	rmap, err := walkExt(targetDir, ".go")
	if err != nil {
		return nil, err
	}
	rSize := len(rmap)

	fmap := make(map[string][]string)

	var mu sync.Mutex
	done, errCh := make(chan struct{}), make(chan error)

	for _, fpath := range rmap {

		go func(fpath string) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, fpath, nil, parser.ImportsOnly)
			if err != nil {
				errCh <- err
				return
			}
			for _, elem := range f.Imports {
				mu.Lock()
				pv := strings.TrimSpace(strings.Replace(elem.Path.Value, `"`, "", -1))
				if _, ok := fmap[pv]; !ok {
					fmap[pv] = []string{fpath}
				} else {
					fmap[pv] = append(fmap[pv], fpath)
				}
				mu.Unlock()
			}
			done <- struct{}{}
		}(fpath)

	}

	i := 0
	for {
		select {
		case e := <-errCh:
			return nil, e
		case <-done:
			i++
			if i == rSize {
				close(done)
				return fmap, nil
			}
		}
	}
}

// NonStdImports get all import paths that are not Go standard package.
func NonStdImports(goroot, targetDir string) (map[string][]string, error) {
	rmap, err := Imports(targetDir)
	if err != nil {
		return nil, err
	}
	for k := range rmap {
		if _, ok := GoStandardPackageMap[k]; ok {
			delete(rmap, k)
		}
	}
	return rmap, nil
}

func apiFiles(goroot string) (map[string]struct{}, error) {
	rmap, err := walkExt(filepath.Join(goroot, "api"), ".txt")
	if err != nil {
		return nil, err
	}
	am := make(map[string]struct{})
	for k := range rmap {
		base := filepath.Base(k.Name())
		if strings.HasPrefix(base, "go") {
			am[filepath.Join(goroot, "api", base)] = struct{}{}
		}
	}
	return am, nil
}

func mustOpen(name string) io.Reader {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	return f
}

var sym = regexp.MustCompile(`^pkg (\S+).*?, (?:var|func|type|const) ([A-Z]\w*)`)

// StdPkg returns all standard packages.
// Copied from https://github.com/golang/tools/blob/master/imports/mkstdlib.go.
func StdPkg(goroot string) (map[string]struct{}, error) {
	if goroot == "" {
		return nil, fmt.Errorf("got empty GOROOT")
	}

	am, err := apiFiles(goroot)
	if err != nil {
		return nil, err
	}

	rds := make([]io.Reader, 0)
	for k := range am {
		rds = append(rds, mustOpen(k))
	}
	f := io.MultiReader(rds...)
	pkgs := map[string]struct{}{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		l := sc.Text()
		has := func(v string) bool { return strings.Contains(l, v) }
		if has("struct, ") || has("interface, ") || has(", method (") {
			continue
		}
		if !strings.HasPrefix(l, "pkg ") {
			continue
		}
		if m := sym.FindStringSubmatch(l); m != nil {
			s := m[0]
			s = strings.Replace(s, "pkg ", "", -1)
			s = strings.Replace(s, ",", "", -1)
			s = strings.Split(s, " ")[0]
			pkgs[s] = struct{}{}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	return pkgs, nil
}

// StdPkgNaive returns all lists of Go standard packages.
// Alternatively look at https://github.com/golang/tools/blob/master/imports/mkstdlib.go
func StdPkgNaive(goroot string) (map[string]struct{}, error) {
	if goroot == "" {
		return nil, fmt.Errorf("got empty GOROOT")
	}

	stdpkgPath := filepath.Join(goroot, "src")
	rmap, err := walkDir(stdpkgPath)
	if err != nil {
		return nil, err
	}

	smap := make(map[string]struct{})
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
			smap[stdName] = struct{}{}
		}
	}

	return smap, nil
}
