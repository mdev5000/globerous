package globerous_test

import (
	"fmt"
	"github.com/mdev5000/globerous"
	"os"
	"path/filepath"
	"strings"
)

func Example_simpleExample() {
	path, err := filepath.Abs(filepath.Join("testdata", "examplesfs"))
	if err != nil {
		panic(err)
	}
	fs := globerous.NewOSGlobFs()

	compiler := globerous.NewCompiler(globerous.GlobPlusPartCompiler)
	matcher := compiler.MustCompile("*", "*", "*.txt")

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Walk matched files and folders.
	err = globerous.WalkSimple(fs, matcher, path, func(dir string, info os.FileInfo) error {
		fmt.Println(strings.TrimPrefix(dir, wd), info.Name())
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Output:
	///testdata/examplesfs/first/nested first.txt
	///testdata/examplesfs/second/nested second.txt
}
