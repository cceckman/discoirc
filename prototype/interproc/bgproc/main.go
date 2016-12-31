// 2016-12-29 cceckman <charles@cceckman.com>
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	flog "github.com/cceckman/discoirc/prototype/log"
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
	defer l.Close()

	ctx := signalContext()
	p := primes.NewMemoizingPrimer()

	// Listen for new connections in a background thread;
	// close the listener in main().
	go func() {
		for cid := 0; true; cid++ {

			log.Println("Awaiting connection...")
			newCon, err := l.AcceptUnix()
			if err != nil {
				log.Println(err)
			}
			log.Printf("Got a new connection (%d)\n", cid)

			// Start a background writer.
			go func(ctx context.Context, n int, conn net.Conn) {
				defer conn.Close()

				// Explicitly respond to channel getting closed.
				// May be a no-op, if the write fails first.
				newCtx, cancel := context.WithCancel(ctx)
				go func() {
					defer cancel()

					b := []byte{}
					err := conn.SetReadDeadline(time.Time{}) // set no deadline
					if err != nil {
						log.Println(err)
					}

					log.Printf("watching connection %d for close\n", n)
					// Block until "read" gets "EOF".
					for {
						if _, err := conn.Read(b); err != nil {
							log.Printf("connection %d is closed: %v.\n", n, err)
							// Stream closed.
							return
						}
					}
				}()

				// Rate-limit writes, one per second, just for fun.
				ticker := time.NewTicker(time.Second)
				defer ticker.Stop()

				c := make(chan int)
				p.PrimesUpTo(1000000, c)

				for prime := range c {
					select {
					case <-newCtx.Done():
						return
					case <-ticker.C:
						fmt.Printf("Tick on %d: %d\n", n, prime)
						conn.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))

						_, err := fmt.Fprintf(conn, "%07d\n", prime)
						if err != nil {
							switch e := err.(type) {
							case net.Error:
								if e.Temporary(){
									continue
								} else {
									log.Printf("permanent write error for channel %d is not temporary", n)
								}
							}
							// "default" case
							log.Printf("write error on channel %d: %v\n", n, err)
							return
						}
					}
				}
			}(ctx, cid, newCon)
		}
	}()

	// Wait for the context to be done before exiting.
	<-ctx.Done()
}
