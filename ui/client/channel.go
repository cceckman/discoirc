package client

import (
	"fmt"

	"github.com/cceckman/discoirc/data"

	"github.com/marcusolsson/tui-go"
)

// NewChannel returns a new Channel view.
func NewChannel(network *Network, name string) *Channel {
	r := &Channel{
		network:         network,
		name:            name,
		indicatorWidget: newIndicator(),
		nameWidget:      tui.NewLabel(name),
		modeWidget:      tui.NewLabel(""),
		// TODO: replace these with a "localized-compressed" widget,
		// which shrinks e.g. "messages" / "msgs" / "✉" as space is needed,
		// in an appropriately localized fashion.
		unreadWidget:  tui.NewLabel("✉ ?"),
		membersWidget: tui.NewLabel("? ☺"),
	}

	r.Widget = tui.NewHBox(
		r.indicatorWidget,
		tui.NewVBox(
			tui.NewHBox(r.nameWidget, tui.NewSpacer(), r.modeWidget),
			tui.NewHBox(r.unreadWidget, tui.NewSpacer(), r.membersWidget),
		),
	)

	return r
}

// Channel is a view (Widget) giving an overview of a channel.
type Channel struct {
	tui.Widget
	network *Network
	name    string

	focus           bool
	indicatorWidget *indicator
	nameWidget      *tui.Label
	modeWidget      *tui.Label
	unreadWidget    *tui.Label
	membersWidget   *tui.Label
}

// UpdateChannel updates the view with the provided channel state.
func (c *Channel) UpdateChannel(ch data.ChannelState) {
	c.modeWidget.SetText(ch.ChannelMode)
	c.unreadWidget.SetText(fmt.Sprintf("✉ %d", ch.Unread))
	c.membersWidget.SetText(fmt.Sprintf("%d ☺", ch.Members))
}

// SetFocused indicates the user's focus is on the Channel.
// The Channel should provide a visual indicator of this focus and respond to
// key events.
func (c *Channel) SetFocused(focus bool) {
	c.focus = focus
	if focus {
		c.indicatorWidget.SetFill('|')
	} else {
		c.indicatorWidget.SetFill(' ')
	}
}

// IsFocused returns true if the user's focus is on the channel.
func (c *Channel) IsFocused() bool {
	return c.focus
}

// Name gives the channel's name.
func (c *Channel) Name() string {
	return c.name
}

// OnKeyEvent handles key presses, if the channel is selected.
func (c *Channel) OnKeyEvent(ev tui.KeyEvent) {
	if !c.focus {
		return
	}

	ctl := c.network.client.controller
	if ev.Key == tui.KeyEnter && ctl != nil {
		ctl.ActivateChannel(c.network.name, c.name)
	}
}
