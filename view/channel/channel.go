// Package channel implements the Views and ViewModels for the Channel display.
package channel

import (
	"errors"

	"github.com/cceckman/discoirc/model"
	"github.com/cceckman/discoirc/view"
	"github.com/jroimartin/gocui"
)

const (
	// View names.
	ChannelInputView    = "channel input"
	ChannelStatusView   = "channel status"
	ChannelContentsView = "channel contents"
)

// Channel is the ViewModel for the Channel view.
type Channel struct {
	*view.Context

	Connection, Channel string

	// channel is not necessarily populated until 'connected' is closed.p
	// connected blocks some operations until the channel is properly connected from the client side.
	connected chan struct{}
	channel   model.Channel
}

func (vm *Channel) validate() error {
	if vm.Gui == nil {
		return errors.New("no Gui provided")
	}
	if vm.Log == nil {
		return errors.New("no Logger provided")
	}
	if vm.Connection == "" {
		return errors.New("no Connection provided")
	}
	if vm.Channel == "" {
		return errors.New("no Channel provided")
	}
	return nil
}

func (vm *Channel) Start() error {
	if err := vm.validate(); err != nil {
		return err
	}
	// Start client connection.
	vm.connected = make(chan struct{})
	go func() {
		defer close(vm.connected)
		vm.channel = vm.Backend.Connection(vm.Connection).Channel(vm.Channel)
	}()

	// Attach ViewModels.
	vm.Gui.SetManager(
		vm.NewInput(),
		vm.NewStatus(),
		vm.NewContents(),
		QuitManager,
	)

	return nil
}

// QuitManager is a Manager that provides a Ctrl+C quit handler.
var QuitManager gocui.Manager = gocui.ManagerFunc(func(g *gocui.Gui) error {
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		return gocui.ErrQuit
	})
	return nil
})
