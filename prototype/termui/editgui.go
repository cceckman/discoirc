//
// editgui.go
// provides GUI objects, etc. for this prototype

package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/jroimartin/gocui"
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
		notice: make(chan string, 1),
		input: make(chan string, 1),
	}
}

// ModelView is a view manager.
type ModelView struct {
	ui *gocui.Gui

	notice chan string
	input chan string
}

// Type enforcement.
var _ gocui.Manager = &ModelView{}

// Start begins operations that run outside the main thread. It should be run in a background thread (i.e. go m.Start())
func (m *ModelView) Start(ctx context.Context) {
	go m.WatchInput(ctx)
}

// WatchInput watches the input channel
// Its input is non-blocking
func (m *ModelView) WatchInput(ctx context.Context) {
	buffer := make([]string, 2)
	

}

// Layout implements gocui.Manager for ModelView.
func (m *ModelView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil

}
