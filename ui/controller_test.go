package ui_test

import (
	"testing"

	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/cceckman/discoirc/ui/client"
	"github.com/cceckman/discoirc/ui/testhelper"
)

func TestActivateChannel(t *testing.T) {
	t.Parallel()
	u := testhelper.NewUI()

	ctl := ui.New(u, testhelper.NewBackend())

	ctl.ActivateChannel("foonet", "#barchan")
	if _, ok := u.Root.(*channel.View); !ok {
		t.Errorf("unexpected view at UI root: got: %+v want: client.View", u.Root)
	}
}

func TestActivateClient(t *testing.T) {
	t.Parallel()
	u := testhelper.NewUI()

	ctl := ui.New(u, testhelper.NewBackend())

	ctl.ActivateClient()
	if _, ok := u.Root.(client.View); !ok {
		t.Errorf("unexpected view at UI root: got: %+v want: client.View", u.Root)
	}
}

func TestEndToEnd(t *testing.T) {
	t.Parallel()
	u := testhelper.NewUI()
	surface := tui.NewTestSurface(30, 10)
	u.Painter = tui.NewPainter(surface, tui.NewTheme())
	be := testhelper.NewBackend()

	ctl := ui.New(u, be)
	ctl.ActivateClient()

	ch := &data.ChannelStateEvent{
		EventID: data.EventID{
			Scope: data.Scope{
				Net:  "HamNet",
				Name: "#hamlet",
			},
		},
		ChannelState: data.ChannelState{
			Mode:        "i",
			Topic:       "The Battlements",
			LastMessage: testhelper.Events[2].ID().Seq,
		},
	}
	net := &data.NetworkStateEvent{
		EventID: data.EventID{
			Scope: data.Scope{Net: "HamNet"},
		},
		NetworkState: data.NetworkState{
			Nick:  "yorick",
			State: data.Connecting,
		},
	}
	be.Receiver.Receive(ch)
	be.Receiver.Receive(net)
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
