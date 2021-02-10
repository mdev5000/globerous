package globerous_test

import (
	"fmt"
	"github.com/mdev5000/globerous"
)

// Resolve and then test that a path matches the matcher.
func ExampleMatchesPath_matchingAPathBasedOnTheOs() {
	compiler := globerous.NewCompiler(globerous.GlobPlusPartCompiler)
	matcher := compiler.MustCompile("some", "*.txt")
	path, err := globerous.GetPathFromString(globerous.NewOSGlobFs(), "/some/root", "some/path.txt")
	if err != nil {
		panic(err)
	}
	match, partialMatch, err := globerous.MatchesPath(matcher, path)
	if err != nil {
		panic(err)
	}
	fmt.Println(match)
	fmt.Println(partialMatch)
}

// You can match based on a just a path strings instead of retrieving the actual path information. This is fine as long
// as none of your matchers require information other than IsDir() and Name() from the io.FileInfo (ex. the file size).
func ExampleMatchesPath_matchingAStringPath() {
	compiler := globerous.NewCompiler(globerous.GlobPlusPartCompiler)
	matcher := compiler.MustCompile("first", "*", "test.txt")
	path := globerous.FakePath("first/second/test.txt", false)
	match, partialMatch, err := globerous.MatchesPath(matcher, path)
	if err != nil {
		panic(err)
	}
	fmt.Println(match)
	fmt.Println(partialMatch)
}
