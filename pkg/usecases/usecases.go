package usecases

import (
	"github.com/geshido/i3-tools/pkg/storage"
	"github.com/geshido/i3-tools/pkg/usecases/container_alttab"
	"github.com/geshido/i3-tools/pkg/usecases/container_move"
	"github.com/geshido/i3-tools/pkg/usecases/promote_window"
	"github.com/geshido/i3-tools/pkg/usecases/scratchpad_print"
	"github.com/geshido/i3-tools/pkg/usecases/scratchpad_toggle"
	"github.com/geshido/i3-tools/pkg/usecases/workspace_rename"
	"github.com/geshido/i3-tools/pkg/usecases/workspace_switch"
)

type Registry struct {
	ScratchpadPrint  scratchpad_print.Usecase
	ScratchpadToggle scratchpad_toggle.Usecase
	WorkspaceSwitch  workspace_switch.Usecase
	WorkspaceRename  workspace_rename.Usecase
	ContainerMove    container_move.Usecase
	ContainerAltTab  container_alttab.Usecase
	PromoteWindow    promote_window.Usecase
}

func Build() Registry {
	alttabStorage, err := storage.NewFileStorage("/tmp/alttab")
	if err != nil {
		panic(err)
	}

	return Registry{
		ScratchpadPrint:  scratchpad_print.Usecase{},
		ScratchpadToggle: scratchpad_toggle.Usecase{},
		WorkspaceSwitch:  workspace_switch.Usecase{},
		WorkspaceRename:  workspace_rename.Usecase{},
		ContainerMove:    container_move.Usecase{},
		ContainerAltTab: container_alttab.Usecase{
			Storage: alttabStorage,
		},
		PromoteWindow: promote_window.Usecase{},
	}
}
