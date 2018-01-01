package ui_test

import (
	"testing"

	"github.com/cceckman/discoirc/ui"
	_ "github.com/cceckman/discoirc/ui/channel"
	"github.com/cceckman/discoirc/ui/client"
	"github.com/cceckman/discoirc/ui/mocks"
)

func TestActivateChannel(t *testing.T) {
	/* TODO: Disabled until implemented
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	u := discomocks.NewUI()
	u.Update(func() {
		u.SetWidget(widgets.NewSplash())
	})

	ctl := ui.New(ctx, u)

	u.Add(1) // Expect one update to change root, keybindings, etc.
	ctl.ActivateChannel("foonet", "#discoirc")
	u.RunSync(func() {
		if _, ok := u.Root.(channel.View); !ok {
			t.Errorf("unexpected view at UI root: got: %+v want: controller.View", u.Root)
		}
	})
	*/
}

func TestActivateClient(t *testing.T) {
	u := mocks.NewUI()

	ctl := ui.New(u, mocks.NewBackend())

	ctl.ActivateClient()
	// Panic: negative WG counter
	u.RunSync(func() {
		if _, ok := u.Root.(client.View); !ok {
			t.Errorf("unexpected view at UI root: got: %+v want: client.View", u.Root)
		}
	})
}
