package crawler

import (
	"strings"
)

type Node struct {
	Name     string `json:"name"`
	Children []Node `json:"children"`
}

// AddToTree - Add new nodes to site url tree.
func AddToTree(root []Node, names []string) []Node {
	if len(names) > 0 {
		var i int
		name := names[0]
		for i = 0; i < len(root); i++ {
			// already in tree
			if root[i].Name == name {
				break
			}
		}

		if i == len(root) {
			// Patch empty root string
			if name == "" {
				name = "/"
			}
			root = append(root, Node{Name: name})
		}

		root[i].Children = AddToTree(root[i].Children, names[1:])
	}

	return root
}

// BuildTree - Create a JSON tree from a list of urls.
func BuildTree(urls []string) []Node {
	var tree []Node
	for i := range urls {
		tree = AddToTree(tree, strings.Split(urls[i], "/"))
	}

	return tree
}
