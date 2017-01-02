// Package client provides a multi-server IRC client.
package client

import (
	"context"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"log"
	"sync"

	"github.com/cceckman/discoirc/prototype/bufchan"
)

// C is a client connected to (potentially several) IRC servers and rooms.
type C interface {
	LoadConfigs() map[string]*irc.Config
	Connect() []error
	ConnectedNetworks() []string
	Disconnect() []error

	// Listen returns a channel on which messages are relayed.
	Listen(ctx context.Context) <-chan interface{}
	// Set the current target to channel or nick 'target' on the given network.
	// SetTarget(network, target error) error
	// Send(msg string) error
}

type client struct {
	networks map[string]network
	sync.RWMutex

	receive chan<- interface{}
	bufchan.Broadcaster
}

type network struct {
	conn *irc.Conn
}

func NewClient() C {
	bc := bufchan.NewBroadcaster()
	c := &client{
		Broadcaster: bc,
		networks:          make(map[string]network),
		receive:           bc.Send(),
	}

	return c
}

// Handle returns an irc.Handler for messages on the given network.
func (c *client) Handle(network string) irc.HandlerFunc {
	return func(conn *irc.Conn, line *irc.Line) {
		c.RLock()
		_, ok := c.networks[network]
		c.RUnlock()
		if !ok {
			// Disconnected before this handler was processed. Return silently.
			return
		}

		text := fmt.Sprintf("[%s] %s to %s : %s", network, line.Cmd, line.Target(), line.Text())
		c.receive <- text
	}
}

func (c *client) Connect() []error {
	// Connect to multiple servers in parallel, collect errors.
	errs := make(chan error)
	var wg sync.WaitGroup

	for name, cfg := range c.LoadConfigs() {
		wg.Add(1)
		log.Println("Launching connector for network", name)
		go func(name string, cfg *irc.Config) {
			// Mark connection as completed when we're done here.
			defer wg.Done()

			log.Println("Attempting to connect to", name)
			conn := irc.Client(cfg)

			c.Lock()
			c.networks[name] = network{conn: conn}
			c.Unlock()

			c.attachHandlers(name, conn)

			if err := conn.Connect(); err != nil {
				e := fmt.Errorf("error connecting to %s: %v", name, err)
				log.Println(e)
				errs <- e
			}
		}(name, cfg)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	var results []error
	for err := range errs {
		results = append(results, err)
	}

	return results
}

func (c *client) attachHandlers(name string, conn *irc.Conn) error {
	conn.HandleFunc(
		irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			msg := fmt.Sprintf("[%s] Connected", name)
			log.Println(msg)

			// Shouldn't have to background this, since receive should be ~non-blocking.
			// TODO put back in the event thread.
			go func() {
				c.receive <- msg
				log.Printf("[%s] wrote to receive channel", name)
			}()
		},
	)
	conn.HandleFunc(
		irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			msg := fmt.Sprintf("[%s] Disconnected", name)
			log.Println(msg)

			c.Lock()
			delete(c.networks, name)
			c.Unlock()

			// Shouldn't have to background this, since receive should be ~non-blocking.
			// TODO put back in the event thread.
			go func() {
				c.receive <- msg
				log.Printf("[%s] wrote to receive channel", name)
			}()
		},
	)

	handle := c.Handle(name)
	for _, event := range []string{
		irc.REGISTER, irc.CAP, irc.CTCP, irc.CTCPREPLY, irc.ERROR,
		irc.MODE, irc.NOTICE, irc.OPER, irc.PASS, irc.PING, irc.PONG,
		irc.PRIVMSG, irc.QUIT, irc.USER, irc.VERSION, irc.VHOST, irc.WHO,
		irc.WHOIS,
	} {
		// conn.Handle(event, handle)
		_ = event
		_ = handle
	}
	return nil
}

func (c *client) ConnectedNetworks() []string {
	c.RLock()
	defer c.RUnlock()

	result := make([]string, len(c.networks))
	i := 0
	for name, _ := range c.networks {
		result[i] = name
		i++
	}
	return result
}

func (c *client) Disconnect() []error {
	errs := make(chan error)
	var wg sync.WaitGroup

	for name, net := range c.networks {
		// close internal variables
		name, net := name, net
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Close() presumably invokes the "disconnected" handler
			if err := net.conn.Close(); err != nil {
				errs <- fmt.Errorf("error disconnecting from %s: %v", name, err)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	var results []error
	for err := range errs {
		results = append(results, err)
	}

	c.RLock()
	// Internal check; make sure they are actually disconnected.
	for name, _ := range c.networks {
		results = append(
			results,
			fmt.Errorf("network %s still present in connected nets table after Close()", name),
		)
	}
	c.RUnlock()

	close(c.receive)

	return results

}
