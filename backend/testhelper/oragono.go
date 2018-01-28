package testhelper

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sync"
)

// RealServer is an IRCd running as a process on localhost.
type RealServer struct {
	Address string

	cmd        *exec.Cmd
	stdout     bytes.Buffer
	stderr     bytes.Buffer
	collectors sync.WaitGroup

	cancel func()
}

// NewServer starts a new IRCd as a background process.
// It will be terminated when the context completes.
func NewServer(ctx context.Context) (*RealServer, error) {
	oragono, err := exec.LookPath("oragono")
	if err != nil {
		return nil, fmt.Errorf("could not find oragono executable: %v", err)
	}

	s := &RealServer{
		// TODO: include v4, v6, tls, non-tls addresses
		Address: config.Server.Listen[0],
	}

	ctx, s.cancel = context.WithCancel(ctx)

	s.cmd = exec.CommandContext(ctx, oragono, "run", "--conf", configPath)

	// Collect output for logs
	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	s.collectors.Add(1)
	go func() {
		s.stdout.ReadFrom(stdout)
		s.collectors.Done()
	}()
	go func() {
		s.stderr.ReadFrom(stderr)
		s.collectors.Done()
	}()

	if err := s.cmd.Start(); err != nil {
		return nil, err
	}
	return s, nil
}

// Output ensures the server is shut down, and returns its output.
func (s *RealServer) Output() (stdout, stderr []byte) {
	s.cancel()
	s.cmd.Wait()
	s.collectors.Wait()

	stdout = s.stdout.Bytes()
	stderr = s.stderr.Bytes()
	return
}
