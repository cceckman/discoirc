// Package main provides a fake EventProvider that listens on -socket for clients, and sends regular
// events to a few channels.
// TODO: Uses 'prototype/bufchan.Broadcaster'
// TODO: Uses 'prototype/sigcontext'
// TODO: Doesn't support the backfill option.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/cceckman/discoirc/api/stream"
	"github.com/cceckman/discoirc/prototype/bufchan"
	"github.com/cceckman/discoirc/prototype/sigcontext"
)

var (
	help   = flag.Bool("help", false, "Display a usage message.")
	socket = flag.String("socket", "", "Socket address to write events to.")
)

type ticker int

func (t *ticker) Count() string {
	s := strconv.Itoa(int(*t))
	(*t) = (*t) + 1
	return s
}

type Server struct {
	D      time.Duration
	Events bufchan.Broadcaster
}

var _ stream.EventProviderServer = &Server{}

// Subscribe is the handler for EventProvider.Subscribe.
// TODO: it does not support the backfill option.
func (s *Server) Subscribe(
	req *stream.SubscribeRequest,
	sender stream.EventProvider_SubscribeServer,
) error {
	if s.Events == nil {
		return fmt.Errorf("event Broadcaster not initialized")
	}
	go func() {
		// The sender's context closes when the networked receiver closes out;
		// hence, the listener closes.
		listener := s.Events.Listen(sender.Context())
		for i := range listener {
			e := i.(*stream.Event)
			if req.Filter.Exec(e) {
				resp := &stream.SubscribeResponse{
					Event: e,
				}
				err := sender.Send(resp)
				if err != nil {
					log.Printf("error sending to listener stream: %v", err)
				}
			}

		}
	}()
	return nil
}

// Tick writes the message returned by 'msg' to the server.
func (s *Server) Tick(ctx context.Context, t *stream.Id, msg func() string) {
	ticker := time.NewTicker(s.D)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e := &stream.Event{
				Stream: t,
				Text:   msg(),
			}
			select {
			case <-ctx.Done():
				return
			case s.Events.Send() <- e:
				log.Println("sent ", e)
			}
		}
	}
}

func NewServer(ctx context.Context, d time.Duration) *Server {
	b := bufchan.NewBroadcaster()
	go func() {
		s := b.Send()
		<-ctx.Done()
		close(s)
	}()
	return &Server{
		D:      d,
		Events: b,
	}
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
	lis, err := net.Listen("unix", *socket)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := NewServer(ctx, time.Second*2)

	s := grpc.NewServer()
	stream.RegisterEventProviderServer(s, server)
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to start serving: %v", err)
	}

	log.Println("Running...")
	log.Println("Starting tickers...")

	var hundred ticker = 100
	server.Tick(ctx, &stream.Id{
		Plugin:  "discoirc",
		Network: "hundrednet",
		Channel: "hundreds",
	}, hundred.Count)

	var thousand ticker = 1000
	server.Tick(ctx, &stream.Id{
		Plugin:  "discoirc",
		Network: "thousandnet",
		Channel: "thousands",
	}, (&thousand).Count)

	<-ctx.Done()
	// Closes down listener as well.
	s.GracefulStop()

}
