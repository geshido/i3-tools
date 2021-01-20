package i3tools

import (
	"context"
	"go.i3wm.org/i3/v4"
)

func Subscribe(ctx context.Context, types ...i3.EventType) (chan i3.Event, func()) {
	c, cancel := context.WithCancel(ctx)
	ch := make(chan i3.Event)

	go func() {
		reciever := i3.Subscribe(types...)
		defer func() { _ = reciever.Close() }()
		defer close(ch)

		for {
			if !reciever.Next() {
				return
			}

			select {
			case ch <- reciever.Event():
			case <-c.Done():
				return
			}
		}
	}()

	return ch, cancel
}
