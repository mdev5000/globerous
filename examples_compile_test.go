package globerous_test

import (
	"fmt"
	"github.com/mdev5000/globerous"
	"os"
	"path/filepath"
	"strings"
)

type myCustomMatcher struct {
	text string
	next globerous.Matcher
}

func (m myCustomMatcher) Matches(info os.FileInfo) (globerous.Matcher, bool, error) {
	if info.Name() == m.text {
		return m.next, true, nil
	}
	return nil, false, nil
}

func newMyCustomMatcher(text string, next globerous.Matcher) globerous.Matcher {
	return &myCustomMatcher{text, next}
}

func ExampleNewCompiler_globCompiler() {
	fs := globerous.NewOSGlobFs()

	compiler := globerous.NewCompiler(globerous.GlobPlusPartCompiler)
	matcher := compiler.MustCompile("some", "*", "test.txt")

	fmt.Println("Matches:")
	err := globerous.WalkSimple(fs, matcher, "/", func(dir string, info os.FileInfo) error {
		fmt.Println(" ", filepath.Join(dir, info.Name()))
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func ExampleNewCompiler_regexCompiler() {
	fs := globerous.NewOSGlobFs()

	compiler := globerous.NewCompiler(globerous.GlobPlusPartCompiler)
	matcher := compiler.MustCompile("^some$", "^.+$", "^.{3,10}\\.txt$")

	fmt.Println("Matches:")
	err := globerous.WalkSimple(fs, matcher, "/", func(dir string, info os.FileInfo) error {
		fmt.Println(" ", filepath.Join(dir, info.Name()))
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func ExampleNewCompiler_hybridCompiler() {
	fs := globerous.NewOSGlobFs()

	compiler := globerous.NewCompiler(globerous.HybridGlobRegexPartCompiler)

	// uses regex when a part starts with ^ or ends with $, otherwise uses glob syntax.
	matcher := compiler.MustCompile("^first|second$", "**", "tests", "*_test.py")

	fmt.Println("Matches:")
	err := globerous.WalkSimple(fs, matcher, "/", func(dir string, info os.FileInfo) error {
		fmt.Println(" ", filepath.Join(dir, info.Name()))
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func ExampleNewCompiler_customCompiler() {
	fs := globerous.NewOSGlobFs()

	compiler := globerous.NewCompiler(func(partStr string, next globerous.NextPartFn) (globerous.Matcher, error) {
		childMatcher, err := next()
		if err != nil {
			return nil, err
		}
		switch partStr {
		case "**":
			return globerous.AnyRecursiveMatcher(childMatcher), nil
		case "*":
			return globerous.AnyMatcher(childMatcher), nil
		}
		return newMyCustomMatcher(partStr, childMatcher), nil
	})

	matcher := compiler.MustCompile("some", "*", "test.txt")

	fmt.Println("Matches:")
	err := globerous.WalkSimple(fs, matcher, "/", func(dir string, info os.FileInfo) error {
		fmt.Println(" ", filepath.Join(dir, info.Name()))
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func ExampleMatchCompiler_Compile_matchingForwardSlash() {
	fs := globerous.NewOSGlobFs()

	compiler := globerous.NewCompiler(globerous.GlobPlusPartCompiler)

	// split the forward slashes before passing to the compiler.
	parts := strings.Split("some_folder/*/test.txt", "/")
	matcher := compiler.MustCompile(parts...)

	fmt.Println("Matches:")
	err := globerous.WalkSimple(fs, matcher, "/", func(dir string, info os.FileInfo) error {
		fmt.Println(" ", filepath.Join(dir, info.Name()))
		return nil
	})
	if err != nil {
		panic(err)
	}
}
