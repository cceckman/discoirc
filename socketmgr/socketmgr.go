// Package socketmgr provides an interface for finding & attaching to client/server sockets.
package socketmgr

import (
	"flag"
	"io/ioutil"
	"log"
	_ "net"
	"os"
	"path"
	"strings"
)

var (
	defaultPaths = []string{
		"$SOCKDIR",
		"$HOME/.discoirc/sockets",
		"/tmp/discoirc/sockets",
	}
	target = flag.String(
		"socket", "", `
	Unix domain socket to find a server on.
	Used if explicitly specified; otherwise, searches, in order:
	`+strings.Join(defaultPaths, "\n  "),
	)
)

// T manages sockets between a discoirc client and server.
type T interface {
	// Get lists the available sockets, based on the environment.
	Get() []string
	// Make creates a new socket, and returns an open listener for it.
	// Make() (net.Listener, string, error)
}

func New() T {
	if !flag.Parsed() {
		flag.Parse()
	}
	return &sm{
		target: *target,
		paths:  defaultPaths,
	}
}

type sm struct {
	target string
	paths  []string
}

// resolvePaths gives, in order, the list of paths to search / create sockets in.
// They may be directories or exact paths.
func (s *sm) resolvePaths() []string {
	var r []string

	if s.target != "" {
		r = []string{s.target}
	} else {
		r = s.paths
	}

	for i, p := range r {
		r[i] = os.ExpandEnv(p)
	}

	return r
}

// Get gets zero or more sockets which appear to be available.
func (s *sm) Get() []string {
	r := []string{}

	for _, p := range s.resolvePaths() {
		info, err := os.Stat(p)
		if err != nil {
			log.Printf("couldn't read %s: %v", p, err)
			continue
		}
		if validatePath(info) {
			r = append(r, p)
		}

		if info.IsDir() {
			// List files within the path.
			files, err := ioutil.ReadDir(p)
			if err != nil {
				log.Printf("couldn't list files in %s: %v", p, err)
				continue
			}
			for _, file := range files {
				if validatePath(file) {
					r = append(r, path.Join(p, file.Name()))
				}
			}
		}
	}
	return r
}

// validatePath determines whether the path appears to be an available socket.
func validatePath(f os.FileInfo) bool {
	ok := true
	ok = ok && (f.Mode()&os.ModeSocket != 0)
	// TODO: check ownership
	return ok
}
