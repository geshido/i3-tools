package container_move

import (
	"fmt"
	"go.i3wm.org/i3/v4"
)

type Usecase struct{}
type selector interface {
	Select() (string, error)
}

func (u Usecase) Run(workspaceSelector selector) error {
	wsName, err := workspaceSelector.Select()
	if err != nil {
		return fmt.Errorf("cannot get workspace name: %w", err)
	}
	if wsName == "" {
		return nil
	}

	if _, err := i3.RunCommand(fmt.Sprintf("move container to workspace %s", wsName)); err != nil {
		return fmt.Errorf("cannot execute move command: %w", err)
	}
	return nil
}
