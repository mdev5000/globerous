package main

import (
	"fmt"
	"github.com/mdev5000/globerous"
	"os"
	"path/filepath"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fs := globerous.NewOSGlobFs()

	compiler := globerous.NewCompiler(globerous.GlobPlusPartCompiler)
	matcher := compiler.MustCompile("testdata", "*", "*.txt")

	wd, err := os.Getwd()
	must(err)
	dir := filepath.Join(wd, "examples", "simpleglob")

	fmt.Println("Matches:")
	must(globerous.WalkSimple(fs, matcher, dir, func(dir string, info os.FileInfo) error {
		fmt.Println(" ", filepath.Join(dir, info.Name()))
		return nil
	}))
}
