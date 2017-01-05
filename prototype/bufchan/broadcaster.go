package bufchan

import (
	"context"
)

// A Broadcaster duplicates values Send()ed to it to all Listen()ers.
type Broadcaster interface {
	Send() chan<- interface{}
	Listen(context.Context) <-chan interface{}
}

/*
Based on @rogpeppe and @keroproxx520 at
https://rogpeppe.wordpress.com/2009/12/01/concurrent-idioms-1-broadcasting-values-in-go-with-linked-channels/
and its comments.
*/

//
type receiver struct {
	next chan broadcast
}

type broadcast struct {
	next chan broadcast
	v    interface{}
}

type broadcaster struct {
	next  chan broadcast
	sendc chan<- interface{}
}

func NewBroadcaster() Broadcaster {
	next := make(chan broadcast, 1)
	sendc := make(chan interface{})
	b := &broadcaster{
		next:  next,
		sendc: sendc,
	}
	go func() {
		for v := range sendc {
			c := make(chan broadcast, 1)
			newb := broadcast{next: c, v: v}
			b.next <- newb
			b.next = c
		}
		// Input channel closed. Send empty value to listeners.
		// The channel won't be closed, though it may be GC'd.
		b.next <- broadcast{}
	}()

	return b
}

func (b *broadcaster) Send() chan<- interface{} {
	return b.sendc
}

func (b *broadcaster) Listen(ctx context.Context) <-chan interface{} {
	out := make(chan interface{})
	// TODO: is there not a race condition in reading this value?
	r := receiver{b.next}

	go func() {
		defer close(out)

		for {
			// Wait for either cancellation, or a new Broadcast to come down the pipe.
			select {
			case <-ctx.Done():
				return
			case b := <-r.next:
				// New value came down the pipe.
				// Take the value,
				v := b.v
				// give the broadcast back, for the next listener to take,
				r.next <- b
				// and update ourselves.
				r.next = b.next

				// if 'next' is nil, we're done.
				if b.next == nil {
					return
				}

				// Otherwise, send along to *our* listener.

				select {
				case <-ctx.Done():
					return
				case out <- v:
					// sent the value to the listener.
				}

			}
		}
	}()

	return out
}
