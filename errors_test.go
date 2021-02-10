package globerous

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestGlobChildMatcherError_UnwrapsNilWhenNoError(t *testing.T) {
	me := NewGlobChildMatcherError()
	require.Nil(t, errors.Unwrap(me))
}

// @todo add test that unwrap unwraps

func TestGlobChildMatcherError_CanUnwrapsFirstError(t *testing.T) {
	me := NewGlobChildMatcherError()
	me.Add(fmt.Errorf("something bad happened"))
	me.Add(fmt.Errorf("another bad thing happened"))
	err := errors.Unwrap(me)
	require.EqualError(t, err, "something bad happened")
}

func TestGlobChildMatcherError_Error_PrintsForLineOfErrors(t *testing.T) {
	me := NewGlobChildMatcherError()
	me.Add(fmt.Errorf("something bad happened"))
	me.Add(fmt.Errorf("another bad thing happened"))
	me.Add(fmt.Errorf("multiple line\nerror"))
	require.EqualError(t, me, strings.TrimSpace(`
Multiple errors occurred:
- something bad happened
- another bad thing happened
- multiple line...
`)+"\n")
}

func TestGlobChildMatcherError_Error_WhenOnlyOneErrorPrintsThatError(t *testing.T) {
	me := NewGlobChildMatcherError()
	me.Add(fmt.Errorf("something bad happened"))
	require.EqualError(t, me, "something bad happened")
}

func TestGlobChildMatcherError_Error_CanBeUsedAsOriginalError(t *testing.T) {
	me := NewGlobChildMatcherError()
	me.Add(fmt.Errorf("something bad happened"))
	var err error = me
	var me2 *GlobChildMatcherError
	require.True(t, errors.As(err, &me2))
}
