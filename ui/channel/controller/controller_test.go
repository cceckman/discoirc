package controller_test

import (
	"context"
	"sync"

	_ "github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel/controller"
	"github.com/cceckman/discoirc/ui/channel/mocks"
	_ "github.com/marcusolsson/tui-go"

	"testing"
)

// UpdateCounter is a controller.UIUpdater that completes one task whenever an Update is completed.
type UpdateCounter struct {
	sync.WaitGroup

	counting bool
	incoming chan func()
}

func NewUpdateCounter() *UpdateCounter {
	r := &UpdateCounter{
		incoming: make(chan func()),
	}
	go func() {
		for f := range r.incoming {
			f()
			if r.counting {
				r.Done()
			}
		}
	}()
	return r
}

func (u *UpdateCounter) Add(delta int) {
	u.Update(func() {
		// Add one, to count off this operation.
		u.counting = true
		u.WaitGroup.Add(delta + 1)
	})
}

func (u *UpdateCounter) Update(f func()) {
	u.incoming <- f
}

// RunSync runs the enclosed method in the same thread as other updates, but waits until it completes.
func (u *UpdateCounter) RunSync(f func()) {
	u.Wait()
	u.Add(1)
	u.Update(f)
	u.Wait()
}

func TestController_Resize(t *testing.T) {
	ui := NewUpdateCounter()

	m := &mocks.Model{}
	v := &mocks.View{}
	_ = controller.New(context.Background(), ui, v, m)

	ui.Add(2)
	// Resize when no events available.
	v.Controller.Resize(10)

	ui.RunSync(func() {
		if len(v.Events) > 0 {
			t.Errorf("wrong number of events: got: %v want: none", v.Events)
		}
	})
}
