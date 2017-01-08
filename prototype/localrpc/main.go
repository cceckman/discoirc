// 2017-01-07 cceckman <charles@cceckman.com>
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"
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
)

func dialSock(addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("unix", addr, timeout)
}

func ticker(ctx context.Context, sock string) {
	conn, err := grpc.DialContext(
		ctx,
		sock,
		grpc.WithDialer(dialSock),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Printf("error in connecting to %s: %v", sock, err)
		return
	}

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
				log.Printf("error in request #%d on %s: %v", i, sock, err)
			} else {
				log.Printf("received response: %v", rsp)
			}
		}
	}
}

func runPlugin(ctx context.Context, sockpath, bin string) {
	// TODO: Should probably re-run on failure.
	cmd := exec.CommandContext(ctx, bin, "-socket", sockpath)
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

	// Start output in background thread, let it run.
	go func() {
		if err := cmd.Start(); err != nil {
			log.Printf("plugin %s: error starting: %v", bin, err)
		}

		if err := cmd.Wait(); err != nil {
			log.Printf("plugin %s: error waiting: %v", bin, err)
		}
	}()

	// Start a thread to write to it.
	go ticker(ctx, sockpath)
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
	sockdir := os.TempDir()

	plist := strings.Split(*plugins, ",")
	for i, s := range plist {
		sockpath := path.Join(sockdir, fmt.Sprintf("plugin%d", i))
		runPlugin(ctx, sockpath, s)
	}

	<-ctx.Done()
}
