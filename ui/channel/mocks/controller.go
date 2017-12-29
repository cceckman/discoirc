package mocks

import (
	"fmt"

	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
)

// UIController is a mock channel.Controller.
type UIController struct {
	Received []string
	Size     int
	HasQuit  bool
}

func (c *UIController) Input(s string) {
	c.Received = append(c.Received, s)
}

func (c *UIController) Resize(n int) {
	c.Size = n
	// TODO: implement optional side effects?
}

func (c *UIController) Quit() {
	c.HasQuit = true
}

func (c *UIController) UpdateMeta( _ data.Channel) {
	panic(fmt.Errorf("unsupported call in mock instance"))
}

func (c *UIController) UpdateContents( _ data.Event) {
	panic(fmt.Errorf("unsupported in mock instance"))
}

var _ channel.Controller = &UIController{}
