// This file provides
// Eventually, I want this to be "load config from a proto;"
// but for now, just produce things.
package client

import (
	"crypto/tls"

	irc "github.com/fluffle/goirc/client"
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
