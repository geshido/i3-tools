package main

import (
	"flag"
	"fmt"
	"github.com/geshido/i3-tools/pkg/i3tools"
	"go.i3wm.org/i3/v4"
	"log"
	"os"
)

const (
	OpDetect = "detect"
	OpToggle = "toggle"
)

func main() {
	var (
		op             string
		scratchPadName string
	)
	flag.StringVar(&op, "op", OpDetect, "operation on scratchpad: detect - to check if it exists, toggle - to toggle scratchpad")
	flag.StringVar(&scratchPadName, "s", "", "scratchpad name")
	flag.Parse()

	if scratchPadName == "" {
		flag.Usage()
		os.Exit(1)
	}
	if !itemInStrings([]string{OpDetect, OpToggle}, op) {
		flag.Usage()
		os.Exit(1)
	}

	tree, err := i3.GetTree()
	handleError(err)

	var foundNode *i3.Node
	nodes := i3tools.FindAll(tree.Root, func(node *i3.Node) bool {
		return node.Type == i3.Con && itemInStrings(node.Marks, scratchPadName)
	})
	if len(nodes) > 0 {
		foundNode = nodes[0]
	}
	if foundNode == nil {
		fmt.Println("")
		return
	}

	switch op {
	case OpDetect:
		fmt.Println(foundNode.Name)
	case OpToggle:
		_, _ = i3.RunCommand(fmt.Sprintf("[con_id=\"%d\"] scratchpad show", foundNode.ID))
	}

}

func itemInStrings(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
