package gomp

import (
	"os"
	"testing"
)

func TestGetStdPkg(t *testing.T) {
	t.Log("GOROOT:", os.Getenv("GOROOT"))
	// if it cannot find this directory
	// it will find the GOROOT environment variable
	rmap, err := GetStdPkg("/usr/local/go")
	if err != nil {
		t.Fatal(err)
	}

	// this may not be accurate for gvm install set-up like Travis CI
	for key := range rmap {
		t.Logf("%s", key)
	}
}

func TestGetImports(t *testing.T) {
	t.Log("GOROOT:", os.Getenv("GOROOT"))
	rmap, err := GetImports(".")
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range rmap {
		t.Logf("%s | %+v\n", k, v)
	}
}

func TestGetNonStdImports(t *testing.T) {
	t.Log("GOROOT:", os.Getenv("GOROOT"))
	rmap, err := GetNonStdImports("/usr/local/go", ".")
	if err != nil {
		t.Fatal(err)
	}
	for key := range rmap {
		t.Logf("%s", key)
	}
}
