package mocks

import (
	"github.com/cceckman/discoirc/ui/channel"
)

// UIController is a mock channel.UIController.
type UIController struct {
	Received []string
	Size     int
}

func (c *UIController) Input(s string) {
	c.Received = append(c.Received, s)
}

func (c *UIController) Resize(n int) {
	c.Size = n
	// TODO: implement optional side effects?
}

var _ channel.UIController = &UIController{}
