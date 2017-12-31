package mocks

import (
	"sync"
)

// UpdateCounter is a controller.UIUpdater that can queues, and can synchronize against, outstanding requests.
type UpdateCounter struct {
	wg sync.WaitGroup

	incoming chan func()
}

// NewUpdateCounter returns a new UpdateCounter.
func NewUpdateCounter() *UpdateCounter {
	u := &UpdateCounter{
		incoming: make(chan func(), 1),
	}
	go func() {
		for f := range u.incoming {
			f()
			u.wg.Done()
		}
	}()
	return u
}

// Update queues f to run at a later time.
// Updates are processed in the order they are received.
func (u *UpdateCounter) Update(f func()) {
	u.wg.Add(1)
	u.incoming <- f
}

// Wait blocks until all pending Updates have run.
func (u *UpdateCounter) Wait() {
	u.wg.Wait()
}

// Close cleans up this UpdateCounter.
func (u *UpdateCounter) Close() {
	close(u.incoming)
}

// RunSync runs its callback as an Update, but waits for it to complete.
// Note that it does not wait for *all* updates to complete.
func (u *UpdateCounter) RunSync(f func()) {
	blk := make(chan struct{})
	u.Update(func() {
		f()
		close(blk)
	})
	<-blk
}
