package globerous

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type alwaysMatchesMatcher struct{}

func (a alwaysMatchesMatcher) Matches(os.FileInfo) (Matcher, bool, error) {
	return a, true, nil
}

func newAlwaysMatchesMatcher() Matcher {
	return &alwaysMatchesMatcher{}
}

type neverMatchesMatcher struct{}

func (n neverMatchesMatcher) Matches(os.FileInfo) (Matcher, bool, error) {
	return nil, false, nil
}

func newNeverMatchesMatcher() Matcher {
	return &neverMatchesMatcher{}
}

func tGetMaxDepthCount(t *testing.T, m Matcher, maxDepth int) int {
	for i := 0; i < maxDepth; i++ {
		m, _, _ = m.Matches(newFakeFile("first"))
		if m == nil {
			return i + 1
		}
	}
	require.Nil(t, fmt.Errorf("hit maximum depth of %d", maxDepth))
	return maxDepth
}

func Test_MaxDepthMatcher_LimitsToMaxDepth(t *testing.T) {
	matcher := MaxDepthMatcher(4, newAlwaysMatchesMatcher())
	require.Equal(t, tGetMaxDepthCount(t, matcher, 5), 4)
}

func Test_MaxDepthMatcher_ExitsAsSoonAsTheChildMatcherExits(t *testing.T) {
	matcher := MaxDepthMatcher(4, newNeverMatchesMatcher())
	require.Equal(t, tGetMaxDepthCount(t, matcher, 5), 1)
}

func Test_FileOrFolderMatcher_CanMatchOnlyFilesOrFolders(t *testing.T) {
	fileMatcher := FileOrFolderMatcher(false, newAlwaysMatchesMatcher())
	_, walk, _ := fileMatcher.Matches(newFakeFile("file1"))
	require.True(t, walk)
	_, walk, _ = fileMatcher.Matches(newFakeFolder("folder1"))
	require.False(t, walk)

	folderMatcher := FileOrFolderMatcher(true, newAlwaysMatchesMatcher())
	_, walk, _ = folderMatcher.Matches(newFakeFolder("folder2"))
	require.True(t, walk)
	_, walk, _ = folderMatcher.Matches(newFakeFile("file2"))
	require.False(t, walk)
}
