package globerous

type NextPartFn = func() (Matcher, error)

// A pattern for parsing a set of matching parts when creating a MatchCompiler. See the custom example in the
// NewCompiler docs for an example of how this works.
type PartCompiler = func(part string, next NextPartFn) (Matcher, error)

// Compiles a list of "matching parts" into a Matcher that can be used to match and filter file paths. How the matching
// parts are interpreted differs based on how the compiler is initialized. See NewCompiler for details on options
// available.
type MatchCompiler struct {
	partCompiler PartCompiler
}

// Creates a matching MatchCompiler for generating a Matcher, based on a PartCompiler. Providing a different PartCompiler
// effects how the matcher strings is interpreted and converted to a Matcher. A custom PartCompiler can be provided to
// customize the behaviour further.
func NewCompiler(partCompiler PartCompiler) *MatchCompiler {
	return &MatchCompiler{
		partCompiler: partCompiler,
	}
}

// Same as Compile, but will panic if an error occurs during compilation.
func (s MatchCompiler) MustCompile(parts ...string) Matcher {
	m, err := s.Compile(parts...)
	if err != nil {
		panic(err)
	}
	return m
}

// Compile a set of matching strings into a Matcher. How is determined by the compiler type, see the NewCompiler docs
//for more details.
//
// Note the "/" is not considered a separator. However to accomplish this you can use strings.Split before passing it
//as shown in the example MatchingForwardSlash.
func (s MatchCompiler) Compile(parts ...string) (Matcher, error) {
	return s.compile(parts)
}

func (s MatchCompiler) compile(parts []string) (Matcher, error) {
	if len(parts) == 0 {
		return nil, nil
	}
	part := parts[0]
	childParts := parts[1:]
	matcher, err := s.partCompiler(part, func() (Matcher, error) {
		return s.compile(childParts)
	})
	return matcher, err
}
