// Package demo implements the discoirc non-UI portions with demo data.
package demo

import (
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

func (d *Demo) Subscribe(c backend.StateReceiver) {
	d.Lock()
	defer d.Unlock()
	d.receiver = c
	d.filter = nil

	// If the receiver actually exists, update state.
	if d.receiver != nil {
		d.updateAll()
	}
	return
}

func (d *Demo) SubscribeFiltered(recv backend.FilteredStateReceiver) {
	d.Lock()
	defer d.Unlock()
	d.receiver = recv
	d.filter = recv.Filter

	d.updateAll()
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

func (d *Demo) updateNetwork(network string) {
	if d.filter != nil {
		net, _ := d.filter()
		if net != network {
			return
		}
	}

	if d.receiver != nil {
		d.receiver.UpdateNetwork(*d.nets[network])
	}
}

func (d *Demo) TickNetwork(network string) {
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
	d.updateNetwork(network)
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

func (d *Demo) TickChannel(network, channel string) {
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

	d.updateNetwork(network)
	d.updateChannel(network, channel)
}

func (d *Demo) updateChannel(network, channel string) {
	if d.filter != nil {
		net, ch := d.filter()
		if net != network || ch != channel {
			return
		}
	}
	if d.receiver != nil {
		d.receiver.UpdateChannel(*d.chans[ChanIdent{
			Network: network,
			Channel: channel,
		}])
	}
}

func (d *Demo) updateAll() {
	if d.receiver == nil {
		return
	}
	for id, _ := range d.chans {
		d.updateNetwork(id.Network)
		d.updateChannel(id.Network, id.Channel)
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
