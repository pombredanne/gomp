package walk

import "testing"

func TestWalkExt(t *testing.T) {
	rmap, err := walkExt(".", ".go")
	if err != nil {
		t.Error(err)
	}
	if len(rmap) != 8 {
		t.Errorf("expected to have total 8 go files here but got %d", len(rmap))
	}
}

func TestWalkDir(t *testing.T) {
	rmap, err := walkDir(".")
	if err != nil {
		t.Error(err)
	}
	if len(rmap) != 3 {
		t.Errorf("expected to have 3 sub-directories but got %d", len(rmap))
	}
}
