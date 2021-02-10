package globerous

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Determines whether a certain file or folder should be explored. If nextMatcher is nil, then children of a folder
// will not be explored. If walk is set to true, the file or folder will be output to the WalkFn passed to the Walk
// function.
//
// Note a Matcher can also return itself if the same matcher should be used to match sub-levels. A good example of where
// this might be useful is the "**" matcher.
type Matcher interface {
	Matches(info os.FileInfo) (nextMatcher Matcher, walk bool, err error)
}

type WalkFn = func(dir string, info os.FileInfo, error error) error
type WalkFnSimple = func(dir string, info os.FileInfo) error

type GlobFs interface {
	ReadDir(dir string) ([]os.FileInfo, error)
}

// Print a list of files matched by the matcher out to a io.Writer, split by newline.
func Print(fs GlobFs, matcher Matcher, dir string, w io.Writer) error {
	return WalkSimple(fs, matcher, dir, func(dir string, info os.FileInfo) error {
		fullPath := filepath.Join(dir, info.Name())
		_, err := fmt.Fprintln(w, fullPath)
		return err
	})
}

// Get a list of all files and folders matched by the matcher.
func List(fs GlobFs, matcher Matcher, dir string) ([]string, error) {
	var out []string
	err := WalkSimple(fs, matcher, dir, func(dir string, info os.FileInfo) error {
		out = append(out, filepath.Join(dir, info.Name()))
		return nil
	})
	return out, err
}

func WalkSimple(fs GlobFs, matcher Matcher, dir string, walkFn WalkFnSimple) error {
	return Walk(fs, matcher, dir, func(dir string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return walkFn(dir, info)
	})
}

// Walk the fs for files and folders matched by the matcher.
func Walk(fs GlobFs, matcher Matcher, dir string, walkFn WalkFn) error {
	pathStack := &pathStack{}
	topDir := dir
	files, err := fs.ReadDir(topDir)
	if err != nil {
		if err := walkFn("", nil, err); err != nil {
			return err
		}
	}
	pathStack.Push(topDir, files, matcher)
	currentPath := pathStack.Pop()
	for currentPath != nil {
		if err := doMatch(fs, pathStack, currentPath, walkFn); err != nil {
			return err
		}
		currentPath = pathStack.Pop()
	}
	return nil
}

func doMatch(fs GlobFs, pathStack *pathStack, currentPath *pathDetails, walkFn WalkFn) error {
	matcher := currentPath.matcher
	for _, info := range currentPath.files {
		parentPath := currentPath.parentPath
		relativePathFromRoot := filepath.Join(parentPath, info.Name())
		if info.IsDir() {
			nextMatcher, walk, err := matcher.Matches(info)
			if err != nil {
				if err := walkFn("", nil, err); err != nil {
					return err
				}
			}
			if walk {
				if err := walkFn(parentPath, info, nil); err != nil {
					return err
				}
				continue
			}
			if nextMatcher == nil {
				continue
			}
			files, err := fs.ReadDir(relativePathFromRoot)
			if err != nil {
				if err := walkFn("", nil, err); err != nil {
					return err
				}
			}
			pathStack.Push(relativePathFromRoot, files, nextMatcher)
		} else {
			_, walk, err := matcher.Matches(info)
			if err != nil {
				if err := walkFn("", nil, err); err != nil {
					return err
				}
			}
			if walk {
				if err := walkFn(parentPath, info, nil); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
