package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	workingDir string
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	var err error
	workingDir, err = os.Getwd()
	must(err)
}

func embedFileCode(path string) string {
	fullPath, err := filepath.Abs(path)
	must(err)
	if !strings.HasPrefix(fullPath, workingDir) {
		panic(fmt.Errorf("cannot include files outside of the working directory"))
	}
	contentB, err := ioutil.ReadFile(fullPath)
	if err != nil {
		panic(fmt.Errorf("failed to embed file '%s'", path))
	}

	return "```go\n" + string(contentB) + "\n```\n"
}

func renderMarkdown(contents string, out io.Writer) {
	funcMap := template.FuncMap{
		"embed": embedFileCode,
	}
	tpl, err := template.New("embedTpl").Funcs(funcMap).Parse(contents)
	must(err)
	must(tpl.Execute(out, nil))
}

func main() {
	if len(os.Args) != 3 {
		panic(fmt.Errorf("must specify the markdown file to embed"))
	}
	markdownFileIn := os.Args[1]
	markdownFileOut := os.Args[2]

	in, err := ioutil.ReadFile(markdownFileIn)
	must(err)

	out, err := os.OpenFile(markdownFileOut, os.O_RDWR|os.O_CREATE, 0775)
	defer out.Close()
	must(err)
	renderMarkdown(string(in), out)
	out.Sync()
}
