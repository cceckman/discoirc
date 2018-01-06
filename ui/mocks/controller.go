package mocks

// NewController returns a mock Controller.
func NewController() *Controller {
	return &Controller{
		UI: NewUI(),
	}
}

// ActiveView is an enum of top-level views.
type ActiveView int

// These top-level views can be active.
const (
	UnknownView = ActiveView(iota)
	ClientView
	ChannelView
)

// Controller is mock global-level controller.
// Its methods should only be called within the .Update thread.
type Controller struct {
	*UI

	V       ActiveView
	Network string
	Channel string
}

// ActivateClient sets the Controller to the client view.
func (c *Controller) ActivateClient() {
	c.V = ClientView
	c.Network, c.Channel = "", ""
}

// ActivateChannel sets the Controller to the channel view.
func (c *Controller) ActivateChannel(network, channel string) {
	c.V = ChannelView
	c.Network, c.Channel = network, channel
}
