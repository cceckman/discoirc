package view

import (
	"github.com/jroimartin/gocui"
)

func (vm *Channel) NewInput() gocui.Manager {
	return ChannelInput(vm)
}

// ChannelInput is the ViewModel for the Channel window's message input field.
type ChannelInput *Channel

var _ gocui.Editor = *ChannelInput(nil)
var _ gocui.Manager = *ChannelInput(nil)

// Layout sets up the ChannelInput view.
func (c *ChannelInput) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	// No border at the bottom of the terminal, full width.
	ax, ay, bx, by := -1, maxY-2, maxX, maxY
	v, err := g.SetView(ChannelInputView, ax, ay, bx, by)
	switch err {
	case nil:
		return nil
	case gocui.ErrUnknownView:
		c.Log.Printf("%s [start] initial setup", ChannelInputView)
		defer c.Log.Printf("%s [done] initial setup", ChannelInputView)
		v.Frame = false
		v.Editable = true
		v.Editor = c
	default:
		return err
	}
	return nil

	return nil
}

// Send sends a message to the backend if a connection is available. If the message is sent, it clears the View buffer.
func (c *ChannelInput) Send(rawmsg string) {
	// TODO sanitize!
	msg := rawmsg

	// Check if we're connected. If not, just return - without even clearing the field.
	select {
	case <-c.connected:
		// pass
	default:
		return
	}
	c.channel.SendMessage(msg)
	// Callback; clear if the message is sent.
	c.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(ChannelInputView)
		if err == gocui.ErrUnknownView {
			return nil
		} else if err != nil {
			return err
		}

		// If the user has modified the buffer since we sent the message, keep the existing message.
		if v.Buffer() == rawmsg {
			v.Clear()
		}
		return nil
	})
}

// Edit provide editing functionality to the input field.
func (m *ChannelInput) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if ch != 0 && mod == 0 {
		v.EditWrite(ch)
		return
	}
	switch key {
	case gocui.KeySpace:
		v.EditWrite(' ')
	case gocui.KeyBackspace, gocui.KeyBackspace2:
		// TODO check position; don't delete prompt
		v.EditDelete(true)
	case gocui.KeyDelete:
		v.EditDelete(false)
	case gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	case gocui.KeyArrowLeft:
		// TODO check position; don't go past prompt
		v.MoveCursor(01, 0, false)
	case gocui.KeyEnter:
		// This handler is run in the event loop, but may require talking with a network thread
		// to send the message; so launch a goroutine to handle it.
		// Will take a callback to clear if needed.
		go m.Send(v.Buffer())
	case gocui.KeyTab:
		// TODO tab-complete names
	}
	// TODO: Handle history with up/down
}
