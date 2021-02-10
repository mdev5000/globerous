package globerous_test

import (
	"fmt"
	"github.com/mdev5000/globerous"
	"os"
	"path/filepath"
)

func Example_simpleExample() {
	path, err := filepath.Abs(filepath.Join("testdata", "examplesfs"))
	if err != nil {
		panic(err)
	}
	fs := globerous.NewOSGlobFs()

	compiler := globerous.NewCompiler(globerous.GlobPlusPartCompiler)
	matcher := compiler.MustCompile("*", "*", "*.txt")

	// Walk matched files and folders.
	err = globerous.WalkSimple(fs, matcher, path, func(dir string, info os.FileInfo) error {
		fmt.Println(dir, info.Name())
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Output:
	///Users/matt/devtmp/go/globerous/testdata/examplesfs/first/nested first.txt
	///Users/matt/devtmp/go/globerous/testdata/examplesfs/second/nested second.txt
}
