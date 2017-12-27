package controller_test

import (
	"context"
	"testing"

	"github.com/cceckman/discoirc/ui/client/controller"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/cceckman/discoirc/ui/client/mocks"
	discomocks "github.com/cceckman/discoirc/ui/mocks"
)

func TestActivation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ui := discomocks.NewUI()
	ui.Update(func() {
		ui.SetWidget(&mocks.View{})
	})

	ctl := controller.New(ctx, ui)

	ui.Add(1) // Expect one update to change root, keybindings, etc.
	ctl.ActivateChannel("foonet", "#discoirc")
	ui.RunSync(func() {
		root := ui.Root
		if _, ok := root.(channel.View); !ok {
			t.Errorf("unexpected view at UI root: got: %+v want: controller.View", ui.Root)
		}
	})
}
