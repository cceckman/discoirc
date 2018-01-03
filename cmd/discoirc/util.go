package main

import (
	"context"
	"time"

	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/backend/demo"
)

func GetTheme() *tui.Theme {
	t := tui.NewTheme()
	t.SetStyle("reversed", tui.Style{
		Reverse: tui.DecorationOn,
	})
	return t
}

// Toggle is a stub.Channel wrapper that toggles message / metadata updates.
type Toggle struct {
	Net, Chan string
	Duration  time.Duration
	*demo.Demo

	cancelNetwork  func()
	cancelChannel  func()
	cancelMessages func()
}

func (t *Toggle) Network() {
	if t.cancelNetwork != nil {
		t.cancelNetwork()
		t.cancelNetwork = nil
		return
	}

	var ctx context.Context
	ctx, t.cancelNetwork = context.WithCancel(context.Background())
	go func() {
		tick := time.NewTicker(t.Duration)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				t.Demo.TickNetwork(t.Net)
			}
		}
	}()
}

func (t *Toggle) Channel() {
	if t.cancelChannel != nil {
		t.cancelChannel()
		t.cancelChannel = nil
		return
	}

	var ctx context.Context
	ctx, t.cancelChannel = context.WithCancel(context.Background())
	go func() {
		tick := time.NewTicker(t.Duration)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				t.Demo.TickChannel(t.Net, t.Chan)
			}
		}
	}()
}

func (t *Toggle) Messages() {
	if t.cancelMessages != nil {
		t.cancelMessages()
		t.cancelMessages = nil
		return
	}

	var ctx context.Context
	ctx, t.cancelMessages = context.WithCancel(context.Background())
	go func() {
		tick := time.NewTicker(t.Duration)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				t.Demo.TickMessages(t.Net, t.Chan)
			}
		}
	}()
}
