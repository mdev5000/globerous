package globerous

import (
	"regexp"
	"strings"
)

// A hybrid part matcher that uses both regex and glob syntax. The a part is matched as a regex if it starts with a
// ^ or ends with  $. Otherwise glob syntax is used. Also supports the "**" operator.
//
// Examples: "**", "*", "test.txt", "*.txt", "^first|second$"
//
func HybridGlobRegexPartCompiler(partStr string, next NextPartFn) (Matcher, error) {
	childMatcher, err := next()
	if err != nil {
		return nil, err
	}
	switch partStr {
	case "**":
		return AnyRecursiveMatcher(childMatcher), nil
	case "*":
		return AnyMatcher(childMatcher), nil
	}
	if strings.HasPrefix(partStr, "^") || strings.HasSuffix(partStr, "$") {
		return parseRegexToMatcher(partStr, childMatcher)
	}
	return GlobMatcherSm(partStr, childMatcher), nil
}

// Matches strings based on glob syntax, but additionally all support '**', for nested matching.
func GlobPlusPartCompiler(partStr string, next NextPartFn) (Matcher, error) {
	childMatcher, err := next()
	if err != nil {
		return nil, err
	}
	switch partStr {
	case "**":
		return AnyRecursiveMatcher(childMatcher), nil
	case "*":
		return AnyMatcher(childMatcher), nil
	}
	return GlobMatcherSm(partStr, childMatcher), nil
}

func parseRegexToMatcher(partStr string, childMatcher Matcher) (Matcher, error)  {
	partRegex, err := regexp.Compile(partStr)
	return RegexMatcher(partRegex, childMatcher), err
}
