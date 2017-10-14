// Package view provides UI handlers for discocli.
package view

import (
	"log"

	"github.com/jroimartin/gocui"
)

// How to handle this setup?
// Start by attaching something to the UI.
// That "something" includes decision criteria for what the display should be, and needs to
// instantiate and attach a Manager appropriately.
// The decision of what to instantiate will ultimately be provided by the backend.
// So: This needs to
// - Attach to backend OR decide itself what the current mode should be
// - Receive signals that say "change the layout"
//   -> Reattach managers when those come in
//   -> Exit when the signal says "none"

// What do layouts need to initialize?
// - Logger
// - "Next" channel for changing views
// - View-specific intitialization data (channel, etc.)
// - Later: Model source (daemon connection)

// ViewInfo provides the information necessary to initialize one of the supported Views.
type ViewInfo interface {
	// NewManager creates a new gocui.Manager from this ViewInfo.
	// The ViewInfo incorporates whatever view-specific parameters there are; the arguments to this
	// function provide the default logger and "change view" channels.
	NewManager(*log.Logger, chan<- ViewInfo) gocui.Manager
}

// LayoutSwitcher initializes
// the layout.
type LayoutSwitcher struct {
	next chan ViewInfo

	Gui *gocui.Gui
	Log *log.Logger
}

// await handles requests to update the layout.
func (l *LayoutSwitcher) await() {
	for n := range l.next {
		if n == nil {
			l.Log.Print("LayoutSwitcher received request for a nil view, ignoring.")
			continue
		}
		m := n.NewManager(l.Log, l.next)
		l.Gui.SetManager(m)
		// Default keybinding of ctrl-c to exit
		l.Gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit)
	}
}

func Quit(*gocui.Gui, *gocui.View) error {
	return gocui.ErrQuit
}

// StartLayout starts the LayoutSwitcher using the provided initial view and logger.
func StartLayoutSwitcher(g *gocui.Gui, log *log.Logger, initView ViewInfo) *LayoutSwitcher {
	layout := &LayoutSwitcher{
		next: make(chan ViewInfo),
		Gui: g,
		Log: log,
	}
	go layout.await()
	go func() {
		layout.next <- initView
	}()
	return layout
}

func (l *LayoutSwitcher) Done() {
	close(l.next)
}
