// Package plugins defines types and functions for managing and using local gRPC servers.
package plugins

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"path"
	"syscall"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	pluginKillDelay = flag.Int("plugin-kill-timeout-ms", 3000, "Timeout (in milliseconds) to wait for plugins to gracefully stop (SIGTERM) before ungracefully stopping them (SIGKILL).")
)

type Plugin interface {
	Connect() (*grpc.ClientConn, error)
	Name() string
}

type plugin struct {
	name string
	sock string
	done chan struct{}
}

func (p *plugin) Name() string {
	return p.name
}

func (p *plugin) Connect() (*grpc.ClientConn, error) {
	dialer := func(addr string, timeout time.Duration) (net.Conn, error) {
		return net.DialTimeout("unix", addr, timeout)
	}

	conn, err := grpc.Dial(
		p.sock,
		grpc.WithDialer(dialer),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("[%s] error in connecting to %s: %v", p.name, p.sock, err)
	}

	go func() {
		<-p.done
		conn.Close()
	}()

	return conn, nil
}

func Run(ctx context.Context, name, bin string) (Plugin, error) {
	// pluginDone is closed when the plugin process has exited.
	p := &plugin{
		name: name,
		done: make(chan struct{}),
	}

	// Create a new temp dir for this plugin to use.
	if sockdir, err := ioutil.TempDir("", "discoirc"); err != nil {
		return nil, fmt.Errorf("[%s] could not create socket dir: %v", p.name, err)
	} else {
		p.sock = path.Join(sockdir, "socket")
	}

	// Create Command object.
	cmd := exec.Command(bin, "-socket", p.sock)
	if err := logOutput(cmd); err != nil {
		return nil, err
	}

	// Handle the *context*'s exiting; attempt to terminate the plugin.
	go func() {
		<-ctx.Done()
		// Send interrupt signal...
		cmd.Process.Signal(syscall.SIGINT)

		timeout := time.Millisecond * time.Duration(*pluginKillDelay)
		select {
		case <-p.done:
			log.Printf("[%s] appears to have exited safely.", p.name)
			return
		case <-time.After(timeout):
			log.Printf("[%s] hasn't exited before %v, killing", p.name, timeout)
			cmd.Process.Kill()
		}
	}()

	// Start the Command, in the background.
	if err := cmd.Start(); err != nil {
		close(p.done)
		return nil, err
	}

	// Handle the plugin's exiting.
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Printf("[%s] error in waiting for plugin: %v", p.name, err)
		}
		close(p.done)
	}()


	return p, nil
}

func logOutput(cmd *exec.Cmd) error {
	// Log output.
	if stdout, err := cmd.StdoutPipe(); err != nil {
		return err
	} else {
		go func() {
			scan := bufio.NewScanner(stdout)
			for scan.Scan() {
				log.Printf("[%s] stdout: %s", cmd.Path, scan.Text())
			}
		}()
	}
	if stderr, err := cmd.StderrPipe(); err != nil {
		return err
	} else {
		go func() {
			scan := bufio.NewScanner(stderr)
			for scan.Scan() {
				log.Printf("[%s] stderr: %s", cmd.Path, scan.Text())
			}
		}()
	}

	return nil
}
