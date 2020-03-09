package crawler

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Node struct {
	Name     string `json:"name"`
	Children []Node `json:"children"`
}

// AddToTree - Add new nodes to tree.
func AddToTree(root []Node, names []string) []Node {
	if len(names) > 0 {
		var i int
		for i = 0; i < len(root); i++ {

			// already in tree
			if root[i].Name == names[0] {
				break
			}
		}

		if i == len(root) {
			name := names[0]
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
func BuildTree(urls []string) string {
	var tree []Node
	for i := range urls {
		tree = AddToTree(tree, strings.Split(urls[i], "/"))
	}

	b, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		panic(err)
	}

	treeStr := string(b)
	fmt.Print(treeStr)
	return treeStr
}
