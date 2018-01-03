// Package demo implements the discoirc non-UI portions with demo data.
package demo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/data"
)

var _ backend.Backend = &Demo{}

type ChanIdent struct {
	Network, Channel string
}

// Demo provides data and updates to discoirc UI components.
type Demo struct {
	incoming chan func()

	subscriber backend.StateReceiver
	filter     func() (string, string)

	nets     map[string]*data.NetworkState
	chans    map[ChanIdent]*data.ChannelState
	contents map[ChanIdent][]data.Event
}

func New() *Demo {
	d := &Demo{
		nets:     make(map[string]*data.NetworkState),
		chans:    make(map[ChanIdent]*data.ChannelState),
		contents: make(map[ChanIdent][]data.Event),
		incoming: make(chan func()),
	}
	go d.run()
	return d
}

func (d *Demo) Close() {
	close(d.incoming)
}

func (d *Demo) run() {
	for f := range d.incoming {
		f()
	}
}

func (d *Demo) Subscribe(recv backend.StateReceiver) {
	d.incoming <- func() {
		d.subscriber = recv
		d.filter = nil
		d.updateAll()
	}
}

func (d *Demo) SubscribeFiltered(recv backend.FilteredStateReceiver) {
	d.incoming <- func() {
		d.subscriber = recv
		d.filter = recv.Filter
		d.updateAll()
	}
}

func (d *Demo) TickNetwork(network string) {
	d.incoming <- func() {
		d.tickNetwork(network)
	}
}

func (d *Demo) TickChannel(network, channel string) {
	d.incoming <- func() {
		d.tickChannel(network, channel)
	}
}

func (d *Demo) TickMessages(network, channel string) {
	d.incoming <- func() {
		d.tickMessages(network, channel)
	}
}

func (d *Demo) Send(network, channel string, messages string) {
	d.incoming <- func() {
		d.send(network, channel, messages)
	}
}

// Join waits until all other issued operations (ticks, sends, etc.)
// have completd
//
// This can be used by tests to block until pending operations have completed.
func (d *Demo) Join() {
	blk := make(chan struct{})
	d.incoming <- func() {
		close(blk)
	}
	<-blk
}

func (d *Demo) ensureNetwork(network string) {
	net := d.nets[network]
	if net == nil {
		d.nets[network] = &data.NetworkState{
			Network: network,
		}
		net = d.nets[network]
	}
}

func (d *Demo) ensureChannel(network, channel string) {
	d.ensureNetwork(network)
	chId := ChanIdent{Network: network, Channel: channel}
	ch := d.chans[chId]
	if ch == nil {
		d.chans[chId] = &data.ChannelState{
			Network: network,
			Channel: channel,
			Unread:  0,
		}
		ch = d.chans[chId]
	}
}

func tickNick(nick string) string {
	base := "nicholas"
	suffix := strings.TrimPrefix(nick, base)
	// OK to ignore err; val defaulting to 0 is correct.
	val, _ := strconv.Atoi(suffix)
	return fmt.Sprintf("%s%d", base, val+1)
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

	if d.subscriber != nil {
		d.subscriber.UpdateNetwork(*d.nets[network])
	}
}

func (d *Demo) tickNetwork(network string) {
	d.ensureNetwork(network)
	net := d.nets[network]
	net.State = tickConnState(net.State)
	net.Nick = tickNick(net.Nick)
	d.updateAll()
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

func (d *Demo) tickChannel(network, channel string) {
	d.ensureChannel(network, channel)
	ch := d.chans[ChanIdent{
		Network: network,
		Channel: channel,
	}]
	ch.ChannelMode = tickMode(ch.ChannelMode)
	ch.UserMode = tickMode(ch.UserMode)
	ch.Presence = tickPresence(ch.Presence)
	ch.Topic = tickTopic(ch.Topic)
	ch.Members += 1

	d.updateAll()
}

func (d *Demo) updateChannel(network, channel string) {
	if d.filter != nil {
		net, ch := d.filter()
		if net != network || ch != channel {
			return
		}
	}
	if d.subscriber != nil {
		d.subscriber.UpdateChannel(*d.chans[ChanIdent{
			Network: network,
			Channel: channel,
		}])
	}
}

func (d *Demo) updateAll() {
	for net, _ := range d.nets {
		d.updateNetwork(net)
	}
	for id, _ := range d.chans {
		d.updateChannel(id.Network, id.Channel)
	}
}

var messages = []string{
	"Shall I compare thee to a summer’s day?",
	"Thou art more lovely and more temperate.",
	"Rough winds do shake the darling buds of May,",
	"And summer’s lease hath all too short a date.",
	"Sometime too hot the eye of heaven shines,",
	"And often is his gold complexion dimmed;",
	"And every fair from fair sometime declines,",
	"By chance, or nature’s changing course, untrimmed;",
	"But thy eternal summer shall not fade,",
	"Nor lose possession of that fair thou ow’st,",
	"Nor shall death brag thou wand’rest in his shade,",
	"When in eternal lines to Time thou grow’st.",
	"So long as men can breathe, or eyes can see,",
	"So long lives this, and this gives life to thee.",
}
var speakers = []string{
	"troilus",
	"cressida",
	"aeneas",
	"dido",
	"antonio",
	"sebastian",
	"gentleman2",
}

func (d *Demo) tickMessages(network, channel string) {
	d.ensureChannel(network, channel)
	id := ChanIdent{
		Network: network,
		Channel: channel,
	}

	// Construct a message using the absolute sequence number; as if
	// each character were reciting the sonnet in turn.
	seq := len(d.contents[id])
	msg := messages[seq%len(messages)]

	// what iteration of the sonnet are we on?
	iteration := seq / len(messages)
	// different speaker for each iteration.
	speaker := speakers[iteration%len(speakers)]
	d.appendMessage(network, channel, speaker, msg)
}

func (d *Demo) send(network, channel string, message string) {
	d.ensureNetwork(network)
	nick := d.nets[network].Nick
	d.appendMessage(network, channel, nick, message)
}

func (d *Demo) appendMessage(network, channel string, speaker, message string) {
	d.ensureChannel(network, channel)
	id := ChanIdent{
		Network: network,
		Channel: channel,
	}
	last := d.chans[id].LastMessage.EventID
	var next data.Event

	// Construct a new message.
	next.Seq = last.Seq + 1
	next.Epoch = last.Epoch
	next.Contents = fmt.Sprintf(
		"<%s> %s",
		speaker, message,
	)

	d.contents[id] = append(d.contents[id], next)
	d.chans[id].Unread += 1
	d.chans[id].LastMessage = next

	d.updateAll()
}

func (d *Demo) EventsBefore(network, channel string, n int, last data.EventID) []data.Event {
	id := ChanIdent{
		Network: network,
		Channel: channel,
	}
	r := make(chan []data.Event)
	d.incoming <- func() {
		r <- data.NewEvents(d.contents[id]).SelectSizeMax(uint(n), last)
		close(r)
		if v, ok := d.chans[id]; ok {
			v.Unread = 0
			d.updateAll()
		}
	}

	return <-r
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
