// 2016-12-29 cceckman <charles@cceckman.com>
package main

import (
	"flag"
	"fmt"
	"os"

	flog "github.com/cceckman/discoirc/prototype/termui/log"
	"log"

	"github.com/cceckman/primes"

	"time"
	"net"
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


	// Create a Unix domain socket.
	network := "unix"
	addr := "/tmp/discod"
	l, err := net.ListenUnix(network, &net.UnixAddr{addr, network})
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(addr)

	p := primes.NewMemoizingPrimer()
	for cid := 0; true; cid++{
		log.Println("Awaiting connection...")
		newCon, err := l.AcceptUnix()
		if err != nil {
			log.Println(err)
		}
		log.Printf("Got a new connection (%d)\n", cid)

		// Start a background writer.
		go func(n int, conn net.Conn) {
			// TODO: needs to reliably handle closed other end.
			defer conn.Close()

			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			c := make(chan int)
			p.PrimesUpTo(1000000, c)

			for prime := range c {
				<-ticker.C
				fmt.Printf("Tick on %d: %d\n", n, prime)
				fmt.Fprintf(conn, "%07d\n", prime)
			}
		}(cid, newCon)
	}
}
