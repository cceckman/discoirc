package mocks

import (
	"sync"

	"github.com/marcusolsson/tui-go"
)

// UI implements a subset of the tui.UI functionality for use in tests.
type UI struct {
	UpdateCounter

	root tui.Widget
}

func (ui *UI) SetWidget(w tui.Widget) {
	ui.root = w
}

// GetWidget safely gets the root Widget, after all pending updates have run.
func (ui *UI) GetWidget() tui.Widget {
	r := make(chan tui.Widget)
	ui.RunSync(func() {
		r <-ui.root
	})
	return <-r
}

// UpdateCounter is a controller.UIUpdater that can queues, and can synchronize against, outstanding requests.
type UpdateCounter struct {
	sync.WaitGroup

	startCounting sync.Once
	counting      chan struct{}
	incoming      chan func()
}

// NewUpdateCounter returns a new UpdateCounter.
func NewUpdateCounter() *UpdateCounter {
	r := &UpdateCounter{
		counting: make(chan struct{}),
		incoming: make(chan func(), 1),
	}
	go func() {
		for f := range r.incoming {
			f()
			select {
			case _ = <-r.counting:
				r.Done()
			default:
				// do nothing
			}
		}
	}()
	return r
}

// Add adds delta to the count of expected updates,
// and enables tracking
func (u *UpdateCounter) Add(delta int) {
	u.startCounting.Do(func() { close(u.counting) })
	u.WaitGroup.Add(delta)
}

// Update queues f to run at a later time.
func (u *UpdateCounter) Update(f func()) {
	u.incoming <- f
}

// RunSync waits until all other updates are complete, then runs itself and waits for completion.
func (u *UpdateCounter) RunSync(f func()) {
	u.Wait()
	u.Add(1)
	u.Update(f)
	u.Wait()
}
