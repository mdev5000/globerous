package globerous

import (
	"os"
	"regexp"
)

type regexMatcher struct {
	r    *regexp.Regexp
	next Matcher
}

func (r regexMatcher) Matches(info os.FileInfo) (nextMatcher Matcher, walk bool, err error) {
	if !r.r.MatchString(info.Name()) {
		return nil, false, nil
	}
	return r.next, r.next == nil, nil
}

// Matches a file or folder based on a regex expression.
func RegexMatcher(regexp *regexp.Regexp, next Matcher) Matcher {
	return &regexMatcher{
		r:    regexp,
		next: next,
	}
}
