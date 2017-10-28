package model

import (
	"sort"
)

// Client is the main interface to the DiscoIRC client.
// It hosts connections and their properties.
type Client interface {
	Connection(name string) Connection
	Connections() []string
}

type DumbClient map[string]Connection

func (f DumbClient) Connection(name string) Connection {
	return f[name]
}

func (f DumbClient) Connections() []string {
	result := make([]string, len(f))
	i := 0
	for k := range f {
		result[i] = k
		i++
	}
	sort.Strings(result)
	return result
}
