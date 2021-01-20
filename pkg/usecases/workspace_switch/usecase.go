package workspace_switch

import (
	"fmt"
	"go.i3wm.org/i3/v4"
)

type Usecase struct {
}

type selector interface {
	Select() (string, error)
}

func (u Usecase) Run(selector selector) error {
	wsName, err := selector.Select()
	if err != nil {
		return fmt.Errorf("cannot get workspace name: %w", err)
	}
	if wsName == "" {
		return nil
	}

	if _, err := i3.RunCommand(fmt.Sprintf("workspace %s", wsName)); err != nil {
		return fmt.Errorf("cannot run i3 command: %w", err)
	}

	return nil
}
