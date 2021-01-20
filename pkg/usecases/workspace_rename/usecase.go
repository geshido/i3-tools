package workspace_rename

import (
	"fmt"
	"go.i3wm.org/i3/v4"
	"strings"
)

type Usecase struct{}

type selector interface {
	Select() (string, error)
}

func (u Usecase) Run(nameSelector selector) error {
	wsName, err := nameSelector.Select()
	if err != nil {
		return fmt.Errorf("cannot get new name: %w", err)
	}
	if strings.TrimSpace(wsName) == "" {
		return nil
	}

	if _, err := i3.RunCommand(fmt.Sprintf("rename workspace to %s", wsName)); err != nil {
		return fmt.Errorf("cannot execute rename command: %w", err)
	}
	return nil
}
