package mocks

import (
	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/marcusolsson/tui-go"
)

// View implements channel.View for testing.
type View struct {
	tui.Widget

	Topic      string
	Nick       string
	Connection string
	Presence   string
	Mode       string
	Events     []data.Event

	Renderer   channel.EventRenderer
	Controller channel.UIController
}

func (v *View) SetTopic(s string)                   { v.Topic = s }
func (v *View) SetNick(s string)                    { v.Nick = s }
func (v *View) SetConnection(s string)              { v.Connection = s }
func (v *View) SetPresence(s string)                { v.Presence = s }
func (v *View) SetMode(s string)                    { v.Mode = s }
func (v *View) SetEvents(s []data.Event)            { v.Events = s }
func (v *View) SetRenderer(s channel.EventRenderer) { v.Renderer = s }
func (v *View) Attach(s channel.UIController)         { v.Controller = s }
