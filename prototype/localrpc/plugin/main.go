// A plugin; that is, a SimpleService server.
package main

import (
	"strconv"
	// Generated code uses the legacy package.
	// "context"
	"golang.org/x/net/context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	service "github.com/cceckman/discoirc/prototype/localrpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	help = flag.Bool("help", false, "Display a usage message.")
	port = flag.Int("port", 8001, "Port to listen on")
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
	// Above is boilerplate.

	// TODO: Listen on a local socket instead.
	lis, err := net.Listen("tcp", strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	service.RegisterSimpleServiceServer(s, &echoServer{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type echoServer struct{}

var _ service.SimpleServiceServer = &echoServer{} // Say that three times fast.

func (_ *echoServer) Do(ctx context.Context, req *service.MyRequest) (*service.MyResponse, error) {
	return &service.MyResponse{Event: req.Event}, nil
}
