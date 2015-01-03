package walk

import (
	"os"
	"testing"
)

func TestStdPkg(t *testing.T) {
	t.Log("GOROOT:", os.Getenv("GOROOT"))
	// if it cannot find this directory
	// it will find the GOROOT environment variable
	rmap, err := StdPkg("/usr/local/go")
	if err != nil {
		t.Fatal(err)
	}

	// this may not be accurate for gvm install set-up like Travis CI
	for key := range rmap {
		t.Logf("%s", key)
	}
}

func TestImports(t *testing.T) {
	t.Log("GOROOT:", os.Getenv("GOROOT"))
	rmap, err := Imports(".")
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range rmap {
		t.Logf("%s | %+v\n", k, v)
	}
}

func TestNonStdImports(t *testing.T) {
	t.Log("GOROOT:", os.Getenv("GOROOT"))
	rmap, err := NonStdImports("/usr/local/go", ".")
	if err != nil {
		t.Fatal(err)
	}
	for key := range rmap {
		t.Logf("%s", key)
	}
}
