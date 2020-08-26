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
		activeColor    string
		truncate       int
	)
	flag.StringVar(&op, "op", OpDetect, "operation on scratchpad: detect - to check if it exists, toggle - to toggle scratchpad")
	flag.StringVar(&scratchPadName, "s", "", "scratchpad name")
	flag.StringVar(&activeColor, "highlight", "fff", "RGB color to highlight visible scratchpad window name")
	flag.IntVar(&truncate, "truncate", 20, "truncate window title. 0 - to disable.")
	flag.Parse()

	if scratchPadName == "" {
		flag.Usage()
		os.Exit(1)
	}
	if !itemInStrings([]string{OpDetect, OpToggle}, op) {
		flag.Usage()
		os.Exit(1)
	}

	switch op {
	case OpToggle:
		tree, err := i3.GetTree()
		handleError(err)

		foundNode := findScratchNode(tree, scratchPadName)
		if foundNode == nil {
			return
		}
		_, _ = i3.RunCommand(fmt.Sprintf("[con_id=\"%d\"] scratchpad show", foundNode.ID))
	case OpDetect:
		reciever := i3.Subscribe(i3.WorkspaceEventType, i3.WindowEventType, i3.OutputEventType)
		ch := make(chan i3.Event)
		go func() {
			defer func() { _ = reciever.Close() }()
			defer close(ch)
			for {
				if !reciever.Next() {
					return
				}
				ch <- reciever.Event()
			}
		}()
		printScratchpadIfExists(scratchPadName, truncate, activeColor)
		for range ch {
			printScratchpadIfExists(scratchPadName, truncate, activeColor)
		}
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

func truncateString(s string, max int) string {
	if max == 0 {
		return s
	}
	runes := []rune(s)
	if len(runes) < max {
		return s
	}
	return string(runes[:max-1])
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func findScratchNode(tree i3.Tree, scratchPadName string) *i3.Node {
	var foundNode *i3.Node
	nodes := i3tools.FindAll(tree.Root, func(node *i3.Node) bool {
		return node.Type == i3.Con && itemInStrings(node.Marks, scratchPadName)
	})
	if len(nodes) > 0 {
		foundNode = nodes[0]
	}
	return foundNode
}

func printScratchpadIfExists(scratchPadName string, truncate int, activeColor string) {

	tree, err := i3.GetTree()
	handleError(err)

	foundNode := findScratchNode(tree, scratchPadName)
	if foundNode == nil {
		fmt.Println("")
		return
	}
	// search for foundNode workspace in currently open workspaces
	visible := false
	foundWorkspaces := i3tools.FindAll(tree.Root, func(node *i3.Node) bool {
		if node.Type != i3.WorkspaceNode {
			return false
		}

		if len(i3tools.FindAll(node, func(n *i3.Node) bool { return n.ID == foundNode.ID })) > 0 {
			return true
		}

		return false
	})
	var foundWorkspace *i3.Node
	if len(foundWorkspaces) > 0 {
		foundWorkspace = foundWorkspaces[0]
	}
	if foundWorkspace != nil {
		// get list of outputs and check if found workspace is visible now
		outputs, err := i3.GetOutputs()
		handleError(err)
		for _, output := range outputs {
			if output.CurrentWorkspace == foundWorkspace.Name {
				visible = true
				break
			}
		}
	}
	truncated := truncateString(foundNode.Name, truncate)
	if visible {
		fmt.Printf("%%{F#%s}%s%%{F-}\n", activeColor, truncated)
	} else {
		fmt.Println(truncated)
	}
}
