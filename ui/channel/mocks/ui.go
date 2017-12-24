package mocks

import (
	"sync"
)

// UpdateCounter is a controller.UIUpdater that can queues, and can synchronize against, outstanding requests.
type UpdateCounter struct {
	sync.WaitGroup

	counting bool
	incoming chan func()
}

// NewUpdateCounter returns a new UpdateCounter.
func NewUpdateCounter() *UpdateCounter {
	r := &UpdateCounter{
		incoming: make(chan func(), 1),
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

// Add adds delta to the count of expected updates,
// and enables tracking
func (u *UpdateCounter) Add(delta int) {
	u.Update(func() {
		// Add one, to count off this operation.
		u.counting = true
		u.WaitGroup.Add(delta + 1)
	})
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
