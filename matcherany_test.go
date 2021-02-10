package globerous

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type alwaysFailsMatcher struct {
	err error
}

func (a alwaysFailsMatcher) Matches(info os.FileInfo) (Matcher, bool, error) {
	return nil, false, a.err
}
func newAlwaysFailsMatcher(err error) Matcher {
	return &alwaysFailsMatcher{err}
}

func TestRecursiveAnyMatcher_WhenMultipleMatchersFailWhensAllTheErrorsThatOccurs(t *testing.T) {
	am := recursiveAnyMatcher{}
	am.next = newAlwaysFailsMatcher(fmt.Errorf("next failure message"))
	am.children = []Matcher{
		newAlwaysFailsMatcher(fmt.Errorf("first child failure message")),
		newAlwaysFailsMatcher(fmt.Errorf("second child failure message")),
	}
	_, _, err := am.Matches(newFakeFile("aFile"))

	var mErr *GlobChildMatcherError

	require.True(t, errors.As(err, &mErr))
	require.EqualError(t, mErr.Errors[0], "next failure message")
	require.EqualError(t, mErr.Errors[1], "first child failure message")
	require.EqualError(t, mErr.Errors[2], "second child failure message")
}
