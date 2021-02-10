package globerous

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Tests if a match matches a specific path. Useful for pattern matching on path in utilities like watchers.
func MatchesPath(m Matcher, path []os.FileInfo) (matches, partialMatch bool, err error) {
	currentMatcher := m
	lastIndex := len(path) - 1
	for i, part := range path {
		nextMatcher, walk, matchErr := currentMatcher.Matches(part)
		if matchErr != nil {
			return false, false, matchErr
		}
		if walk {
			if i == lastIndex {
				matches = true
				return
			} else {
				partialMatch = true
			}
		}
		if nextMatcher == nil {
			return
		}
		currentMatcher = nextMatcher
	}
	return
}

// Convert a path into an array usable with the MatchesPath function.
func GetPathFromString(fs GlobFs, root string, path string) ([]os.FileInfo, error) {
	return GetPathFromArray(fs, root, strings.Split(path, "/"))
}

// Convert a path array into an array usable with the MatchesPath function.
func GetPathFromArray(fs GlobFs, root string, path []string) ([]os.FileInfo, error) {
	parent := root
	var out []os.FileInfo
	for _, part := range path {
		f, err := findFileWithName(fs, parent, part)
		if err != nil {
			return nil, nil
		}
		if f == nil {
			return nil, fmt.Errorf("path does not exist in filesystem")
		}
		parent = filepath.Join(parent, f.Name())
		out = append(out, f)
	}
	return out, nil
}

func findFileWithName(fs GlobFs, parent string, part string) (os.FileInfo, error) {
	files, err := fs.ReadDir(parent)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.Name() == part {
			return f, nil
		}
	}
	return nil, nil
}

type fakeFileInfo struct {
	name  string
	isDir bool
}

func (f fakeFileInfo) Name() string {
	return f.name
}

func (f fakeFileInfo) IsDir() bool {
	return f.isDir
}

func (f fakeFileInfo) Size() int64 {
	return 0
}

func (f fakeFileInfo) Mode() os.FileMode {
	return 0
}

func (f fakeFileInfo) ModTime() time.Time {
	return time.Now()
}

func (f fakeFileInfo) Sys() interface{} {
	return nil
}

// Create a fake path to match with the MatchesPath function.
//
// Note fake defaults are created for the os.FileInfo meaning that matchers that require methods other than Name()
// and IsDir() may not function correctly.
//
// See MatchesPath docs for an example of using this function.
//
func FakePath(path string, isDir bool) []os.FileInfo {
	var out []os.FileInfo
	parts := strings.Split(path, "/")
	lastIndex := len(parts) - 1
	for _, part := range parts[:lastIndex] {
		out = append(out, &fakeFileInfo{part, true})
	}
	out = append(out, &fakeFileInfo{parts[lastIndex], isDir})
	return out
}
