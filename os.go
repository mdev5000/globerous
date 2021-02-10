package globerous

import (
	"io/ioutil"
	"os"
)

type osGlob struct{}

func (o osGlob) ReadDir(dir string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dir)
}

func NewOSGlobFs() GlobFs {
	return &osGlob{}
}
