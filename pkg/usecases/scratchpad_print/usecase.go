package scratchpad_print

import (
	"context"
	"fmt"
	"github.com/geshido/i3-tools/pkg/i3tools"
	"go.i3wm.org/i3/v4"
	"strings"
)

type Input struct {
	ScratchPadName string
	Truncate       int
	ActiveColor    string
	RemoveSuffix   string
}
type Usecase struct{}

func (u Usecase) Run(input Input) error {
	events, cancel := i3tools.Subscribe(context.Background(),
		i3.WorkspaceEventType,
		i3.WindowEventType,
		i3.OutputEventType,
	)
	defer cancel()
	if err := printScratchpadIfExists(input); err != nil {
		return fmt.Errorf("cannot detect scratchpad: %w", err)
	}

	for range events {
		if err := printScratchpadIfExists(input); err != nil {
			return fmt.Errorf("cannot detect scratchpad: %w", err)
		}
	}
	return nil
}
func printScratchpadIfExists(input Input) error {

	tree, err := i3.GetTree()
	if err != nil {
		return fmt.Errorf("cannot get i3 tree: %w", err)
	}

	foundNode := i3tools.FindScratchNode(tree, input.ScratchPadName)
	if foundNode == nil {
		fmt.Println("")
		return nil
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
		if err != nil {
			return fmt.Errorf("cannot get outputs: %w", err)
		}
		for _, output := range outputs {
			if output.CurrentWorkspace == foundWorkspace.Name {
				visible = true
				break
			}
		}
	}
	out := foundNode.Name
	if strings.TrimSpace(input.RemoveSuffix) != "" {
		out = strings.TrimSuffix(out, input.RemoveSuffix)
	}

	var truncated bool
	out, truncated = truncateString(out, input.Truncate)
	if truncated {
		out += "…"
	}

	if visible {
		fmt.Printf("%%{F#%s}%s%%{F-}\n", input.ActiveColor, out)
	} else {
		fmt.Println(out)
	}
	return nil
}

func truncateString(s string, max int) (string, bool) {
	if max == 0 {
		return s, false
	}
	runes := []rune(s)
	if len(runes) < max {
		return s, false
	}
	return string(runes[:max-1]), true
}
