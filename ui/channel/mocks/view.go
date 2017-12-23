package mocks

import (
	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/marcusolsson/tui-go"
)

// View implements channel.View for testing.
type View struct {
	tui.Widget

	topic      string
	nick       string
	connection string
	presence   string
	mode       string
	events     []data.Event

	renderer   channel.EventRenderer
	controller channel.Controller
}

func (v *View) SetTopic(s string)                   { v.topic = s }
func (v *View) SetNick(s string)                    { v.nick = s }
func (v *View) SetConnection(s string)              { v.connection = s }
func (v *View) SetPresence(s string)                { v.presence = s }
func (v *View) SetMode(s string)                    { v.mode = s }
func (v *View) SetEvents(s []data.Event)            { v.events = s }
func (v *View) SetRenderer(s channel.EventRenderer) { v.renderer = s }
func (v *View) Attach(s channel.Controller)         { v.controller = s }
