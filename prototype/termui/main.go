// 2016-12-31 cceckman <charles@cceckman.com>
package main

import(
	"flag"
	"fmt"
	"os"

	flog "github.com/cceckman/discoirc/prototype/log"
	"log"

	"github.com/jroimartin/gocui"
)

var(
	help = flag.Bool("help", false, "Display a usage message.")
)

func main() {
	// Set up flags, respond to --help.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s:	 \nUsage:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(1)
	}
	// Initialize logging.
	if err := flog.Init(); err != nil {
		// Reset before writing any more messages.
		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
	flog.LogArgs()

	// Start GUI.
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	if err := SetupUI(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
