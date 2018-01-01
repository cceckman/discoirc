package ui_test

import (
	"testing"

	"github.com/cceckman/discoirc/ui"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/cceckman/discoirc/ui/client"
	"github.com/cceckman/discoirc/ui/mocks"
)

func TestActivateChannel(t *testing.T) {
	u := mocks.NewUI()

	ctl := ui.New(u, mocks.NewBackend())

	ctl.ActivateChannel("foonet", "#barchan")
	u.RunSync(func() {
		if _, ok := u.Root.(*channel.View); !ok {
			t.Errorf("unexpected view at UI root: got: %+v want: client.View", u.Root)
		}
	})
}

func TestActivateClient(t *testing.T) {
	u := mocks.NewUI()

	ctl := ui.New(u, mocks.NewBackend())

	ctl.ActivateClient()
	u.RunSync(func() {
		if _, ok := u.Root.(client.View); !ok {
			t.Errorf("unexpected view at UI root: got: %+v want: client.View", u.Root)
		}
	})
}
