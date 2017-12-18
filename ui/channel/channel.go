// Package channel includes top-level types and interfaces for the channel cnotents UI.
package channel

import (
	"context"
	"github.com/cceckman/discoirc/data"
	"github.com/marcusolsson/tui-go"
)

type Model interface {
  // Includes nick, channelMode, connection state, topic
  Channel(ctx context.Context) <-chan data.Channel
  // Returns up to N events starting at EpochId
  EventsStartingAt(start data.EventID, n int) []data.Event
  // Returns up to N events ending at EpochId
  EventsEndingAt(end data.EventID, n int) []data.Event
  // Streams events starting at EpochId
  Follow(ctx context.Context, start data.EventID) <-chan data.Event
  Send(e data.Event) error
}

type View interface {
  tui.Widget

  SetTopic(string)
  SetNick(string)
  SetChannelMode(string)
	SetConnection(string)
  SetPresence(string)
  // Renders the messages.
  SetMessages([]data.Event)

  // Advanced / deferred: special rendering for widgets,
  // e.g. a closure on a hilighter.
  SetRenderer(func(data.Event) tui.Widget)
}

type Controller interface {
  // Parse input string, do appropriate stuff to backend
  Input(string)

  // Resize indicates the number of lines now available for messages.
  Resize(n int)

  // (Asynchronous) scroll commands
  ScrollUp()
  ScrollDown()

	// TODO: Deferred: Localization of connection / presence
}


