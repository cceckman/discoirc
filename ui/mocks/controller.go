package mocks

// TODO:
// - Unify on a single Controller, rather than having Controller, UI, UpdateCounter
// - Add Context to constructor to clean up old threads when test is done
func NewController() *Controller {
	return &Controller{
		UI: NewUI(),
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
// Its methods should only be called within the .Update thread.
type Controller struct {
	*UI

	V       ActiveView
	Network string
	Channel string
}

func (c *Controller) ActivateClient() {
	c.V = ClientView
	c.Network, c.Channel = "", ""
}

func (c *Controller) ActivateChannel(network, channel string) {
	c.V = ChannelView
	c.Network, c.Channel = network, channel
}
