package main

import (
	"context"
	"time"

	"github.com/cceckman/discoirc/backend/stub"
	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/marcusolsson/tui-go"
)

func GetStubChannel(ctx context.Context, ui channel.UIControl, network, chanName, nick string) *stub.Channel {
	m := stub.NewChannel(ctx)
	m.SetMeta(data.Channel{
		Name: chanName,
		Connection: data.Connection{
			Network: network,
			Nick:    nick,
		},
	})

	_ = channel.New(ctx, ui, channel.NewView(), m)
	return m
}

func GetTheme() *tui.Theme {
	t := tui.NewTheme()
	t.SetStyle("reversed", tui.Style{
		Reverse: tui.DecorationOn,
	})
	return t
}

// Toggle is a stub.Channel wrapper that toggles message / metadata updates.
type Toggle struct {
	Duration time.Duration
	*stub.Channel

	cancelMessages func()
	cancelMeta     func()
}

func (t *Toggle) Messages(ctx context.Context) {
	if t.cancelMessages != nil {
		t.cancelMessages()
		t.cancelMessages = nil
		return
	}

	ctx, t.cancelMessages = context.WithCancel(ctx)
	go t.GenerateMessages(ctx, t.Duration)
}

func (t *Toggle) Meta(ctx context.Context) {
	if t.cancelMeta != nil {
		t.cancelMeta()
		t.cancelMeta = nil
		return
	}

	ctx, t.cancelMeta = context.WithCancel(ctx)
	go t.GenerateMeta(ctx, t.Duration)
}
