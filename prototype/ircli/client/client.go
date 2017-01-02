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
	Listen(ctx context.Context) <-chan string
	// Set the current target to channel or nick 'target' on the given network.
	// SetTarget(network, target error) error
	// Send(msg string) error
}

type client struct {
	networks map[string]network
	sync.RWMutex

	receive chan<- string
	bufchan.StringBroadcaster
}

type network struct {
	conn *irc.Conn
}

func NewClient() C {
	bc := bufchan.NewStringBroadcaster()
	c := &client{
		StringBroadcaster: bc,
		networks:          make(map[string]network),
		receive:           bc.Send(),
	}

	return c
}

// Handle returns an irc.Handler for messages on the given network.
func (c *client) Handle(network string) irc.HandlerFunc {
	return func(conn *irc.Conn, line *irc.Line) {
		c.RLock()
		defer c.RUnlock()

		_, ok := c.networks[network]
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

	for network, cfg := range c.LoadConfigs() {
		wg.Add(1)
		log.Println("Launching connector for network", network)
		go func(name string, cfg *irc.Config) {
			// Mark connection as completed when we're done here.
			defer wg.Done()

			log.Println("Attempting to connect to", network)
			cli := irc.Client(cfg)
			c.attachHandlers(name, cli)

			if err := cli.Connect(); err != nil {
				e := fmt.Errorf("error connecting to %s: %v", name, err)
				log.Println(e)
				errs <- e
			}
		}(network, cfg)
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

func (c *client) attachHandlers(name string, cli *irc.Conn) error {
	cli.HandleFunc(
		irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			c.receive <- fmt.Sprintf("[%s] Starting connection", name)
			c.Lock()
			defer c.Unlock()

			c.networks[name] = network{conn: conn}
			c.receive <- fmt.Sprintf("[%s] Connected", name)
		},
	)
	cli.HandleFunc(
		irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			c.Lock()
			defer c.Unlock()

			delete(c.networks, name)
			c.receive <- fmt.Sprintf("[%s] Disconnected", name)
		},
	)
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
	defer c.RUnlock()
	// Internal check; make sure they are actually disconnected.
	for name, _ := range c.networks {
		results = append(
			results,
			fmt.Errorf("network %s still present in connected nets table after Close()", name),
		)
	}

	close(c.receive)

	return results

}
