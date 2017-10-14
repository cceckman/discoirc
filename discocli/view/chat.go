// Package chat provides the Chat view/model/viewmodel for the IRC channel view.
package view

import (
	"log"

	"github.com/jroimartin/gocui"
)

type Manager struct {
	*ChatViewInfo

	Log *log.Logger
	Done chan<- ViewInfo
}

// ChatViewInfo is the normal view of a channel or PM thread: scrolling text, an input field, etc.
type ChatViewInfo struct {
	Connection, Channel string
}

func (vi *ChatViewInfo) NewManager(log *log.Logger, done chan<- ViewInfo) gocui.Manager {
	return &Manager{
		ChatViewInfo: vi,
		Log: log,
		Done: done,
	}
}

var _ gocui.Manager = &Manager{}

func (m *Manager) Layout(g *gocui.Gui) error {
	return nil
}


