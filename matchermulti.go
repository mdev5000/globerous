package globerous

import (
	"os"
)

type multiMatcher struct {
	matchers []Matcher
}

func (m multiMatcher) Matches(info os.FileInfo) (Matcher, bool, error) {
	walk := false
	var nextMatchers []Matcher
	for _, matcher := range m.matchers {
		nextChild, childWalk, err := matcher.Matches(info)
		if err != nil {
			return matcher, false, err
		}
		walk = walk || childWalk
		if nextChild != nil {
			nextMatchers = append(nextMatchers, nextChild)
		}
	}
	return multiMatcher{matchers: nextMatchers}, walk, nil
}

// Matches a file or folder if any of the child matchers are a match.
func MultiMatcher(matchers ...Matcher) Matcher {
	return multiMatcher{matchers: matchers}
}
