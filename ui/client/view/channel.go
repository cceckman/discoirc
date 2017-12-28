package view

import (
	"fmt"

	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/ui/client"
)

var _ client.ChannelView = &Channel{}

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

func (c *Channel) SetFocused(focus bool) {
	c.focus = focus
	if focus {
		c.indicatorWidget.SetFill('|')
	} else {
		c.indicatorWidget.SetFill(' ')
	}
}

func (c *Channel) IsFocused() bool {
	return c.focus
}

func (c *Channel) SetMode(m string) {
	c.modeWidget.SetText(m)
}

func (c *Channel) SetUnread(n int) {
	c.unreadWidget.SetText(fmt.Sprintf("✉ %d", n))
}

func (c *Channel) SetMembers(n int) {
	c.membersWidget.SetText(fmt.Sprintf("%d ☺", n))
}

func (c *Channel) Name() string {
	return c.name
}

func (c *Channel) OnKeyEvent(ev tui.KeyEvent) {
	if !c.focus {
		return
	}

	ctl := c.network.client.controller
	if ev.Key == tui.KeyEnter && ctl != nil {
		ctl.ActivateChannel(c.network.name, c.name)
	}
}
