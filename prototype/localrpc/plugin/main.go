// A plugin; that is, a SimpleService server.
package main

import (
	// Generated code uses the legacy package.
	// "context"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"net"
	"os"

	service "github.com/cceckman/discoirc/prototype/localrpc/proto"
	"github.com/cceckman/discoirc/prototype/sigcontext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	help = flag.Bool("help", false, "Display a usage message.")
	sock = flag.String("socket", "", "UNIX domain socket to listen on")
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

	ctx := sigcontext.New()

	if *sock == "" {
		log.Fatal("no socket specified")
	}
	lis, err := net.Listen("unix", *sock)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	service.RegisterSimpleServiceServer(s, &echoServer{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("Running...")

	<-ctx.Done()
	// Closes down listener as well.
	s.GracefulStop()
}

type echoServer struct{}

var _ service.SimpleServiceServer = &echoServer{} // Say that three times fast.

func (_ *echoServer) Do(ctx context.Context, req *service.MyRequest) (*service.MyResponse, error) {
	return &service.MyResponse{Event: req.Event}, nil
}
