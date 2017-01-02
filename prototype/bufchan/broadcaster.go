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
		// Input channel closed; propagate to clients with an empty broadcast.
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
				// update our pointer to the 'next' channel.
				r.next = b.next

				if r.next == nil {
					// There will be no 'next' value, i.e. input channel is closed.
					return
				}

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

// StringBroadcaster is a Broadcaster with methods wrapped to use string channels.
type StringBroadcaster interface {
	Send() chan<- string
	Listen(context.Context) <-chan string
}

type stringBroadcaster struct {
	ssend chan string
	b     Broadcaster
}

func NewStringBroadcaster() StringBroadcaster {
	r := &stringBroadcaster{
		ssend: make(chan string),
		b:     NewBroadcaster(),
	}
	// Mirror sending channel.
	go func() {
		send := r.Send()
		defer close(send)
		for str := range r.ssend {
			send <- str
		}
	}()

	return r
}

func (s *stringBroadcaster) Send() chan<- string {
	return s.ssend
}

func (s *stringBroadcaster) Listen(ctx context.Context) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for v := range s.b.Listen(ctx) {
			// Yes, panic if it's not a string. Only way to send is via the typed Send, above.
			out <- v.(string)
		}
	}()

	return out
}
