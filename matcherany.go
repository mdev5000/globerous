package globerous

import (
	"os"
)

type recursiveAnyMatcher struct {
	next     Matcher
	children []Matcher
}

// Matcher for the ** glob pattern.
//
// Rule is matches the next matcher in the current or any child directories.
//
// Examples:
// 		For the glob pattern /**/first/text.txt
//
//		Yes 	/first/text.txt
//		Yes 	/first/first/text.txt
//		Yes 	/ducks/first/text.txt
//		Yes 	first/ducks/first/text.txt
//
func (r recursiveAnyMatcher) Matches(info os.FileInfo) (Matcher, bool, error) {
	var children []Matcher
	errors := NewGlobChildMatcherError()

	nextMatcher, walk, err := r.next.Matches(info)
	if err != nil {
		errors.Add(err)
	} else if nextMatcher != nil {
		children = append(children, nextMatcher)
	}

	for _, matcher := range r.children {
		childNextMatcher, walkChild, err := matcher.Matches(info)
		if err != nil {
			errors.Add(err)
			continue
		}
		walk = walk || walkChild
		if childNextMatcher != nil {
			children = append(children, childNextMatcher)
		}
	}

	return recursiveAnyMatcher{
		next:     r.next,
		children: children,
	}, walk, errors.Result()
}

// Matcher handling the glob pattern "**". Will match at all directory levels and delegates to child matches, unless it
// is the last matcher.
//
// Ex. for the pattern "**/first/test.txt", all of the following will match:
//
// - first/text.txt
//
// - another/first/test.txt
//
// - first/another/first/test.txt
func AnyRecursiveMatcher(child Matcher) Matcher {
	// @todo add case for when child is nil.
	return recursiveAnyMatcher{
		next: child,
	}
}

type anyMatcher struct {
	next Matcher
}

func (r *anyMatcher) Matches(os.FileInfo) (Matcher, bool, error) {
	return r.next, r.next == nil, nil
}

// Matcher handling the glob pattern "*". Will match any file or folder, but a file or folder must exist at that level.
//
// For example, for "*/test.txt", the following will match:
//
// - first/test.txt
//
// - second/test.txt
//
// But this will not:
//
// - test.txt
func AnyMatcher(child Matcher) Matcher {
	return &anyMatcher{
		next: child,
	}
}
