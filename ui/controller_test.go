package ui_test

import (
	"testing"

	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/data"
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

func TestEndToEnd(t *testing.T) {
	u := mocks.NewUI()
	surface := tui.NewTestSurface(20, 5)
	u.Update(func() {
		u.Painter = tui.NewPainter(surface, tui.NewTheme())
	})
	be := mocks.NewBackend()

	var ctl *ui.Controller

	// Root creation happens in main thread
	u.Update(func() { ctl = ui.New(u, be) })
	// Client activation happens in response to user events, in the main thread
	u.Update(func() { ctl.ActivateClient() })
	u.Wait()

	be.Receiver.UpdateChannel(data.ChannelState{
		Network:     "HamNet",
		Channel:     "#hamlet",
		ChannelMode: "i",
	})
	// Simulate selection
	u.Type("jj")

	wantContents := `
 HamNet: ?          
|#hamlet           i
|✉ 0             0 ☺
                    
                    
`
	u.RunSync(func() {
		got := surface.String()
		if got != wantContents {
			t.Errorf("unexpected contents:\ngot = \n%s\n--\nwant = \n%s\n--", got, wantContents)
		}

	})

	// Simulate activation
	u.Type("\n")

	u.RunSync(func() {
		if _, ok := u.Root.(*channel.View); !ok {
			t.Errorf("unexpected view at UI root: got: %+v want: client.View", u.Root)
		}

	})

}
