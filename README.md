# File pattern matching for go lang

A package for running flexible file matching patterns to filter and match files.

```go
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

```


## Supported globbing patterns

A limited set of globbing patterns are supported specifically:

- prefix wildcard: `*.txt`
- suffix wildcard: `file.*`
- recursive wildcard: `**`

For more control of patterns use the **regex** or **hybrid** match compiler. See the
next section for more details.


## Running regex patterns

You can run regex patterns with either the `globerous.RegexPartCompiler` or via the
`globerous.HybridGlobRegexPartCompiler`. The latter allows use of both regex and glob
pattern matching.

Here's an example of the hybrid compiler:

```go
package globerous_test

import (
	"fmt"
	"github.com/mdev5000/globerous"
	"os"
	"path/filepath"
)

// The hybrid compiler allows use of both simple globbing patterns and regex.
func ExampleHybridGlobRegexPartCompiler_hybridMatching() {
	path, err := filepath.Abs(filepath.Join("testdata", "examplesfs"))
	if err != nil {
		panic(err)
	}

	compiler := globerous.NewCompiler(globerous.HybridGlobRegexPartCompiler)

	// If a part starts with ^ or ends with $ then it is handled as regex instead of as a glob.
	matcher := compiler.MustCompile("^first|third$", "**", "*.txt")

	fmt.Println("Matches:")
	err = globerous.WalkSimple(globerous.NewOSGlobFs(), matcher, path, func(dir string, info os.FileInfo) error {
		fmt.Println(" ", filepath.Join(dir, info.Name()))
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Output:
	// Matches:
	//   /Users/matt/devtmp/go/globerous/testdata/examplesfs/first/nested/first.txt
	//   /Users/matt/devtmp/go/globerous/testdata/examplesfs/third/nested/deeper/third.txt
}

```


## Virtual filesystems

Globerous has support for the **blang vfs** library (https://github.com/blang/vfs). This
allows you to virtualize the filesystem and for example use an in-memory fs.

```go
package globerous_test

import (
	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
	"github.com/mdev5000/globerous"
	"os"
	"path/filepath"
)

// Globerous works out of the box with the blang vfs package (https://github.com/blang/vfs).
func Example_vfs() {
	fs := memfs.Create()
	root := filepath.Join("some", "path")
	if err := vfs.MkdirAll(fs, root, 0777); err != nil {
		panic(err)
	}
	if err := vfs.WriteFile(fs, filepath.Join(root, "first.txt"), nil, 0777); err != nil {
		panic(err)
	}
	if err := vfs.WriteFile(fs, filepath.Join(root, "second.txt"), nil, 0777); err != nil {
		panic(err)
	}
	compiler := globerous.NewCompiler(globerous.GlobPlusPartCompiler)
	if err := globerous.Print(fs, compiler.MustCompile("*", "path", "*.txt"), "/", os.Stdout); err != nil {
		panic(err)
	}
	// Output:
	///some/path/first.txt
	///some/path/second.txt
}

```


## Custom pattern matching

You can extend or replace the pattern matching used by using your own part compiler.

```go
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

```


You can also use the raw matchers instead of compiling.

```go
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

```


# Testing paths

You can test paths either using real filesystem information or simple path information.

```go
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

```

