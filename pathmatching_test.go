package globerous

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_MatchesPath_MatchesCorrectlyWithFakeStrings(t *testing.T) {
	path := FakePath("some/path/file.txt", false)
	matcher := tGlobMatcher(t, "some", "*", "file.txt")
	match, partialMatch, err := MatchesPath(matcher, path)
	require.Nil(t, err)
	require.True(t, match)
	require.False(t, partialMatch)
}

func Test_MatchesPath_CanIndicatePartialMatch(t *testing.T) {
	path := FakePath("some/path/file.txt", false)
	matcher := tGlobMatcher(t, "some", "*")
	match, partialMatch, err := MatchesPath(matcher, path)
	require.Nil(t, err)
	require.False(t, match)
	require.True(t, partialMatch)
}

// @todo test GetPathFromString
