// Package main provides an EventProvider client that listens for all Events that the EventProvider server at -socket sends.
// TODO: Uses 'prototype/sigcontext'
// TODO: Doesn't exercise the backfill option.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/cceckman/discoirc/api/stream"
	"github.com/cceckman/discoirc/prototype/sigcontext"
)

var (
	help    = flag.Bool("help", false, "Display a usage message.")
	socket  = flag.String("socket", "", "Socket address to read events from.")
	timeout = flag.Int("timeout", 3000, "(ms) Timeout when attempting to connect to the server at --socket.")
)

// Dial is a dialer function for GRPC calls.
func Dial(addr string, to time.Duration) (net.Conn, error) {
	return net.DialTimeout("unix", addr, to)
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
	// Above is boilerplate.

	ctx := sigcontext.New()

	if *socket == "" {
		log.Fatal("no socket specified")
	}

	conn, err := grpc.Dial(
		*socket,
		grpc.WithDialer(Dial),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal("could not connect to EventProvider: ", err)
	}
	log.Println("connected to EventProvider")

	client := stream.NewEventProviderClient(conn)
	sub, err := client.Subscribe(ctx, &stream.SubscribeRequest{
		Filter: &stream.Filter{
			Matches: []*stream.Match{&stream.Match{}},
		},
	})
	if err != nil {
		log.Fatal("error in subscribe request: ", err)
	}
	log.Println("established Subscribe channel")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; true; i++ {
			log.Println("waiting for subscription...")
			resp, err := sub.Recv()
			log.Print("received event ", i, " ", resp)
			if err != nil {
				log.Print("error in receive stream: ", err)
				return
			}
		}
		wg.Done()
	}()

	<-ctx.Done()
	log.Println("Shutting down...")
	wg.Wait()
}
