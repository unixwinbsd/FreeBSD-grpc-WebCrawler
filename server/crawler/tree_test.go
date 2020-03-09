package crawler

import (
	"io/ioutil"
	"strings"
	"testing"
)

// TestTreeBuilder - Test the BuildTree
func TestTreeBuilder(t *testing.T) {
	content, err := ioutil.ReadFile("./tree_test.txt")
	if err != nil {
		t.Errorf(err.Error())
	}

	urls := strings.Split(string(content), "\n")
	res := BuildTree(urls)
	if res == "" {
		t.Error("tree should not be an empty string")
	}
}
