package walk

import (
	"os"
	"path/filepath"
	"strings"
)

// walkExt returns all FileInfos with specific extension.
// Make sure to prefix the extension name with dot.
// For example, to find all go files, pass ".go".
func walkExt(targetDir, ext string) (map[os.FileInfo]string, error) {
	rmap := make(map[os.FileInfo]string)
	visit := func(path string, f os.FileInfo, err error) error {
		if f != nil {
			if !f.IsDir() {
				if filepath.Ext(path) == ext {
					if !filepath.HasPrefix(path, ".") && !strings.Contains(path, "/.") {
						if _, ok := rmap[f]; !ok {
							wd, err := os.Getwd()
							if err != nil {
								return err
							}
							thepath := filepath.Join(wd, strings.Replace(path, wd, "", -1))
							rmap[f] = thepath
						}
					}
				}
			}
		}
		return nil
	}
	err := filepath.Walk(targetDir, visit)
	if err != nil {
		return nil, err
	}
	return rmap, nil
}

// walkDir returns all directories.
func walkDir(targetDir string) (map[os.FileInfo]string, error) {
	rmap := make(map[os.FileInfo]string)
	visit := func(path string, f os.FileInfo, err error) error {
		if f != nil {
			if f.IsDir() {
				if !filepath.HasPrefix(path, ".") && !strings.Contains(path, "/.") {
					if _, ok := rmap[f]; !ok {
						rmap[f] = filepath.Join(targetDir, path)
					}
				}
			}
		}
		return nil
	}
	if err := filepath.Walk(targetDir, visit); err != nil {
		return nil, err
	}
	return rmap, nil
}
