package globerous

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func newFakeFile(path string) os.FileInfo {
	return &fakeFileInfo{path, false}
}

func newFakeFolder(path string) os.FileInfo {
	return &fakeFileInfo{path, true}
}

func Test_GlobMatcher_CorrectlyMatchesNames(t *testing.T) {
	m := GlobMatcherSm("first.txt", nil)
	next, walk, _ := m.Matches(newFakeFile("first.txt"))
	require.Nil(t, next)
	require.True(t, walk)
	next, walk, _ = m.Matches(newFakeFile("firsttxt"))
	require.Nil(t, next)
	require.False(t, walk)
}

func Test_GlobMatcher_CorrectlyMatchesNames2(t *testing.T) {
	m := GlobMatcherSm("first", nil)
	next, walk, _ := m.Matches(newFakeFile("first"))
	require.Nil(t, next)
	require.True(t, walk)
	next, walk, _ = m.Matches(newFakeFile("second"))
	require.Nil(t, next)
	require.False(t, walk)
}

func Test_GlobMatcher_CorrectlyMatchesWildCards(t *testing.T) {
	m := GlobMatcherSm("*.txt", nil)
	next, walk, _ := m.Matches(newFakeFile("first.txt"))
	require.Nil(t, next)
	require.True(t, walk)

	next, walk, _ = m.Matches(newFakeFile("another.txt"))
	require.Nil(t, next)
	require.True(t, walk)

	next, walk, _ = m.Matches(newFakeFile(".txt"))
	require.Nil(t, next)
	require.True(t, walk)
}
