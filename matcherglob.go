package globerous

import (
	"strings"
)

// Matches a glob pattern (https://en.wikipedia.org/wiki/Glob_(programming)). Note supports only a limited subset.
// Specifically only prefixes and suffixes of * are supported.
//
// Example:
//
// - *.txt
//
// - first*
func GlobMatcherSm(partStr string, next Matcher) Matcher {
	if strings.HasPrefix(partStr, "*") {
		return HasSuffixMatcher(partStr[1:], next)
	}
	if strings.HasSuffix(partStr, "*") {
		return HasPrefixMatcher(partStr[:len(partStr)-1], next)
	}
	return EqualMatch(partStr, next)
}
