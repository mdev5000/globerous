package globerous

import (
	"os"
)

type maxDepthMatcher struct {
	maxDepth     int
	currentDepth int
	matcher      Matcher
}

func (m maxDepthMatcher) Matches(info os.FileInfo) (Matcher, bool, error) {
	if m.maxDepth <= m.currentDepth {
		return nil, false, nil
	}
	nextMatcher, walk, err := m.matcher.Matches(info)
	if err != nil {
		return nil, false, err
	}
	if nextMatcher == nil {
		return nil, walk, nil
	}
	return maxDepthMatcher{
		maxDepth:     m.maxDepth,
		currentDepth: m.currentDepth + 1,
		matcher:      nextMatcher,
	}, walk, nil
}

// Limits the matching to a maximum depth. Useful when used with the "**" glob pattern to limit the depth.
func MaxDepthMatcher(maxDepth int, matcher Matcher) Matcher {
	return maxDepthMatcher{
		maxDepth:     maxDepth,
		currentDepth: 1,
		matcher:      matcher,
	}
}

type fileOrFolderMatcher struct {
	matchDir bool
	matcher  Matcher
}

func (f fileOrFolderMatcher) Matches(info os.FileInfo) (Matcher, bool, error) {
	nextMatcher, walk, err := f.matcher.Matches(info)
	if walk {
		return nextMatcher, f.matchDir == info.IsDir(), err
	}
	return nextMatcher, false, err
}

// Will match only files or only folders depending on the value of matchDir. Can be used to limit the output to only
// files or folders.
func FileOrFolderMatcher(matchDir bool, matcher Matcher) Matcher {
	return &fileOrFolderMatcher{matchDir: matchDir, matcher: matcher}
}
