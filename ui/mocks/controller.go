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
// Operations on it should be run within its UI.
type Controller struct {
	*UI

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
