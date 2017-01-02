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
	go func() {
		for msg := range c.Listen(ctx) {
			log.Printf(">%s", msg)
		}
	}()

	log.Println("Starting connection.")
	for err := range c.Connect() {
		log.Printf("ERR: %v", err)
	}
	log.Println("Connection complete.")

	time.Sleep(time.Second * 10)
	log.Println("Done waiting! Cancelling context.")
	cancel()

	time.Sleep(time.Second * 10)
}
