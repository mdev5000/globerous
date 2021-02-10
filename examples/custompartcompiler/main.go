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

type myCustomMatcher struct {
	text string
	next globerous.Matcher
}

func (m myCustomMatcher) Matches(info os.FileInfo) (nextMatcher globerous.Matcher, walk bool, err error) {
	if info.Name() == m.text {
		return m.next, true, nil
	}
	return nil, false, nil
}

func newMyCustomMatcher(text string, next globerous.Matcher) globerous.Matcher {
	return &myCustomMatcher{text, next}
}

func main() {
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

	matcher := compiler.MustCompile("dir1", "*", "test.txt")

	path, err := filepath.Abs(filepath.Join("examples", "testdata"))
	must(err)
	fmt.Println("Matches:")
	must(globerous.Print(globerous.NewOSGlobFs(), matcher, path, os.Stdout))
}
