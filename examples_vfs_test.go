package globerous_test

import (
	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
	"github.com/mdev5000/globerous"
	"os"
	"path/filepath"
)

// Globerous works out of the box with the blang vfs package (https://github.com/blang/vfs).
func Example_vfs() {
	fs := memfs.Create()
	root := filepath.Join("some", "path")
	if err := vfs.MkdirAll(fs, root, 0777); err != nil {
		panic(err)
	}
	if err := vfs.WriteFile(fs, filepath.Join(root, "first.txt"), nil, 0777); err != nil {
		panic(err)
	}
	if err := vfs.WriteFile(fs, filepath.Join(root, "second.txt"), nil, 0777); err != nil {
		panic(err)
	}
	compiler := globerous.NewCompiler(globerous.GlobPlusPartCompiler)
	if err := globerous.Print(fs, compiler.MustCompile("*", "path", "*.txt"), "/", os.Stdout); err != nil {
		panic(err)
	}
	// Output:
	///some/path/first.txt
	///some/path/second.txt
}
