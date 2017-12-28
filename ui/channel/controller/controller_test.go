package controller_test

import (
	"context"
	"strings"

	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel/controller"
	"github.com/cceckman/discoirc/ui/channel/mocks"
	discomocks "github.com/cceckman/discoirc/ui/mocks"

	"testing"
)

func TestController_ResizeNoEvents(t *testing.T) {
	ui := discomocks.NewController()

	m := mocks.NewModel()
	v := &mocks.View{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ui.Add(1) // Update of metadata on attachment.
	_ = controller.New(ctx, ui, v, m)
	ui.RunSync(func() {
		if len(v.Events) > 0 {
			t.Errorf("wrong number of events: got: %v want: none", v.Events)
		}
	})

	// Resize when no events available; should trigger write of zero events.
	ui.Add(1)
	v.Controller.Resize(8)
	ui.RunSync(func() {
		if len(v.Events) > 0 {
			t.Errorf("wrong number of events: got: %v want: none", v.Events)
		}
	})
}

func TestController_ResizeWithEvents(t *testing.T) {
	ui := discomocks.NewController()

	m := mocks.NewModel()
	m.Events = mocks.Events
	v := &mocks.View{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ui.Add(2) // Update of metadata, contents on attachment.
	_ = controller.New(ctx, ui, v, m)
	ui.RunSync(func() {
		// Still should be zero events; size is zero.
		if len(v.Events) > 0 {
			t.Errorf("wrong number of events: got: %v want: none", v.Events)
		}
	})

	// Resizing should pick up N events.
	ui.Add(1)
	v.Controller.Resize(8)
	ui.RunSync(func() {
		if len(v.Events) != 8 {
			t.Errorf("wrong number of events: got: %d want: %d", len(v.Events), 8)
			return
		}
		gotLast := v.Events[len(v.Events)-1].EventID
		wantLast := m.Events[len(m.Events)-1].EventID
		if gotLast != wantLast {
			t.Errorf("wrong last event: got: %v want: %v", gotLast, wantLast)
		}
	})
}

func TestController_ReceiveEvent(t *testing.T) {
	ui := discomocks.NewController()

	m := mocks.NewModel()
	m.Events = mocks.Events
	v := &mocks.View{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ui.Add(2)
	_ = controller.New(ctx, ui, v, m)
	ui.Add(1)
	sz := 8
	v.Controller.Resize(sz)
	ui.RunSync(func() {
		if len(v.Events) != sz {
			t.Errorf("wrong number of events: got: %d want: %d", len(v.Events), sz)
		}
	})

	ui.Add(1)
	message := "my message"
	m.AddEvent(message)
	ui.RunSync(func() {
		if len(v.Events) != sz {
			t.Errorf("wrong number of events: got: %d want: %d", len(v.Events), sz)
		}
		lastContents := v.Events[len(v.Events)-1].Contents
		if lastContents != message {
			t.Errorf("wrong contents of last message: got: %q want: %q", lastContents, message)
		}
	})
}

func TestController_UpdateMeta(t *testing.T) {
	ui := discomocks.NewController()

	m := mocks.NewModel()
	m.Channel = data.Channel{
		Name: "#discoirc",
		Connection: data.Connection{
			Network: "Foonetic",
			Nick:    "discobot3k",
			State:   data.Connecting,
		},
		Topic:    "the IRC client of the past",
		Presence: data.Joined,
		Mode:     "+pdq",
	}
	v := &mocks.View{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	compare := func(field, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("unexpected value for %s: got: %q want: %q", field, got, want)
		}
	}

	ui.Add(1) // Update channel metadata on attachment
	_ = controller.New(ctx, ui, v, m)
	ui.RunSync(func() {
		compare("Topic", v.Topic, m.Channel.Topic)
		compare("Name", v.Name, m.Channel.Name)
		compare("Mode", v.Mode, m.Channel.Mode)
		compare("Nick", v.Nick, m.Channel.Connection.Nick)

		if !strings.Contains(v.Connection, m.Channel.Connection.Network) {
			t.Errorf("missing connection name %q from contents %q", m.Channel.Connection.Network, v.Connection)
		}
		if !strings.Contains(v.Connection, "…") {
			t.Errorf("missing connection state %q from contents %q", m.Channel.Connection.Network, "…")
		}
	})
}

func TestController_Send(t *testing.T) {
	ui := discomocks.NewController()

	m := mocks.NewModel()
	v := &mocks.View{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = controller.New(ctx, ui, v, m)

	msg := "this message"
	v.Controller.Input(msg)
	got := <-m.Received
	if got != msg {
		t.Errorf("expected message to be passed along: got: %q want: %q", got, msg)
	}
}

func TestController_Quit(t *testing.T) {
	ui := discomocks.NewController()

	m := mocks.NewModel()
	v := &mocks.View{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = controller.New(ctx, ui, v, m)

	ui.Add(1)
	v.Controller.Input("/quit")
	ui.RunSync(func() {
		if !ui.HasQuit {
			t.Errorf("unexpected state: UI has not quit")
		}
	})

}

// TestController_Client tests jumping to the client view.
func TestController_Client(t *testing.T) {
	ui := discomocks.NewController()
	ui.V = discomocks.ChannelView

	m := mocks.NewModel()
	v := &mocks.View{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ui.Add(1) // initial set-root
	_ = controller.New(ctx, ui, v, m)
	ui.RunSync(func() {
		if ui.Root != v {
			t.Errorf("unexpected root: got: %v want: %v", ui.Root, v)
		}
	})

	// TODO: support a keybinding rather than command
	ui.Add(1)
	v.Controller.Input("/Client please?")
	var got discomocks.ActiveView
	ui.RunSync(func() {
		got = ui.V
	})
	want := discomocks.ClientView
	if got != discomocks.ClientView {
		t.Errorf("unexpected active view: got: %v want: %d", got, want)
	}

}
