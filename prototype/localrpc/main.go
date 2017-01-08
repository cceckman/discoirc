// 2017-01-07 cceckman <charles@cceckman.com>
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	service "github.com/cceckman/discoirc/prototype/localrpc/proto"
	flog "github.com/cceckman/discoirc/prototype/log"
	"github.com/cceckman/discoirc/prototype/sigcontext"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	help    = flag.Bool("help", false, "Display a usage message.")
	plugins = flag.String("plugins", "plugin/main", "Comma-separated list of plugins to run.")

	pluginKillDelay = flag.Int("plugin-kill-timeout-ms", 3000, "Timeout (in milliseconds) to wait for plugins to gracefully stop (SIGTERM) before ungracefully stopping them (SIGKILL).")
)

func dialSock(addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("unix", addr, timeout)
}

func ticker(ctx context.Context, sock string) {
}

func runPlugin(ctx context.Context, wg *sync.WaitGroup, sockpath *os.File, bin string) {
	wg.Add(1)

	// TODO: Should probably re-run on failure.
	cmd := exec.CommandContext(ctx, bin, "-socket", sockpath.Name())
	log.Printf("plugin %s: running command: [%s] [%s]", bin, cmd.Path, strings.Join(cmd.Args, " "))

	// Capture output using background threads.
	go func() {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Printf("plugin %s: error in getting stdout: %v", bin, err)
		}
		r := bufio.NewScanner(stdout)
		for r.Scan() {
			log.Printf("plugin %s: stdout: %s\n", bin, r.Text())
		}
		log.Printf("plugin %s: stdout done", bin)
	}()
	go func() {
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Printf("plugin %s: error in getting stderr: %v", bin, err)
		}
		r := bufio.NewScanner(stderr)
		for r.Scan() {
			log.Printf("plugin %s: stderr: %s\n", bin, r.Text())
		}
		log.Printf("plugin %s: stderr done", bin)
	}()

	// Start program.
	if err := cmd.Start(); err != nil {
		log.Printf("plugin %s: error starting: %v", bin, err)
		return
	}

	// TERM, then KILL, the process when the context exits.
	go func() {
		// Capture when the program exits safely.
		s := make(chan bool, 1)
		go func() {
			if err := cmd.Wait(); err != nil {
				log.Printf("plugin %s: error waiting: %v", bin, err)
			}
			s <- true
		}()

		<-ctx.Done()
		log.Printf("plugin %s: sending SIGINT", bin)
		if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
			log.Printf("plugin %s: error sending SIGINT: %v", err)
		}

		delay := time.Duration(*pluginKillDelay) * time.Millisecond
		wait := time.After(delay)

		select {
		case <-wait:
			log.Printf("plugin %s: process didn't safely exit after %v, killing", bin, delay)
			if err := cmd.Process.Kill(); err != nil {
				log.Printf("plugin %s: error killing: %v", err)
			}
			return
		case <-s:
			// pass
		}
		log.Printf("plugin %s: exited? %t", cmd.ProcessState.Exited())
		wg.Done()
	}()

	// Start a thread to write to it.
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			// Close & clean socket when done.
			sockpath.Close()
			os.Remove(sockpath.Name())
		}()
		conn, err := grpc.DialContext(
			ctx,
			sockpath.Name(),
			grpc.WithDialer(dialSock),
			grpc.WithInsecure(),
		)
		if err != nil {
			log.Printf("error in connecting to %s: %v", sockpath.Name(), err)
			return
		}
		defer conn.Close()

		cli := service.NewSimpleServiceClient(conn)
		tick := time.NewTicker(time.Second * 2)
		defer tick.Stop()

		for i := 1; true; i++ {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				// Send a message, wait for a response.
				req := &service.MyRequest{
					Event: &service.Event{
						Seq:  int64(i),
						Name: "ping",
						Msg:  fmt.Sprintf("ping #%d!", i),
					},
				}
				if rsp, err := cli.Do(ctx, req); err != nil {
					log.Printf("error in request #%d on %s: %v", i, sockpath.Name(), err)
				} else {
					log.Printf("received response: %v", rsp)
				}
			}
		}
	}()
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s:	 \nUsage:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(1)
	}
	if err := flog.Init(); err != nil {
		// Reset before writing any more messages.
		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
	flog.LogArgs()
	// Above is boilerplate.

	ctx := sigcontext.New()

	plist := strings.Split(*plugins, ",")
	var wg sync.WaitGroup
	for _, s := range plist {
		sockfile, err := ioutil.TempFile("", "discoplugin-")
		if err != nil {
			log.Fatalf("could not make tempfile for socket: %v", err)
		}
		runPlugin(ctx, &wg, sockfile, s)
	}

	// Wait for all plugins to say they're done.
	wg.Wait()
}
