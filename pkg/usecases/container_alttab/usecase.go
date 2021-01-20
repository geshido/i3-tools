package container_alttab

import (
	"context"
	"fmt"
	"github.com/geshido/i3-tools/pkg/i3tools"
	"go.i3wm.org/i3/v4"
	"log"
)

type Usecase struct {
	Storage Storage
}

type Storage interface {
	Store(interface{}) error
	Read(interface{}) error
	Cleanup() error
}

func (u *Usecase) Subscribe() error {
	if u.Storage == nil {
		panic("storage is nil")
	}
	events, cancel := i3tools.Subscribe(context.Background(),
		i3.WindowEventType,
		i3.OutputEventType,
		i3.WorkspaceEventType,
	)
	ids := make(chan i3.NodeID)
	defer cancel()
	defer func() {
		_ = u.Storage.Cleanup()
	}()
	defer close(ids)

	go func() {
		var prevID, curID i3.NodeID
		for id := range ids {
			if curID != id {
				prevID = curID
				curID = id
				fmt.Printf("prev=%d cur=%d\n", prevID, curID)
			}
			if err := u.Storage.Store(prevID); err != nil {
				log.Printf("[ERR] cannot store previous container id: %v", err)
			}
		}
	}()

	for range events {
		tree, err := i3.GetTree()
		if err != nil {
			return fmt.Errorf("cannot get tree: %w", err)
		}
		focused := i3tools.FindAll(tree.Root, func(n *i3.Node) bool {
			return n.Focused
		})
		if len(focused) > 0 {
			ids <- focused[0].ID
		}
	}
	return nil
}

func (u *Usecase) SwitchToPrevious() error {
	var id int64
	if err := u.Storage.Read(&id); err != nil {
		return fmt.Errorf("cannot read previous container id: %w", err)
	}

	if id == 0 {
		fmt.Println("prev con id is 0")
		return nil
	}

	if _, err := i3.RunCommand(fmt.Sprintf("[con_id=\"%d\"] focus", id)); err != nil {
		return fmt.Errorf("cannot focus container with id=%d: %w", id, err)
	}
	return nil
}
