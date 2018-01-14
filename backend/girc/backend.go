package girc

import (
	"fmt"

	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/backend"
)

// Backend provides IRC functionality via the girc library.
type Backend struct {
	Network map[string]*Network

	clientLog data.EventList

	recv backend.StateReceiver
}

func NewBackend() *Backend {
	return &Backend{
		Network: make(map[string]*Network),
	}
}


var _ backend.Backend = &Backend{}
func (b *Backend) EventsBefore(s data.Scope, n int, last data.Seq) data.EventList {
	if s.Net == "" {
		return b.clientLog.SelectSizeMax(n, last)
	}

	net, ok := b.Network[s.Net]
	if ok {
		return net.EventsBefore(s, n, last)
	}

	return nil
}
func (b *Backend) Send(s data.Scope, message string) {
	panic(fmt.Errorf("not implemented"))
}
func (b *Backend) Subscribe(s backend.StateReceiver) {
	b.recv = s
}
