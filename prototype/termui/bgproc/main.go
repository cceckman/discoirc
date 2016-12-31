// 2016-12-29 cceckman <charles@cceckman.com>
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	flog "github.com/cceckman/discoirc/prototype/termui/log"
	"log"

	"github.com/cceckman/primes"

	"net"
	"time"
)

var (
	help = flag.Bool("help", false, "Display a usage message.")
)

// signalContext returns a Context that ends with SIGINT or SIGTERM.
func signalContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Printf("Received signal %s, starting shutdown", sig.String())
		cancel()
	}()

	return ctx
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

	// Create a Unix domain socket.
	network := "unix"
	addr := "/tmp/discod"
	l, err := net.ListenUnix(network, &net.UnixAddr{
		Name: addr,
		Net:  network,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Printf("Removing %s\n", addr)
		err := os.Remove(addr)
		if err != nil {
			log.Printf("couldn't remove %s: %v\n", addr, err)
		}
	}()

	ctx := signalContext()

	p := primes.NewMemoizingPrimer()
	for cid := 0; true; cid++ {
		select {
		case <-ctx.Done():
			break
		default:
			// pass
		}

		log.Println("Awaiting connection...")
		newCon, err := l.AcceptUnix()
		if err != nil {
			log.Println(err)
		}
		log.Printf("Got a new connection (%d)\n", cid)

		// Start a background writer.
		go func(ctx context.Context, n int, conn net.Conn) {
			// TODO: needs to reliably handle closed other end.
			defer conn.Close()

			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			c := make(chan int)
			p.PrimesUpTo(1000000, c)

			for prime := range c {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					// pass, see below.
				}
				fmt.Printf("Tick on %d: %d\n", n, prime)
				fmt.Fprintf(conn, "%07d\n", prime)
			}
		}(ctx, cid, newCon)
	}
}
