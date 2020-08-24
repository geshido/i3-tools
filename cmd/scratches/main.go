package main

import (
	"flag"
	"fmt"
	"go.i3wm.org/i3/v4"
	"log"
	"os"
	"sync"
)

func itemInStrings(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func main() {
	var (
		op             string
		scratchPadName string
	)
	flag.StringVar(&op, "op", "detect", "operation on scratchpad: detect - to check if it exists, toggle - to toggle scratchpad")
	flag.StringVar(&scratchPadName, "s", "", "scratchpad name")
	flag.Parse()

	if scratchPadName == "" {
		flag.Usage()
		os.Exit(1)
	}
	if !itemInStrings([]string{"detect", "toggle"}, op) {
		flag.Usage()
		os.Exit(1)
	}

	tree, err := i3.GetTree()
	handleError(err)

	var foundNode *i3.Node
	nodes := findAll(tree.Root, func(node *i3.Node) bool {
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
	case "detect":
		fmt.Println(foundNode.Name)
	case "toggle":
		_, _ = i3.RunCommand(fmt.Sprintf("[con_id=\"%d\"] scratchpad show", foundNode.ID))
	}

}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func findAll(start *i3.Node, predicate func(*i3.Node) bool) []*i3.Node {
	var result []*i3.Node

	if predicate(start) {
		result = append(result, start)
	}

	ch := make(chan *i3.Node)
	wg := sync.WaitGroup{}
	wg.Add(2)
	finder := func(nodes []*i3.Node) {
		defer wg.Done()
		for _, n := range nodes {
			for _, nn := range findAll(n, predicate) {
				ch <- nn
			}
		}
	}
	go finder(start.Nodes)
	go finder(start.FloatingNodes)
	go func() { wg.Wait(); close(ch) }()

	for n := range ch {
		result = append(result, n)
	}

	return result
}
