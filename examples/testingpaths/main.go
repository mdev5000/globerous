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
	dir := filepath.Join(wd, "examples", "testingpaths")

	// Match based on the filesystem info.
	path, err := globerous.GetPathFromString(fs, dir, "testdata/somefolder/first.txt")
	must(err)
	match, partialMatch, err := globerous.MatchesPath(matcher, path)
	must(err)
	fmt.Println("via path:")
	fmt.Println(" match:", match)
	fmt.Println(" partialMatch:", partialMatch)
	fmt.Println("")

	// Match based on string. This is fine as long as your matches don't require extra file information (ex. filesize).
	path = globerous.FakePath("testdata/somefolder/first.txt", false)
	match, partialMatch, err = globerous.MatchesPath(matcher, path)
	must(err)
	fmt.Println("via fake path:")
	fmt.Println(" match:", match)
	fmt.Println(" partialMatch:", partialMatch)
	fmt.Println("")
}
