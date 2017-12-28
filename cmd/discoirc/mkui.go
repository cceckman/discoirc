package main

import (
	"context"
	"time"

	"github.com/cceckman/discoirc/backend/stub"
	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel/controller"
	"github.com/cceckman/discoirc/ui/channel/view"
	"github.com/marcusolsson/tui-go"
)

func GetStubChannel(ctx context.Context, ui controller.UIControl, network, nick, channel string) (*stub.Channel, tui.Widget) {
	m := stub.NewChannel(ctx)
	m.SetMeta(data.Channel{
		Name: channel,
		Connection: data.Connection{
			Network: network,
			Nick:    nick,
		},
	})

	v := view.New()
	_ = controller.New(ctx, ui, v, m)
	return m, v
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
