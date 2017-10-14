package view

import (
	"github.com/jroimartin/gocui"
	"github.com/cceckman/discoirc/discocli/model"
)

// MessageEditor provides a line editor for the Chat input.
type MessageEditor struct {
	Channel model.Channel

	View *gocui.View
}

// NewMessageEditor attaches a MessageEditor to the provided input View.
// It must be called in the UI thread, i.e. by Layout, Update, or a key binding.
func NewMessageEditor(c model.Channel, v *gocui.View) {
	r := &MessageEditor{
		View: v,
	}
	v.Editor = r
	v.Editable = true
}

// Edit provide editing functionality to the input field.
func (m *MessageEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
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
		// TODO send message to controller
	case gocui.KeyTab:
		// TODO tab-complete names
	}
	// TODO: Handle history with up/down
}
