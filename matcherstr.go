package globerous

import (
	"os"
	"strings"
)

type equalMatcher struct {
	str      string
	comparer func(actual, expected string) bool
	next     Matcher
}

func (m equalMatcher) Matches(info os.FileInfo) (Matcher, bool, error) {
	if m.comparer(info.Name(), m.str) {
		return m.next, m.next == nil, nil
	}
	return nil, false, nil
}

// Matches a path part based on a string comparison function.
func StrFnMatcher(actual string, comparer func(actual, expected string) bool, next Matcher) Matcher {
	return &equalMatcher{actual, comparer, next}
}

// Matches a path part if it equals the value of str.
func EqualMatch(str string, next Matcher) Matcher {
	return StrFnMatcher(str, func(actual string, expected string) bool {
		return actual == expected
	}, next)
}

// Matches a path part if it has a prefix matching the prefix value.
func HasPrefixMatcher(prefix string, next Matcher) Matcher {
	return StrFnMatcher(prefix, strings.HasPrefix, next)
}

// Matches a path part if it has a suffix matching the suffix value.
func HasSuffixMatcher(suffix string, next Matcher) Matcher {
	return StrFnMatcher(suffix, strings.HasSuffix, next)
}
