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

	matcher := globerous.EqualMatch(
		"testdata",
		globerous.AnyRecursiveMatcher(
			globerous.HasSuffixMatcher(".txt", nil),
		),
	)

	wd, err := os.Getwd()
	must(err)
	dir := filepath.Join(wd, "examples", "hybridexample")

	fmt.Println("Matches:")
	must(globerous.WalkSimple(fs, matcher, dir, func(dir string, info os.FileInfo) error {
		fmt.Println(" ", filepath.Join(dir, info.Name()))
		return nil
	}))
}
