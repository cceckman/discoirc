// Binary discocli provides a console client for discoirc.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	golog "log"
	"os"

	"github.com/cceckman/discoirc/discocli/log"
	"github.com/cceckman/discoirc/discocli/model"
	"github.com/cceckman/discoirc/discocli/view"
	"github.com/jroimartin/gocui"
)

const (
	usage = `
%s:
	Connect to the discoirc daemon and display contents from it.

	By default (if no flags are provided), this will search for a discoircd connection in the usual
	search paths, and start one if none is found. It will then open the Meta view.

	Specifying a daemon address will force it to connect to that daemon instead; similarly,
	specifying a view will open that view instead.
`
)

var (
	help = flag.Bool("help", false, "Display a usage message.")

	lowcolor = flag.Bool("low-color", false, "Use 8-bit color rather than 256-bit.")
	logpath  = flag.String("log-path", "", "Path to write debug logs to. Use a temporary directory if unset.")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(1)
	}

	logger := LoggerOrDie()
	g := UiOrDie(logger)
	defer g.Close()

	// TODO: Populate the initial view from something else.
	// TODO: Implement Client properly.
	client := model.DumbClient(map[string]model.Connection{
		"testnet": model.DumbConnection(map[string]model.Channel{
			"#testing": model.NewMockChannel("#testing"),
		}),
	})

	ls := view.StartLayoutSwitcher(client, g, logger, &view.ChatViewInfo{
		Connection: "testnet",
		Channel:    "#testing",
	})
	defer ls.Done()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		logger.Print("encountered error, treating as nonfatal: ", err)
	}
}

func LoggerOrDie() *golog.Logger {
	// Initialize logger
	if *logpath == "" {
		tmp, err := ioutil.TempDir("", "")
		if err != nil {
			golog.Fatal("could not create directory for logging: ", err)
		}
		*logpath = tmp
	}
	logger, err := log.New(*logpath)
	if err != nil {
		golog.Fatal("could not open log file: ", err)
	}
	return logger
}

func UiOrDie(logger *golog.Logger) *gocui.Gui {
	// Initialize GUI
	colorMode := gocui.Output256
	if *lowcolor {
		colorMode = gocui.OutputNormal
	}

	g, err := gocui.NewGui(colorMode)
	if err != nil {
		logger.Fatal("could not initialize GUI: ", err)
	}
	return g
}
