package crawler

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestTreeBuilder(t *testing.T) {
	content, err := ioutil.ReadFile("./tree_test.txt")
	if err != nil {
		t.Errorf(err.Error())
	}

	urls := strings.Split(string(content), "\n")
	BuildTree(urls)
}
