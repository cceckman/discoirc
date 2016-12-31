//
// editgui.go
// provides GUI objects, etc. for this prototype

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cceckman/discoirc/prototype/bufchan"
	"github.com/jroimartin/gocui"
)

const (
	messagesView = "messages"
	inputView    = "input"
	noticesView  = "notices"
)

// SetupUI sets up managers for a gocui.Gui, but does not start the main loop.
func SetupUI(g *gocui.Gui) error {
	// Create a context that closing the UI terminates.
	ctx, cancel := context.WithCancel(context.Background())
	_ = ctx

	// Start window layout-er.
	// Note: manager must be provided before setting keybindings (e.g. below.)
	mv := New(g)
	go mv.Start(ctx)
	g.SetManager(mv)

	// Pass through ctrl+c to quit.
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(_ *gocui.Gui, _ *gocui.View) error {
			cancel()
			return gocui.ErrQuit
		},
	); err != nil {
		return err
	}

	return nil
}

// New returns a new ModelView
func New(g *gocui.Gui) *ModelView {
	return &ModelView{
		ui: g,
	}
}

// ModelView is a view manager.
// It should be Start-ed before being attached as a Manager.
type ModelView struct {
	ui *gocui.Gui

	notices  *bufchan.Bufchan
	input    *bufchan.Bufchan
	messages *bufchan.Bufchan
}

// Type enforcement.
var _ gocui.Manager = &ModelView{}

// Start begins operations that run outside the main thread. It should be run in a background thread (i.e. go m.Start())
func (m *ModelView) Start(ctx context.Context) {
	m.notices = bufchan.New(ctx)
	m.input = bufchan.New(ctx)
	m.messages = bufchan.New(ctx)

	// Just testing...
	go func() {
		tick := time.NewTicker(time.Second)
		defer tick.Stop()
		for i := 0; true; i++ {
			select {
			case <-ctx.Done():
				return
			case m.input.In() <- fmt.Sprintf(" %d\n", i):
				// Delay until the tick.
				<-tick.C
			}
		}
	}()

	go m.WatchInput(ctx)
	go m.WriteMessages(ctx)
	//	go m.WriteNotices(ctx)
}

// WatchInput watches the input channel, and demuxes into 'messages' and 'notices'.
func (m *ModelView) WatchInput(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case input := <-m.input.Out():
			if len(input) == 0 {
				continue
			}
			if input[0] == '!' {
				m.notices.In() <- input[1:]
			} else {
				m.messages.In() <- input
			}
		}
	}
}

// WriteMessages listens on the relevant channel, and writes messages to the UI.
func (m *ModelView) WriteMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-m.messages.Out():
			if len(message) == 0 {
				continue
			}
			m.ui.Execute(func(g *gocui.Gui) error {
				if v, err := g.View(messagesView); err == nil {
					fmt.Fprint(v, message)
				} else {
					return err
				}
				return nil
			})
		}
	}
	return
}

// WriteNotices listens on the relevant channel, and writes pop-up notifications to the UI.
func (m *ModelView) WriteNotices(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case notice := <-m.notices.Out():
			if len(notice) == 0 {
				continue
			}
			m.ui.Execute(func(g *gocui.Gui) error {
				maxX, maxY := g.Size()
				l := len(notice) / 2
				if v, err := g.SetView(
					noticesView,
					maxX/2-l-1, maxY/2,
					maxX/2+l+1, maxY/2+2,
				); err != nil {
					if err != gocui.ErrUnknownView {
						return err
					}
					v.Clear()
					g.SetViewOnTop(noticesView)
					fmt.Fprintln(v, notice)
				}
				return nil
			})
		}
	}
}

// Layout implements gocui.Manager for ModelView.
func (m *ModelView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(messagesView, 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true
		v.Title = "Messages"
	}
	return nil
}
