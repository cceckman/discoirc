package demo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cceckman/discoirc/data"
)

func (d *Demo) ensureNetwork(network string) {
	d.Lock()
	defer d.Unlock()

	scope := data.Scope{ Net: network }

	if _, ok := d.nets[scope]; !ok {
		d.nets[scope] = &data.NetworkState{}
	}
}

func (d *Demo) ensureChannel(scope data.Scope) {
	d.ensureNetwork(scope.Net)

	d.Lock()
	defer d.Unlock()

	ch := d.chans[scope]
	if ch == nil {
		d.chans[scope] = &data.ChannelState{
			Unread: 0,
		}
	}
}

// TickNetwork increments the values of the given network.
func (d *Demo) TickNetwork(network string) {
	d.ensureNetwork(network)

	d.Lock()
	defer d.Unlock()

	scope := data.Scope{ Net: network }
	net := d.nets[scope]
	net.State = nextConnState(net.State)
	net.Nick = nextNick(net.Nick)

	go d.updateAll()
}

// TickChannel increments the values of the given channel.
func (d *Demo) TickChannel(network, channel string) {
	scope := data.Scope{
		Net:  network,
		Name: channel,
	}
	d.ensureChannel(scope)

	ch := d.chans[scope]
	ch.Mode = nextMode(ch.Mode)
	ch.Presence = nextPresence(ch.Presence)
	ch.Topic = nextTopic(ch.Topic)
	ch.Members++

	go d.updateAll()
}

// TickMessages adds a message to the channel.
func (d *Demo) TickMessages(network, channel string) {
	scope := data.Scope{
		Net:  network,
		Name: channel,
	}
	d.ensureChannel(scope)

	d.Lock()
	defer d.Unlock()

	// Construct a message using the absolute sequence number; as if
	// each character were reciting the sonnet in turn.
	seq := len(d.contents[scope])
	msg := messages[seq%len(messages)]

	// what iteration of the sonnet are we on?
	iteration := seq / len(messages)
	// different speaker for each iteration.
	speaker := speakers[iteration%len(speakers)]

	// Update unread before appending;
	// only these messages may count as unread.
	d.chans[scope].Unread++
	d.appendMessage(scope, speaker, msg)
}

func nextNick(nick string) string {
	base := "nicholas"
	suffix := strings.TrimPrefix(nick, base)
	// OK to ignore err; val defaulting to 0 is correct.
	val, _ := strconv.Atoi(suffix)
	return fmt.Sprintf("%s%d", base, val+1)
}

func nextConnState(in data.ConnectionState) data.ConnectionState {
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

func nextMode(m string) string {
	modes := []string{"i", "k", "l", "s", ""}
	for i, s := range modes {
		if m == s {
			return modes[(i+1)%len(modes)]
		}
	}
	return modes[0]
}

func nextPresence(p data.Presence) data.Presence {
	switch p {
	case data.Joined:
		return data.NotPresent
	case data.NotPresent:
		return data.Joined
	}
	return data.NotPresent
}

func nextTopic(t string) string {
	topic := strings.Split("The Tragical History of the Life and Death of Doctor Faustus", " ")
	l := len(strings.Split(t, " "))
	l = ((l) % len(topic)) + 1
	return strings.Join(topic[0:l], " ")
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
