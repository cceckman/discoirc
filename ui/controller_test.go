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
	if _, ok := u.Root.(*channel.View); !ok {
		t.Errorf("unexpected view at UI root: got: %+v want: client.View", u.Root)
	}
}

func TestActivateClient(t *testing.T) {
	u := mocks.NewUI()

	ctl := ui.New(u, mocks.NewBackend())

	ctl.ActivateClient()
	if _, ok := u.Root.(client.View); !ok {
		t.Errorf("unexpected view at UI root: got: %+v want: client.View", u.Root)
	}
}

func TestEndToEnd(t *testing.T) {
	u := mocks.NewUI()
	surface := tui.NewTestSurface(30, 10)
	u.Painter = tui.NewPainter(surface, tui.NewTheme())
	be := mocks.NewBackend()

	ctl := ui.New(u, be)
	ctl.ActivateClient()

	ch := data.ChannelState{
		Network:     "HamNet",
		Channel:     "#hamlet",
		ChannelMode: "i",
		Topic:       "The Battlements",
		LastMessage: mocks.Events[2],
	}
	net := data.NetworkState{
		Network: "HamNet",
		Nick:    "yorick",
		State:   data.Connecting,
	}

	be.Receiver.UpdateChannel(ch)
	be.Receiver.UpdateNetwork(net)
	u.Repaint()

	wantContents := `
 HamNet: …              yorick
 #hamlet                     i
 ✉ 0                       0 ☺
                              
                              
                              
                              
                              
                              
                              
`
	got := surface.String()
	if got != wantContents {
		t.Errorf("unexpected contents:\ngot = \n%s\n--\nwant = \n%s\n--", got, wantContents)
	}

	// Simulate selection
	u.Type("jj")
	u.Repaint()

	wantContents = `
 HamNet: …              yorick
|#hamlet                     i
|✉ 0                       0 ☺
                              
                              
                              
                              
                              
                              
                              
`
	got = surface.String()
	if got != wantContents {
		t.Errorf("unexpected contents:\ngot = \n%s\n--\nwant = \n%s\n--", got, wantContents)
	}

	// Simulate activation
	u.Type("\n")
	u.Repaint()
	wantContents = `
                              
                              
                              
                              
                              
                              
                              
                              
HamNet: ? #hamlet:            
< >                           
`

	if _, ok := u.Root.(*channel.View); !ok {
		t.Errorf("unexpected view at UI root: got: %+v want: client.View", u.Root)
	}
	got = surface.String()
	if got != wantContents {
		t.Errorf("unexpected contents:\ngot = \n%s\n--\nwant = \n%s\n--", got, wantContents)
	}

}
