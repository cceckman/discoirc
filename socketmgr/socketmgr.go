// Package socketmgr provides an interface for finding & attaching to client/server sockets.
package socketmgr

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
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

const (
	sockperms = 0750 | os.ModeSocket
)

// T manages sockets between a discoirc client and server.
type T interface {
	// Get lists the available sockets, based on the environment.
	Get() []string
	// Listen finds or makes a socket for listening. It returns a listener and a path, or an error.
	Listen() (net.Listener, string, error)
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

// SetupDirs ensures that the given directories are ready for reading/writing.
// It uses the env-vars expansion of each directory.
func SetupDirs(dirs []string) {
	home := os.ExpandEnv("$HOME")
	for _, d := range dirs {
		p := os.ExpandEnv(d)
		if p == "" {
			continue
		}
		var mode os.FileMode
		if strings.HasPrefix(p, home) {
			mode = 0750
		} else {
			mode = 0777
		}
		err := os.MkdirAll(p, mode)
		if err != nil {
			log.Printf("could not set up socket directory %s -> %s: %v", d, p, err)
		}
	}
}

type sm struct {
	target string
	paths  []string
}

// resolvePaths gives, in order, the list of paths to search / create sockets in.
// They may be directories or exact paths.
func (s *sm) resolvePaths() []string {
	var preExp []string

	if s.target != "" {
		preExp = []string{s.target}
	} else {
		preExp = s.paths
	}

	r := make([]string, 0, len(s.paths))
	for _, p := range preExp {
		x := os.ExpandEnv(p)
		if p != "" {
			r = append(r, x)
		}
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
		if validateForRead(info) {
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
				if validateForRead(file) {
					r = append(r, path.Join(p, file.Name()))
				}
			}
		}
	}
	return r
}

// Listen creates or attaches to a socket along the search path.
func (s *sm) Listen() (net.Listener, string, error) {
	var lasterr error

	for _, p := range s.resolvePaths() {
		// Check if the particular file already exists.
		info, err := os.Stat(p)
		// If it already exists and is a directory, see if any existing sockets are open.
		if err == nil && info.IsDir() {
			lis, path, err := validateDirectoryListen(p)
			if err == nil {
				return lis, path, err
			} else {
				log.Printf("could not use path %s for sockets: %v", p, err)
			}
		}
		// Otherwise, treat as a specific path.
		lis, err := validateFileListen(p)
		if err == nil {
			return lis, p, err
		} else {
			log.Printf("could not use path %s for sockets: %v", p, err)
		}
	}

	return nil, "", fmt.Errorf("could not find a valid socket to listen on!", lasterr)
}

// validateForRead determines whether the path appears to be an available socket.
func validateForRead(f os.FileInfo) bool {
	ok := true
	ok = ok && (f.Mode()&sockperms == sockperms)
	// TODO: check ownership
	return ok
}

// validateFileListen attempts to create a socket on the file name.
// A file is valid if it is, or can be made, a socket.
// If it is valid, this returns a Listener bound to the socket.
func validateFileListen(p string) (net.Listener, error) {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		// Make it.
		f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL, sockperms)
		if err != nil {
			return nil, err
		}
		if err := f.Close(); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// Set permissions.
	if err := os.Chmod(p, sockperms); err != nil {
		return nil, err
	}
	// And attempt to bind.
	return net.ListenUnix("unix", &net.UnixAddr{"unix", p})
}

// validateDirectoryListen looks in a given directory for available sockets.
func validateDirectoryListen(p string) (net.Listener, string, error) {
	files, err := ioutil.ReadDir(p)
	if err != nil {
		// Can't try to open files, won't try to create them.
		return nil, p, fmt.Errorf("couldn't list files in %s: %v", p, err)
	}

	for _, file := range files {
		fpath := path.Join(p, file.Name())
		lis, err := validateFileListen(fpath)
		if err != nil {
			log.Printf("couldn't listen on %s: %v", fpath, err)
		} else {
			return lis, fpath, err
		}
	}

	// No open sockets in directory. create a new file.
	f, err := ioutil.TempFile(p, "discoirc-session-")
	if err != nil {
		return nil, "", err
	}
	if err := f.Close(); err != nil {
		return nil, "", err
	}
	fpath := path.Join(p, f.Name())
	lis, err := validateFileListen(fpath)
	return lis, fpath, err
}
