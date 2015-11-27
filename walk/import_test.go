package walk

import (
	"fmt"
	pathpkg "path"
	"runtime"
	"testing"
)

func TestImports(t *testing.T) {
	rmap, err := Imports(false, true, ".")
	if err != nil {
		t.Error(err)
	}
	for k, v := range rmap {
		t.Logf("Imports: %s | %+v\n", k, v)
	}
}

func TestAPIFiles(t *testing.T) {
	goroot := pathpkg.Clean(runtime.GOROOT())
	am, err := apiFiles(goroot)
	if err != nil {
		t.Error(err)
	}
	for k := range am {
		t.Logf("am: %s\n", k)
	}
}

func TestStdPkg(t *testing.T) {
	goroot := pathpkg.Clean(runtime.GOROOT())
	sm, err := StdPkg(goroot)
	if err != nil {
		t.Error(err)
	}
	for k := range sm {
		_ = k
		// fmt.Printf("sm: %s\n", k)
	}
	nm, err := StdPkgNaive(goroot)
	if err != nil {
		t.Error(err)
	}
	for k := range nm {
		_ = k
		// fmt.Printf("sm: %s\n", k)
	}

	t.Logf("StdPkg: %d\n", len(sm))
	t.Logf("StdPkgNaive: %d\n", len(nm))

	for k := range nm {
		if _, ok := sm[k]; !ok {
			fmt.Println(k, "is not in StdPkg")
		}
	}

	for k := range sm {
		if _, ok := nm[k]; !ok {
			fmt.Println(k, "is not in StdPkgNaive")
		}
	}

}
