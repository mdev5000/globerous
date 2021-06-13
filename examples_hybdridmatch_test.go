package globerous_test

import (
	"fmt"
	"github.com/mdev5000/globerous"
	"os"
	"path/filepath"
	"strings"
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

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println("Matches:")
	err = globerous.WalkSimple(globerous.NewOSGlobFs(), matcher, path, func(dir string, info os.FileInfo) error {
		fmt.Println(" ", filepath.Join(strings.TrimPrefix(dir, wd), info.Name()))
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Output:
	// Matches:
	//   /testdata/examplesfs/first/nested/first.txt
	//   /testdata/examplesfs/third/nested/deeper/third.txt
}
