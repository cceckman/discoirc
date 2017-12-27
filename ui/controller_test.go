package ui_test

import (
	"context"
	"testing"

	"github.com/cceckman/discoirc/ui"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/cceckman/discoirc/ui/widgets"
	discomocks "github.com/cceckman/discoirc/ui/mocks"
)

func TestActivation(t *testing.T) {
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
}
