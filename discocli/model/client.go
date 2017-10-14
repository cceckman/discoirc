package model

// Client is the main interface to the DiscoIRC client.
// It hosts connections and their properties.
type Client interface {
	Connection(name string) Connection
	Connections() []string
}
