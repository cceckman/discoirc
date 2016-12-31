//
// editgui.go
// provides GUI objects, etc. for this prototype

package main

import (
	"context"
	"fmt"
	"github.com/jroimartin/gocui"
)

// SetupUI sets up managers for a gocui.Gui, but does not start the main loop.
func SetupUI(g *gocui.Gui) error {
	// Create a context that closing the UI terminates.
	ctx, cancel := context.WithCancel(context.Background())
	_ = ctx

	// Start window layout-er.
	// Note: manager must be provided before setting keybindings (e.g. below.)
	mv := &ModelView{ui: g}
	go mv.Start()
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

// ModelView is a view manager.
type ModelView struct{
	ui *gocui.Gui
}

// Type enforcement.
var _ gocui.Manager = &ModelView{}

// Start begins operations that update the ModelView. It should be run in a background thread (i.e. go m.Start())
func (m *ModelView) Start() { }

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
