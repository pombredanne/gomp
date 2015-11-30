// gomp lists Go dependencies parsing import paths.
package main

import (
	"fmt"
	"os"
	pathpkg "path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/gyuho/gomp/walk"
	"github.com/spf13/cobra"
)

const (
	cliName        = "gomp"
	cliDescription = "gomp lists Go dependencies parsing import paths."
)

// GlobalFlags contains all the flags defined globally
// and that are to be inherited to all sub-commands.
type GlobalFlags struct {
	// GorootPath is the Go root path. This is required to
	// get all the list of standard Go packages.
	GorootPath string

	// OutputPath is the filepath to store the output.
	OutputPath string

	// ShowOnlyExternal is true when you want to list only non-standard packages.
	ShowOnlyExternal bool

	// IgnoreBuild is true when you want to ignore build tags.
	// By default, gomp already ignores build contraints by platform.
	// If this is true, it handles edges cases like `+build ignore` or
	// 'appengine'.
	IgnoreBuild bool
}

var (
	globalFlags = GlobalFlags{}

	rootCmd = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"goomp", "gom", "gmop"},
		Example:    "'gomp -o imports.txt .' lists all dependencies in the imports.txt file.",
		RunE:       rootCommandFunc,
	}
)

func init() {
	// https://github.com/golang/go/blob/master/src/go/build/build.go#L292
	goRoot := pathpkg.Clean(runtime.GOROOT())
	rootCmd.PersistentFlags().StringVarP(&globalFlags.GorootPath, "goroot", "g", goRoot, "goroot is your GOROOT path. By default, it uses your runtime.GOROOT().")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.OutputPath, "output", "o", "", "output is the path to store the results. By default, it prints out to standard output.")
	rootCmd.PersistentFlags().BoolVarP(&globalFlags.ShowOnlyExternal, "show-external", "e", false, "show-external is true, then it only shows the external dependencies.")
	rootCmd.PersistentFlags().BoolVarP(&globalFlags.IgnoreBuild, "ignore-build", "i", false, "ignore-build is true, when you want to ignore build tags. By default, gomp already ginores build constraints by platform. If this is true, it handles edge cases like '+build ignore' or 'appenengine'.")
}

func init() {
	cobra.EnablePrefixMatching = true
}

func rootCommandFunc(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("gomp accepts only 0 or 1 argument but got %v", args)
	}

	var targetPath string
	if len(args) == 0 {
		targetPath = "."
	} else {
		targetPath = args[0]
	}

	if err := os.Chdir(targetPath); err != nil {
		return err
	}

	goPath := envOr("GOPATH", "")
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	projectPath, err := filepath.Rel(filepath.Join(goPath, "src"), pwd)
	if err != nil {
		return err
	}

	dmap, err := walk.Imports(globalFlags.ShowOnlyExternal, globalFlags.IgnoreBuild, targetPath)
	if err != nil {
		return err
	}

	slice := []string{}
	for k := range dmap {
		if strings.HasPrefix(k, projectPath) {
			continue
		}
		slice = append(slice, k)
	}
	sort.Strings(slice)

	txt := strings.Join(slice, "\n")
	if globalFlags.OutputPath == "" {
		fmt.Fprintln(os.Stdout, txt)
		return nil
	}
	if err := toFileWriteString(txt, globalFlags.OutputPath); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func toFileWriteString(txt, fpath string) error {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		// OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)
		f, err = os.Create(fpath)
		if err != nil {
			return err
		}
	}
	defer f.Close()
	if _, err := f.WriteString(txt); err != nil {
		return err
	}
	return nil
}

// https://github.com/golang/go/blob/master/src/go/build/build.go#L320
func envOr(name, def string) string {
	s := os.Getenv(name)
	if s == "" {
		return def
	}
	return s
}
