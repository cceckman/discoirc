package model

// Connection represents a relationship with an IRC network (or server).
// It may be connected to one or more channels.
type Connection interface {
	Channel(name string) Channel
	Channels() []string

	// TODO add server info
}

type DumbConnection map[string]Channel

func (f DumbConnection) Channel(name string) Channel {
	return f[name]
}

func (f DumbClient) Channels() []string {
	result := make([]string, len(f))
	i := 0
	for k := range f {
		result[i] = k
		i++
	}
	sort.Strings(result)
	return result
}
