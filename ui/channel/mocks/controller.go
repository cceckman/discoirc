package mocks

import (
	"github.com/cceckman/discoirc/ui/channel"
)

// Controller is a mock channel.Controller.
type Controller struct {
	Received []string
	Size     int

	Position int
}

func (c *Controller) Input(s string) {
	c.Received = append(c.Received, s)
}

func (c *Controller) Resize(n int) {
	c.Size = n
	// TODO: implement optional side effects?
}

func (c *Controller) Scroll(up bool) {
	if up {
		c.Position += 1
	} else {
		c.Position -= 1
	}
}

var _ channel.Controller = &Controller{}
