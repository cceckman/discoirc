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
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	Events bufchan.Broadcaster
	clients int32
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

	clientN := atomic.AddInt32(&s.clients, 1)

	log.Println("received Subscribe request", clientN, ":", req)
	defer log.Println("done with Subscribe request", clientN)

	// Must provide a filter, even a blank one.
	if req.Filter == nil {
		err := grpc.Errorf(
			codes.InvalidArgument,
			"SubscribeRequest must include a Filter",
		)
		log.Println("returning error: ", err)
		return err
	}

	// The sender's Context closes when the networked receiver closes out.
	// Set the listener's lifetime to the same.
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

	return nil
}

// Tick writes the message returned by 'msg' to the server.
func (s *Server) Tick(ctx context.Context, d time.Duration, t *stream.Id, msg func() string) {
	ticker := time.NewTicker(d)
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
	/*
	go func() {
		s := b.Send()
		<-ctx.Done()
		close(s)
	}()
	*/
	return &Server{
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
	log.Println("started listening at ", *socket)

	server := NewServer(ctx, time.Second*2)

	s := grpc.NewServer()
	stream.RegisterEventProviderServer(s, server)
	reflection.Register(s)

	log.Println("registered grpc server")

	// Start serving & listening in the background.
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to start serving: %v", err)
		}
	}()

	log.Println("Running...")

	log.Println("Starting tickers...")

	var hundred ticker = 100
	go server.Tick(ctx, time.Second * 2, &stream.Id{
		Plugin:  "discoirc",
		Network: "hundrednet",
		Channel: "hundreds",
	}, hundred.Count)

	var thousand ticker = 1000
	go server.Tick(ctx, time.Second * 3, &stream.Id{
		Plugin:  "discoirc",
		Network: "thousandnet",
		Channel: "thousands",
	}, (&thousand).Count)

	<-ctx.Done()
	// Closes down listener as well.
	s.GracefulStop()

}
