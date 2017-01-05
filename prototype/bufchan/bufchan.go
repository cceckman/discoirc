//
// bufchan.go
// Copyright (C) 2016 cceckman <charles@cceckman.com>
//
// Distributed under terms of the MIT license.
//

// Package bufchan provides channels with unlimited buffering.
package bufchan

import (
	"context"
)

// Bufchan is a channel of strings with an unlimited internal buffer.
// In cases where a UI must interact with a background process (e.g. the network),
// it is necessary that writes between the UI and the background are non-blocking.
// Channels - buffered channels in particular - can help alleviate this; however, channels in Go
// still have a limited size, after which time writes will be blocking operations.
// Bufchan allows channel writes to happen in an approximately always non-blocking manner.
type Bufchan interface {
	In() chan<- interface{}
	Out() <-chan interface{}
}

type bc struct {
	Broadcaster
	out <-chan interface{}
}

func (b *bc) In() chan<- interface{} {
	return b.Send()
}

func (b *bc) Out() <-chan interface{} {
	return b.out
}

func New() Bufchan {
	b := NewBroadcaster()
	r := &bc{
		Broadcaster: b,
		out:         b.Listen(context.Background()),
	}
	return r
}
