package globerous_test

import (
	"fmt"
	"github.com/mdev5000/globerous"
)

func ExampleMaxDepthMatcher() {
	compiler := globerous.NewCompiler(globerous.GlobPlusPartCompiler)
	matcher := compiler.MustCompile("**", "test.txt")
	list, err := globerous.List(globerous.NewOSGlobFs(), globerous.MaxDepthMatcher(4, matcher), "/")
	if err != nil {
		panic(err)
	}
	fmt.Println(list)
}

func ExampleMultiMatcher() {
	compiler := globerous.NewCompiler(globerous.GlobPlusPartCompiler)
	matcher := globerous.MultiMatcher(
		compiler.MustCompile("**", "test", "*.go"),
		compiler.MustCompile("**", "tests", "*.go"),
	)
	files, err := globerous.List(globerous.NewOSGlobFs(), matcher, "/")
	if err != nil {
		panic(err)
	}
	fmt.Println(files)
}

func ExampleFileOrFolderMatcher_matchFilesOnly() {
	compiler := globerous.NewCompiler(globerous.GlobPlusPartCompiler)
	// Only match files matching the pattern **/test.txt
	matcher := globerous.FileOrFolderMatcher(
		false,
		compiler.MustCompile("**", "test.txt"),
	)
	list, err := globerous.List(globerous.NewOSGlobFs(), matcher, "/")
	if err != nil {
		panic(err)
	}
	fmt.Println(list)
}
