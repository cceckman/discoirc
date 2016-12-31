// 2016-12-29 cceckman <charles@cceckman.com>
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	flog "github.com/cceckman/discoirc/prototype/termui/log"
	"log"

	"github.com/jroimartin/gocui"
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

	// Start GUI.
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// Start window layout-er
	g.SetManagerFunc(LayoutPanes)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	// Initialize and start the reader of the other program.
	r := &RemoteReader{
		view: "hello",
		gui:  g,
	}
	_ = r
	go r.Start()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// RemoteReader populates the given view with values read from inputChannel.
type RemoteReader struct {
	view string
	gui  *gocui.Gui
}

func (m RemoteReader) listen() <-chan int {
	out := make(chan int)

	// Attempt reconnections, write results.
	go func() {
		timeout := 1
		for {
			network := "unix"
			addr := "/tmp/discod"
			d, err := net.DialUnix(network, nil, &net.UnixAddr{
				Name: addr,
				Net:  network,
			})
			if err != nil {
				log.Print(err)

				// exponential backoff in timeout
				timeout *= 2
				if timeout > 10 {
					timeout = 10
				}

				// Sleep before retrying.
				time.Sleep(time.Second * time.Duration(timeout))
				continue
			}
			// Success.
			timeout = 1

			for {
				i := 0
				_, err := fmt.Fscanf(d, "%07d\n", &i)
				if err != nil {
					// Try to reconnect.
					d.Close()
					break // the inner loop.
				} else {
					out <- i
				}
			}
		}
	}()

	return out
}

// update continuously reads from the Listen channel and writes events to the GUI
func (m RemoteReader) Start() {
	c := m.listen()

	for i := range c {
		n := i
		// Update to GUI must happen asynchronously; so, dispatch an event.
		m.gui.Execute(func(g *gocui.Gui) error {
			if v, err := g.View(m.view); err != nil {
				if err != gocui.ErrUnknownView {
					log.Println(err)
					return err
				}
			} else {
				v.Clear()
				// OK to update view.
				fmt.Fprintf(v, "%07d\n", n)
			}
			return nil
		})
	}
}

var (
	_ gocui.ManagerFunc = LayoutPanes
)

// LayoutPanes is a ManagerFunc that sizes the "hello" view.
func LayoutPanes(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if _, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
