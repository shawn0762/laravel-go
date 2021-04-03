package routing

import (
	"strings"
	"testing"
)

func TestNode_Search(t *testing.T) {
	root := NewRoot()

	insert(root, "/a/b/c")
	insert(root, "/a/{x}")
	insert(root, "/a/{id?}/{type}")

	s := strings.Trim("/a/s", "/")
	p := strings.Split(s, "/")
	r := root.Search(p, 0)

	if r != nil {
		t.Error("something went wrong")
	}
}

func insert(node *node, pattern string) {
	parts := strings.Split(strings.Trim(pattern, "/"), "/")
	node.insert(nil, pattern, parts, 0)
}
