package globerous

import (
	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func tGlobMatcher(t *testing.T, parts ...string) Matcher {
	compile := NewCompiler(GlobPlusPartCompiler)
	matcher, err := compile.Compile(parts...)
	require.Nil(t, err)
	return matcher
}

func tHybridMatcher(t *testing.T, parts ...string) Matcher {
	compile := NewCompiler(HybridGlobRegexPartCompiler)
	matcher, err := compile.Compile(parts...)
	require.Nil(t, err)
	return matcher
}

func tRequireNotMatches(t *testing.T, m Matcher, pathIsFile bool, path string) {
	tRequireMatchesEquals(t, m, pathIsFile, path, false)
}

func tRequireMatches(t *testing.T, m Matcher, pathIsFile bool, path string) {
	tRequireMatchesEquals(t, m, pathIsFile, path, true)
}

func tRequireMatchesEquals(t *testing.T, m Matcher, pathIsFile bool, path string, expectedMatches bool) {
	fs := memfs.Create()
	parts := strings.Split(path, "/")
	if !pathIsFile {
		require.Nil(t, vfs.MkdirAll(fs, path, 0777))
		return
	}
	folderPath := strings.Join(parts[0:len(parts)-1], "/")
	require.Nil(t, vfs.MkdirAll(fs, folderPath, 0777))
	require.Nil(t, vfs.WriteFile(fs, path, nil, 0777))
	//PrintFs(fs, os.Stdout)
	files, err := List(fs, m, "/")
	require.Nil(t, err)
	if expectedMatches {
		require.Equal(t, files, []string{path})
	} else {
		require.NotEqual(t, files, []string{path})
	}

}

func Test_GlobMatching_WildcardsWork(t *testing.T) {
	m := tGlobMatcher(t, "*", "first", "*")
	tRequireMatches(t, m, true, "/first/first/test.txt")
}

func Test_GlobMatching_MatchAllWildcardsWork(t *testing.T) {
	m := tGlobMatcher(t, "**", "first", "test.txt")

	tRequireMatches(t, m, true, "/first/first/test.txt")
	tRequireMatches(t, m, true, "/first/test.txt")
	tRequireMatches(t, m, true, "/first/second/first/test.txt")

	tRequireNotMatches(t, m, true, "/first/second/first/another.txt")
}

func Test_GlobMatching_MultiMatcher(t *testing.T) {
	compiler := NewCompiler(GlobPlusPartCompiler)
	m := MultiMatcher(
		compiler.MustCompile("first", "test.txt"),
		compiler.MustCompile("second", "nested", "test.txt"),
	)
	tRequireMatches(t, m, true, "/first/test.txt")
	tRequireMatches(t, m, true, "/second/nested/test.txt")

	tRequireNotMatches(t, m, true, "/first/first/test.txt")
	tRequireNotMatches(t, m, true, "/test.txt")
}

func Test_HybridMatching_SupportsRegexSyntax(t *testing.T) {
	m := tHybridMatcher(t, "^first|second$", "nested", "test.txt")

	tRequireMatches(t, m, true, "/first/nested/test.txt")
	tRequireMatches(t, m, true, "/second/nested/test.txt")

	tRequireNotMatches(t, m, true, "/third/nested/test.txt")
	tRequireNotMatches(t, m, true, "/nested/test.txt")
}

func Test_HybridMatching_SupportsMixGlobRegexSyntax(t *testing.T) {
	m := tHybridMatcher(t, "^first|second$", "**", "third", "*", "file.txt")

	tRequireMatches(t, m, true, "/first/third/nested/file.txt")
	tRequireMatches(t, m, true, "/second/third/nested/file.txt")

	tRequireNotMatches(t, m, true, "/third/nested/test.txt")
}
