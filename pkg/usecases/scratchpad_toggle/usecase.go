package scratchpad_toggle

import (
	"fmt"
	"github.com/geshido/i3-tools/pkg/i3tools"
	"go.i3wm.org/i3/v4"
)

type Usecase struct{}

func (u Usecase) Run(scratchPadName string) error {
	tree, err := i3.GetTree()
	if err != nil {
		return fmt.Errorf("cannot get i3 tree: %w", err)
	}

	foundNode := i3tools.FindScratchNode(tree, scratchPadName)
	if foundNode == nil {
		return fmt.Errorf("cannot find scratchpad '%s'", scratchPadName)
	}
	if _, err := i3.RunCommand(fmt.Sprintf("[con_id=\"%d\"] scratchpad show", foundNode.ID)); err != nil {
		return fmt.Errorf("cannot toggle scratchpad: %w", err)
	}
	return nil
}
