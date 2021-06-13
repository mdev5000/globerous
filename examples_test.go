package globerous_test

import (
	"fmt"
	"github.com/mdev5000/globerous"
	"os"
	"path/filepath"
	"strings"
)

func Example_overview() {
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
		// Do something with the files.
		return nil
	})
	if err != nil {
		panic(err)
	}

	// List matched files and folders.
	filesAndFolder, err := globerous.List(fs, matcher, path)
	if err != nil {
		panic(err)
	}
	fmt.Println("List: files and folder:")
	for _, f := range filesAndFolder {
		fmt.Println(strings.TrimPrefix(f, wd))
	}

	// List matched files only.
	files, err := globerous.List(fs, globerous.FileOrFolderMatcher(false, matcher), path)
	if err != nil {
		panic(err)
	}
	fmt.Println("List: files only:")
	for _, f := range files {
		fmt.Println(strings.TrimPrefix(f, wd))
	}

	// List files and folder matching up to a maximum depth of 4.
	files, err = globerous.List(
		fs,
		globerous.MaxDepthMatcher(4, compiler.MustCompile("**", "*.txt")),
		path,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Max Depth:")
	for _, f := range files {
		fmt.Println(strings.TrimPrefix(f, wd))
	}
	// Output:
	// List: files and folder:
	///testdata/examplesfs/first/nested/first.txt
	///testdata/examplesfs/second/nested/second.txt
	// List: files only:
	///testdata/examplesfs/first/nested/first.txt
	///testdata/examplesfs/second/nested/second.txt
	// Max Depth:
	///testdata/examplesfs/first/nested/first.txt
	///testdata/examplesfs/second/nested/second.txt
}

func Example_handlingErrors() {
	// @todo add example
}

// Use the MultiMatcher to combine 2 matching patterns into a single matcher. MultiMatcher acts as the union of 2
// matchers.
func Example_multipleMatches() {
	compiler := globerous.NewCompiler(globerous.GlobPlusPartCompiler)
	matcher := globerous.MultiMatcher(
		compiler.MustCompile("**", "test", "*.go"),
		compiler.MustCompile("**", "tests", "*.go"),
	)
	matches, err := globerous.List(globerous.NewOSGlobFs(), matcher, "/")
	if err != nil {
		panic(err)
	}
	fmt.Println(matches)
}

// You can match based on raw matchers and insert your own where needed.
func Example_usingRawMatchers() {
	// Equivalent to testdata/**/*.txt.
	matcher := globerous.EqualMatch(
		"testdata",
		globerous.AnyRecursiveMatcher(
			globerous.HasSuffixMatcher(".txt", nil),
		),
	)

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir := filepath.Join(wd, "examples", "hybridExample")

	fmt.Println("Matches:")
	err = globerous.WalkSimple(globerous.NewOSGlobFs(), matcher, dir, func(dir string, info os.FileInfo) error {
		fmt.Println(" ", filepath.Join(dir, info.Name()))
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func Example_testingPaths() {
	fs := globerous.NewOSGlobFs()
	compiler := globerous.NewCompiler(globerous.GlobPlusPartCompiler)
	matcher := compiler.MustCompile("**", "*.txt")

	// Match based on the filesystem info.
	path, err := globerous.GetPathFromString(fs, "/", "testdata/some_folder/first.txt")
	if err != nil {
		panic(err)
	}
	match, partialMatch, err := globerous.MatchesPath(matcher, path)
	if err != nil {
		panic(err)
	}
	fmt.Println("via path:")
	fmt.Println(" match:", match)
	fmt.Println(" partialMatch:", partialMatch)
	fmt.Println("")

	// Match based on string. This is fine as long as your matches don't require extra file information (ex. filesize).
	path = globerous.FakePath("testdata/some_folder/first.txt", false)
	match, partialMatch, err = globerous.MatchesPath(matcher, path)
	if err != nil {
		panic(err)
	}
	fmt.Println("via fake path:")
	fmt.Println(" match:", match)
	fmt.Println(" partialMatch:", partialMatch)
	fmt.Println("")
}
