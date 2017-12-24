package controller_test

import (
	"context"

	_ "github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel/controller"
	"github.com/cceckman/discoirc/ui/channel/mocks"
	_ "github.com/marcusolsson/tui-go"

	"testing"
)

func TestController_ResizeNoEvents(t *testing.T) {
	ui := mocks.NewUpdateCounter()

	m := &mocks.Model{}
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
	ui := mocks.NewUpdateCounter()

	m := &mocks.Model{
		Events: mocks.Events,
	}
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
	ui := mocks.NewUpdateCounter()

	m := &mocks.Model{
		Events: mocks.Events,
	}
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

