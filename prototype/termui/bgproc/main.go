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
)

var (
	help = flag.Bool("help", false, "Display a usage message.")

	writeTo = flag.String("target", "", "file (pipe) to write lines to")
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

	if *writeTo == "" {
		log.Fatal("No file specified")
	}

	f, err := os.OpenFile(*writeTo, os.O_WRONLY | os.O_APPEND, 0) //os.ModeNamedPipe)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.Println("Opened for writing; starting counter")
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	p := primes.NewMemoizingPrimer()
	for {
		c := make(chan int)
		p.PrimesUpTo(1000000, c)
		for prime := range c {
			<-ticker.C
			fmt.Printf("Tick: %d\n", prime)
			fmt.Fprintf(f, "%07d\n", prime)
		}
	}
}
