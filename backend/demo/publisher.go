// Package demo implements the discoirc non-UI portions with demo data.
package demo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/data"
)

type ChanIdent struct {
	Network, Channel string
}

// Demo provides data and updates to discoirc UI components.
type Demo struct {
	sync.Mutex

	receiver backend.StateReceiver
	filter   func() (string, string)

	nets  map[string]*data.NetworkState
	chans map[ChanIdent]*data.ChannelState
}

func New() *Demo {
	return &Demo{
		nets:  make(map[string]*data.NetworkState),
		chans: make(map[ChanIdent]*data.ChannelState),
	}
}

func (d *Demo) Subscribe(ctx context.Context, c backend.StateReceiver) {
	d.Lock()
	defer d.Unlock()
	d.receiver = c

	go func() {
		<-ctx.Done()
		d.Lock()
		defer d.Unlock()
		d.receiver = nil
	}()

	// If the context is still active; send initial state.
	select {
	case <-ctx.Done():
		return
	default:
	}
	for _, v := range d.nets {
		d.receiver.UpdateNetwork(*v)
	}
	for _, v := range d.chans {
		d.receiver.UpdateChannel(*v)
	}

	return
}

func tickNick(nick string) string {
	base := "nicholas"
	suffix := strings.TrimPrefix(nick, base)
	// OK to ignore err; val defaulting to 0 is correct.
	val, _ := strconv.Atoi(suffix)
	return fmt.Sprintf("%s%d", base, val)
}

func tickConnState(in data.ConnectionState) data.ConnectionState {
	switch in {
	case data.Disconnected:
		return data.Connecting
	case data.Connecting:
		return data.Connected
	case data.Connected:
		return data.Disconnected
	}
	return in
}

func (d *Demo) updateNetwork(ctx context.Context, network string) {
	// If the context is still active, send update to client.
	select {
	case <-ctx.Done():
		return
	default:
		if d.receiver != nil {
			d.receiver.UpdateNetwork(*d.nets[network])
		}
	}

}

func (d *Demo) TickNetwork(ctx context.Context, network string) {
	d.Lock()
	defer d.Unlock()
	// Update internal state.

	net := d.nets[network]
	if net == nil {
		d.nets[network] = &data.NetworkState{
			Network: network,
		}
		net = d.nets[network]
	}
	net.State = tickConnState(net.State)
	net.Nick = tickNick(net.Nick)
	d.updateNetwork(ctx, network)
}

func tickUMode(m string) string {
	modes := []string{"q", "a", "o", "h", "v", ""}
	for i, s := range modes {
		if m == s {
			return modes[(i+1)%len(modes)]
		}
	}
	return modes[0]
}

func tickMode(m string) string {
	modes := []string{"i", "k", "l", "s", ""}
	for i, s := range modes {
		if m == s {
			return modes[(i+1)%len(modes)]
		}
	}
	return modes[0]
}

func tickPresence(p data.Presence) data.Presence {
	switch p {
	case data.Joined:
		return data.NotPresent
	case data.NotPresent:
		return data.Joined
	}
	return data.NotPresent
}

func tickTopic(t string) string {
	topic := strings.Split("The Tragical History of the Life and Death of Doctor Faustus", " ")
	l := len(strings.Split(t, " "))
	l = ((l + 1) % len(topic)) + 1
	return strings.Join(topic[0:l], " ")
}

func (d *Demo) TickChannel(ctx context.Context, network, channel string) {
	d.Lock()
	defer d.Unlock()

	net := d.nets[network]
	if net == nil {
		d.nets[network] = &data.NetworkState{
			Network: network,
		}
		net = d.nets[network]
	}
	chId := ChanIdent{Network: network, Channel: channel}
	ch := d.chans[chId]
	if ch == nil {
		d.chans[chId] = &data.ChannelState{
			Network: network,
			Channel: channel,
			Unread:  1,
		}
		ch = d.chans[chId]
	}
	ch.ChannelMode = tickMode(ch.ChannelMode)
	ch.UserMode = tickMode(ch.UserMode)
	ch.Presence = tickPresence(ch.Presence)
	ch.Topic = tickTopic(ch.Topic)
	ch.Members += 1
	ch.Unread *= 2

	d.updateNetwork(ctx, network)
	d.updateChannel(ctx, network, channel)
}

func (d *Demo) updateChannel(ctx context.Context, network, channel string) {
	// If the context is still active, send update to client.
	select {
	case <-ctx.Done():
		return
	default:
		if d.receiver != nil {
			d.receiver.UpdateChannel(*d.chans[ChanIdent{
				Network: network,
				Channel: channel,
			}])
		}
	}

}

/*
Inline sketch:

Tests assert that no updates come in after the context is cancelled- essentially,
that ticks are synchronous.

That may be possible to accomplish, as - in the backend/ package - we're not
expecting multiple subscribers to actually need events; we only have one window
open at a time. A subprocess backend's subprocess would need to support that,
but that's a different sort of problem.
*/
