// 2016-12-29 cceckman <charles@cceckman.com>
package main

import (
	"flag"
	"fmt"
	"os"	
	"net"

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
		pane: "hello",
		gui: g,
	}
	r.Start()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// RemoteReader populates the given pane with values read from inputChannel.
type RemoteReader struct {
	pane string
	inputChannel chan int

	gui *gocui.Gui
}

func (m RemoteReader) Start() {
	m.listen()
	go m.update()
}

func (m RemoteReader) listen() {
	network := "unix"
	addr := "/tmp/discod"
	d, err := net.DialUnix(network, nil, &net.UnixAddr{addr, network})
	if err != nil {
		log.Fatal(err)
	}

	m.inputChannel = make(chan int)
	go func(){
		defer d.Close()
		defer close(m.inputChannel)
		for {
			i := 0
			_, err := fmt.Fscanf(d, "%07d\n", &i)
			if err != nil {
				return
			}
			m.inputChannel <- i
		}
	}()
}

// update continuously reads from inputChannel and writes events to the GUI
func (m RemoteReader) update() {
	for n := range m.inputChannel {
		n := n
		// Update to GUI must happen asynchronously; so, dispatch an event.
		update := func(g *gocui.Gui) error {
			if v, err := g.View(m.pane); err != nil {
				if err != gocui.ErrUnknownView {
					return err
				}
			} else {
				// OK to update view.
				fmt.Fprintf(v, "%07d\n", n)
			}
			return nil
		}
		m.gui.Execute(update)
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
