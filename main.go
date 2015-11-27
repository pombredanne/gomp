// gomp is a simple command line tool for go import paths.
package main

import (
	"fmt"
	"os"
	pathpkg "path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/gyuho/gomp/walk"
	"github.com/spf13/cobra"
)

const (
	cliName            = "gomp"
	cliDescription     = "A simple command line tool for go import paths."
	cliDescriptionLong = `gomp can list all non-standard packages in your projects.
This can be useful for checking all the external dependencies.
`
)

// GlobalFlags contains all the flags defined globally
// and that are to be inherited to all sub-commands.
type GlobalFlags struct {
	// GorootPath is the Go root path. This is required to
	// get all the list of standard Go packages.
	GorootPath string

	// OutputPath is the filepath to store the output.
	OutputPath string
}

var (
	tabOut      *tabwriter.Writer
	globalFlags = GlobalFlags{}

	rootCmd = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		Long:       cliDescriptionLong,
		SuggestFor: []string{"goomp", "gom", "gmop"},
		Example:    "'gomp -o imports.txt .' lists all the external dependencies in imports.txt file.",
		RunE:       rootCommandFunc,
	}
)

func init() {
	// https://github.com/golang/go/blob/master/src/go/build/build.go#L292
	goRoot := pathpkg.Clean(runtime.GOROOT())
	rootCmd.PersistentFlags().StringVarP(&globalFlags.GorootPath, "goroot", "g", goRoot, "goroot is your GOROOT path. By default, it uses your runtime.GOROOT().")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.OutputPath, "output", "o", "", "output is the path to store the results. By default, it prints out to standard output.")
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

	dmap, err := walk.NonStdImports(globalFlags.GorootPath, targetPath)
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
