package mocks


// ControllerUI is the functionality that a Controller
type ControllerUI interface {
	Update(func())
}

func NewController(ui ControllerUI) *Controller {
	return &Controller{
		UI: ui,
	}
}

// ActiveView indicates which view is active.
type ActiveView int

const (
	UnknownView = ActiveView(iota)
	ClientView
	ChannelView
)


// Controller is mock global-level controller.
// Operations on it should be run within its UI.
type Controller struct {
	UI ControllerUI

	V ActiveView
	Network string
	Channel string
}

func (c *Controller) ActivateClient() {
	c.UI.Update(func() {
		c.V = ClientView
		c.Network, c.Channel = "", ""
	})
}

func (c *Controller) ActivateChannel(network, channel string) {
	c.UI.Update(func() {
		c.V = ChannelView
		c.Network, c.Channel = network, channel
	})
}
