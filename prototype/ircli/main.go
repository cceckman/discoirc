// 2017-01-01 cceckman <charles@cceckman.com>
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	flog "github.com/cceckman/discoirc/prototype/log"
	"log"

	"github.com/cceckman/discoirc/prototype/ircli/client"
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

	var wg sync.WaitGroup

	// Start it up.
	c := client.NewClient()
	log.Println("Created client.")

	ctx := signalContext()

	wg.Add(1)
	go func() {
		for msg := range c.Listen(ctx) {
			log.Printf(">%v", msg)
		}
		log.Println("Listen channel done")
		wg.Done()
	}()

	log.Println("Starting connection.")
	if errs := c.Connect(); len(errs) > 0 {
		for _, err := range errs {
			log.Printf("ERR: %v", err)
		}
	}

	log.Println("Connection complete. Connected networks: ")
	for _, net := range c.ConnectedNetworks() {
		log.Println("\t", net)
	}

	wg.Add(1)
	go func() {
		<-ctx.Done()
		if errs := c.Disconnect(); len(errs) > 0 {
			for _, err := range errs {
				log.Printf("ERR: %v", err)
			}
		}
		wg.Done()
	}()

	wg.Wait()
}
