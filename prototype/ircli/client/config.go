// This file provides
// Eventually, I want this to be "load config from a proto;"
// but for now, just produce things.
package client

import (
	"crypto/tls"
	"fmt"

	irc "github.com/fluffle/goirc/client"
	"sync"
)

// Load returns one config per connection, named by network.
func (_ *client) LoadConfigs() map[string]*irc.Config {
	result := map[string]*irc.Config{}

	var cfg *irc.Config
	var server string

	server = "irc.foonetic.net"
	cfg = irc.NewConfig("eckbot")
	cfg.SSL = true
	cfg.SSLConfig = &tls.Config{
		ServerName: server,
	}
	cfg.Server = server + ":" + "6667"
	result["Foonetic"] = cfg

	return result
}

func (c *client) Connect() <-chan error {
	// Connect to multiple servers in parallel, disco
	errs := make(chan error)
	var wg sync.WaitGroup

	for network, cfg := range c.LoadConfigs() {
		wg.Add(1)
		go func(name string, cfg *irc.Config) {
			cli := irc.Client(cfg)
			c.attachHandlers(name, cli)

			if err := c.Connect(); err != nil {
				errs <- fmt.Errorf("error connecting to %s: %v", name, err)
			}
			// Indicate that the connection is complete.
			wg.Done()
		}(network, cfg)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()
	return errs
}

func (c *client) attachHandlers(name string, cli *irc.Conn) error {
	cli.HandleFunc(
		irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			c.Lock()
			defer c.Unlock()

			c.connections[name] = conn
		},
	)
	cli.HandleFunc(
		irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			c.Lock()
			defer c.Unlock()

			delete(c.connections, name)
		},
	)
	return nil
}

func (c *client) ConnectedNetworks() []string {
	c.Lock()
	defer c.Unlock()

	result := make([]string, len(c.connections))
	i := 0
	for net, _ := range c.connections {
		result[i] = net
		i++
	}
	return result
}

