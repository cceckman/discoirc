// 2017-01-07 cceckman <charles@cceckman.com>
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cceckman/discoirc/prototype/localrpc/plugins"
	service "github.com/cceckman/discoirc/prototype/localrpc/proto"
	flog "github.com/cceckman/discoirc/prototype/log"
	"github.com/cceckman/discoirc/prototype/sigcontext"
	"golang.org/x/net/context"
)

var (
	help       = flag.Bool("help", false, "Display a usage message.")
	pluginList = flag.String("plugins", "plugin/main", "Comma-separated list of plugins to run.")
)

func TickAgainst(ctx context.Context, p plugins.Plugin) {
	conn, err := p.Connect()
	if err != nil {
		log.Printf("error in creating connection to %s: %v", p.Name(), err)
		return
	}

	cli := service.NewSimpleServiceClient(conn)
	tick := time.NewTicker(time.Second * 2)
	defer tick.Stop()

	for i := 1; true; i++ {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			req := &service.MyRequest{
				Event: &service.Event{
					Seq:  int64(i),
					Name: p.Name(),
					Msg:  fmt.Sprintf("tick %d for %s", i, p.Name()),
				},
			}
			log.Printf("[%s] sending request: %v", p.Name(), req)
			if res, err := cli.Do(ctx, req); err != nil {
				log.Printf("[%s] error for request [%v] : %v", p.Name(), req, err)
			} else {
				log.Printf("[%s] got response: %v", p.Name(), res)
			}
		}
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
	if err := flog.Init(); err != nil {
		// Reset before writing any more messages.
		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
	flog.LogArgs()
	// Above is boilerplate.

	ctx, cancel := context.WithCancel(sigcontext.New())

	plist := strings.Split(*pluginList, ",")
	for i, s := range plist {
		p, err := plugins.Run(ctx, fmt.Sprintf("%d", i), s)
		if err != nil {
			log.Printf("error starting plugin %d: %v", i, err)
			cancel()
			os.Exit(1)
		}
		go TickAgainst(ctx, p)
	}

	<-ctx.Done()
}
