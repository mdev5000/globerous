package globerous

import (
	"os"
)

type pathStack struct {
	node *pathNode
	last *pathNode
}

type pathDetails struct {
	parentPath string
	files      []os.FileInfo
	matcher    Matcher
}

type pathNode struct {
	path *pathDetails
	next *pathNode
}

func (p *pathStack) Push(parentPath string, files []os.FileInfo, matcher Matcher) {
	path := &pathDetails{
		parentPath: parentPath,
		files:      files,
		matcher:    matcher,
	}
	node := &pathNode{path: path}
	if p.node == nil || p.last == nil {
		p.node = node
		p.last = node
	} else {
		p.last.next = node
		p.last = node
	}
}

func (p *pathStack) Pop() *pathDetails {
	if p.node == nil {
		return nil
	}
	path := p.node.path
	p.node = p.node.next
	return path
}

//func (p *pathStack) String() string {
//	node := p.node
//	lines := ""
//	for node != nil {
//		path := node.path
//		var files []string
//		for _, f := range path.files {
//			if f.IsDir() {
//				files = append(files, f.Name() + "/")
//			} else {
//				files = append(files, f.Name())
//			}
//		}
//		lines += fmt.Sprintf("%s: %s\n", path.matcher.String(), strings.Join(files, ","))
//		node = node.next
//	}
//	return lines
//}
