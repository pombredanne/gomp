package gomp

import "testing"

func TestGetStdPkg(t *testing.T) {
	rmap, err := GetStdPkg("/usr/local/go")
	if err != nil {
		t.Fatal(err)
	}
	for key := range rmap {
		t.Logf("%s", key)
	}
}

func TestGetImports(t *testing.T) {
	rmap, err := GetImports(".")
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range rmap {
		t.Logf("%s | %+v\n", k, v)
	}
}

func TestGetNonStdImports(t *testing.T) {
	rmap, err := GetNonStdImports("/usr/local/go", ".")
	if err != nil {
		t.Fatal(err)
	}
	for key := range rmap {
		t.Logf("%s", key)
	}
}
