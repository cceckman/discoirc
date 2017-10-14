// Package chat provides the Chat view/model/viewmodel for the IRC channel view.
package view

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"
)

// ChatManager is the ViewModel for the Chat view.
type ChatManager struct {
	*ChatViewInfo

	Log  *log.Logger
	done chan<- ViewInfo
}

// ChatViewInfo is the normal view of a channel or PM thread: scrolling text, an input field, etc.
type ChatViewInfo struct {
	Connection, Channel string
}

func (vi *ChatViewInfo) NewManager(log *log.Logger, done chan<- ViewInfo) gocui.Manager {
	return &ChatManager{
		ChatViewInfo: vi,
		Log:          log,
		done:         done,
	}
}

var _ gocui.Manager = &ChatManager{}

// Layout sets up the Chat view. It creates new views as necessary, including starting threads.
func (m *ChatManager) Layout(g *gocui.Gui) error {
	m.Log.Print("Chat: [start] layout")
	defer m.Log.Print("Chat: [done] layout")
	// Create three views: input, status, messages.
	// Create them in that order, since 'input' is fixed-len, but 'status' and 'messages' may need to flex.
	if err := m.layoutInput(g); err != nil {
		return err
	}
	if err := m.layoutStatus(g); err != nil {
		return err
	}
	if err := m.layoutMessages(g); err != nil {
		return err
	}
	g.SetCurrentView("input")
	return nil
}

func (m *ChatManager) layoutInput(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	// These are good!
	ax, ay, bx, by := -1, maxY-2, maxX, maxY
	v, err := g.SetView("input", ax, ay, bx, by)
	switch err {
	case nil:
		return nil
	case gocui.ErrUnknownView:
		m.Log.Print("Chat/input: [start] initial layout")
		defer m.Log.Print("Chat/input: [done] initial layout")
		v.Editable = true
		v.Frame = false
		v.FgColor = gocui.ColorBlue
		// TODO handle editor behavior
		fmt.Fprint(v, "your text here")
	default:
		return err
	}
	return nil
}

func (m *ChatManager) layoutStatus(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	ax, ay, bx, by := -1, maxY-3, maxX, maxY-1
	if v, err := g.SetView("status", ax, ay, bx, by); err != nil {
		if err != gocui.ErrUnknownView {
			// unknown error, percolate it up
			return err
		}
		// Initialize the view.
		m.Log.Print("Chat/status: [start] initial layout")
		defer m.Log.Print("Chat/status: [done] initial layout")

		v.Editable = false
		v.Frame = false
		// TODO more color customization
		v.BgColor, v.FgColor = gocui.ColorWhite, gocui.ColorBlack
		// TODO attach controller/model here, instead of static init.
		if _, err := fmt.Fprintf(v, "%s / %s", m.Connection, m.Channel); err != nil {
			return err
		}
	}
	return nil
}

func (m *ChatManager) layoutMessages(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	ax, ay, bx, by := 0, 0, maxX-1, maxY-3
	if v, err := g.SetView("messages", ax, ay, bx, by); err != nil {
		if err != gocui.ErrUnknownView {
			// unknown error, percolate it up
			return err
		}
		// Initialize the view.
		m.Log.Print("Chat/messages: [start] initial layout")
		defer m.Log.Print("Chat/messages: [done] initial layout")
		v.Editable = false
		v.Frame = true

		// TODO attach controller/model here, instead of fake init.
		v.Autoscroll = true
		go func() {
			m.Log.Print("Chat/messages: [start] counting bottles")
			defer m.Log.Print("Chat/messages: [done] counting bottles")
			max := 99
			ctx, cancel := context.WithCancel(context.Background())
			for i := max; i >= 0; i-- {
				time.Sleep(time.Millisecond * 500)
				select {
				case <-ctx.Done():
					max = 0
				default:
					// do nothing
				}

				msg := fmt.Sprintf("\n%d bottles of beer on the wall, %d bottles of beer...", i, i)

				m.Log.Print("Chat/messages: [start] handler for: ", msg)
				g.Update(func(g *gocui.Gui) error {
					v, err := g.View("messages")
					if err == gocui.ErrUnknownView {
						cancel()
						return nil
					} else if err != nil {
						return err
					}
					m.Log.Print(msg)
					if _, err := fmt.Fprint(v, msg); err != nil {
						return err
					}
					return nil
				})
			}
		}()
	}
	return nil
}
