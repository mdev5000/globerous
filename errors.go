package globerous

import (
	"errors"
	"strings"
)

// Certain matchers contain multiple child matchers and, in most cases, even if an error occurs in one, the other
// should still be run. This way the user can deal with specific errors that occur while matching and expected the
// matching to still occur consistently. Because of this, multiple child matchers can fail and need to be tracked
// accordingly.
type GlobChildMatcherError struct {
	Errors []error
}

func NewGlobChildMatcherError() *GlobChildMatcherError {
	return &GlobChildMatcherError{}
}

// If an error occurred returns itself, otherwise returns nil.
func (m *GlobChildMatcherError) Result() error {
	if len(m.Errors) == 0 {
		return nil
	}
	return m
}

func (m *GlobChildMatcherError) Add(err error) {
	m.Errors = append(m.Errors, err)
}

func (m *GlobChildMatcherError) Unwrap() error {
	if len(m.Errors) == 0 {
		return nil
	}
	err := m.Errors[0]
	errW := errors.Unwrap(err)
	if errW != nil {
		return errW
	}
	return err
}

func (m *GlobChildMatcherError) Error() string {
	if len(m.Errors) == 0 {
		return ""
	}
	if len(m.Errors) == 1 {
		return m.Errors[0].Error()
	}
	b := strings.Builder{}
	b.WriteString("Multiple errors occurred:\n")
	for _, err := range m.Errors {
		b.WriteString("- ")
		errMsg := err.Error()
		parts := strings.Split(errMsg, "\n")
		b.WriteString(parts[0])
		if len(parts) > 1 {
			b.WriteString("...")
		}
		b.WriteString("\n")
	}
	return b.String()
}

var _ error = &GlobChildMatcherError{}
