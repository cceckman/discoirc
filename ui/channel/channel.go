// Package channel includes top-level types and interfaces for the channel cnotents UI.
package channel

import (
	"context"
	"github.com/cceckman/discoirc/data"
	"github.com/marcusolsson/tui-go"
)

type Model interface {
  // Includes nick, channelMode, connection state, topic
  ConnectionState(ctx context.Context) <-data.ChannelState
  // Returns up to N events starting at EpochId
  EventsStartingAt(start data.EventID, n int) []Event
  // Returns up to N events ending at EpochId
  EventsEndingAt(end data.EventID, n int) []Event
  // Streams events starting at EpochId
  Follow(start data.EventID) <-data.Event
  Send(e data.Event) error
}

type View interface {
  tui.Widget

  SetTopic(string)
  SetNick(string)
  SetChannelMode(string)
  SetConnected(string)
  // Renders the messages.
  SetMessages([]Event)

  // Advanced / deferred: special rendering for widgets,
  // e.g. a closure on a hilighter.
  SetRenderer(func(Event) tui.Widget)
}

type Controller interface {
  // Parse input string, do appropriate stuff to backend
  Input(string)

  // Resize indicates the number of lines now available for messages.
  Resize(n int)

  // (Asynchronous) scroll commands
  ScrollUp()
  ScrollDown()
}


