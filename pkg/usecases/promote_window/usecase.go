package promote_window

import (
	"fmt"
	"github.com/geshido/i3-tools/pkg/i3tools"
	"go.i3wm.org/i3/v4"
)

type Usecase struct{}

func (Usecase) Run() error {
	wsList, err := i3.GetWorkspaces()
	if err != nil {
		return fmt.Errorf("cannot get workspaces: %w", err)
	}

	tree, err := i3.GetTree()
	if err != nil {
		return fmt.Errorf("cannot get tree: %w", err)
	}

	for _, ws := range wsList {
		if !ws.Focused {
			continue
		}

		wsNode := i3tools.FindByID(tree.Root, i3.NodeID(ws.ID))
		if wsNode == nil {
			continue
		}

		master := findBiggestWindow(wsNode)
		if master == nil {
			continue
		}

		_, err := i3.RunCommand(fmt.Sprintf("swap container with con_id %d", master.ID))
		if err != nil {
			return fmt.Errorf("cannot run command: %w", err)
		}
	}
	return nil
}

func findBiggestWindow(container *i3.Node) *i3.Node {
	var maxLeaf *i3.Node
	maxArea := int64(0)
	for _, leaf := range i3tools.Leaves(container) {
		area := leaf.Rect.Width * leaf.Rect.Height
		if !leaf.Focused && area > maxArea {
			maxArea = area
			maxLeaf = leaf
		}
	}

	return maxLeaf
}
