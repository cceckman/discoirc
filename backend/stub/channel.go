package stub

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
)

func NewChannel(ctx context.Context) *Channel {
	r := &Channel{
		received:   make(chan string, 1),
		metaUpdate: make(chan struct{}, 1),
	}
	go r.awaitMessages(ctx)
	go r.awaitMeta(ctx)
	return r
}

var names = []string{"hamlet", "othello", "macbeth"}
var quotes = []string{
	"o that this too too solid flesh would melt",
	"to post with such dexterity to incestuous sheets",
	"she had eyes and chose me",
	"tomorrow and tomorrow and tomorrow",
}

type Channel struct {
	received   chan string
	metaUpdate chan struct{}

	sync.RWMutex
	controller channel.ModelController
	meta       data.Channel
	events     data.EventList
}

// GenerateMeta generates metadata updates at a given rate, until the context is cancelled.
func (c *Channel) GenerateMeta(ctx context.Context, d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for i, j := 0, 0; true; i, j = i+1%len(names), j+1%len(quotes) {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m := c.GetMeta()
			m.Topic = quotes[j]
			m.Mode = fmt.Sprintf("+%s", names[i])
			switch m.Presence {
			case data.Joined:
				m.Presence = data.NotPresent
			case data.NotPresent:
				m.Presence = data.Joined
			}
			switch m.Connection.State {
			case data.Disconnected:
				m.Connection.State = data.Connecting
			case data.Connecting:
				m.Connection.State = data.Connected
			case data.Connected:
				m.Connection.State = data.Disconnected
			}
			c.SetMeta(m)
		}
	}

}

// Generate generates fake messages at a given rate, until the context is cancelled.
func (c *Channel) GenerateMessages(ctx context.Context, d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for i, j := 0, 0; true; i, j = i+1%len(names), j+1%len(quotes) {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.SendFor(names[i], quotes[j])
		}
	}
}

// SetMeta sets the channel's metadata.
func (c *Channel) SetMeta(d data.Channel) {
	c.Lock()
	defer c.Unlock()
	c.meta = d
	select {
	case c.metaUpdate <- struct{}{}:
	case <-c.metaUpdate:
		c.metaUpdate <- struct{}{}
	}
}

// GetMeta gets the current metadata of the channel.
func (c *Channel) GetMeta() data.Channel {
	c.RLock()
	defer c.RUnlock()
	return c.meta
}

func (c *Channel) awaitMeta(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.metaUpdate:
			func() {
				c.RLock()
				defer c.RUnlock()
				c.controller.UpdateMeta(c.meta)
			}()
		}
	}
}

func (c *Channel) awaitMessages(ctx context.Context) {
	epoch := 0
	seq := uint(0)
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c.received:
			e := data.Event{
				EventID:  data.EventID{Epoch: epoch, Seq: seq},
				Contents: msg,
			}
			func() {
				c.Lock()
				defer c.Unlock()
				c.events = data.NewEvents(append(c.events, e))
				if c.controller != nil {
					c.controller.UpdateContents(e)
				}
			}()
			seq++
		}
	}
}

func (c *Channel) EventsEndingAt(end data.EventID, n int) []data.Event {
	c.RLock()
	defer c.RUnlock()
	return c.events.SelectSizeMax(uint(n), end)
}

func (c *Channel) Attach(m channel.ModelController) {
	c.Lock()
	defer c.Unlock()
	c.controller = m

	c.controller.UpdateMeta(c.meta)
	if len(c.events) > 0 {
		c.controller.UpdateContents(c.events[len(c.events)-1])
	}
}

func (c *Channel) Send(msg string) error {
	return c.SendFor(c.meta.Connection.Nick, msg)
}

func (c *Channel) SendFor(nick, msg string) error {
	c.received <- fmt.Sprintf("<%s> %s", nick, msg)
	return nil
}
