// 2017-01-01 cceckman <charles@cceckman.com>
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	flog "github.com/cceckman/discoirc/prototype/log"
	"log"

	"github.com/cceckman/discoirc/prototype/ircli/client"
)

var (
	help = flag.Bool("help", false, "Display a usage message.")
)

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

	// Start it up.
	c := client.NewClient()
	log.Println("Created client.")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for msg := range c.Listen(ctx) {
			log.Printf(">%v", msg)
		}
		log.Println("Listen channel done")
	}()

	log.Println("Starting connection.")
	if errs := c.Connect(); len(errs) > 0 {
		for _, err := range errs {
			log.Printf("ERR: %v", err)
		}
		os.Exit(1)
	}

	log.Println("Connection complete. Connected networks: ")
	for _, net := range c.ConnectedNetworks() {
		log.Println("\t", net)
	}

	n := 5
	log.Printf("Waiting %d seconds", n)
	time.Sleep(time.Second * time.Duration(n))
	log.Println("Done waiting! Disconnecting.")
	if errs := c.Disconnect(); len(errs) > 0 {
		for _, err := range errs {
			log.Printf("ERR: %v", err)
		}
		os.Exit(1)
	}

	time.Sleep(time.Second * 10)
}
