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

	compiler := globerous.NewCompiler(globerous.HybridGlobRegexPartCompiler)

	// If a part starts with ^ or ends with $ then it is handled as regex instead of as a glob.
	matcher := compiler.MustCompile("testdata", "^another|somefolder$", "*.txt")

	wd, err := os.Getwd()
	must(err)
	dir := filepath.Join(wd, "examples", "hybridexample")

	fmt.Println("Matches:")
	must(globerous.WalkSimple(fs, matcher, dir, func(dir string, info os.FileInfo) error {
		fmt.Println(" ", filepath.Join(dir, info.Name()))
		return nil
	}))
}
