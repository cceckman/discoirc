// Binary discocli provides a console client for discoirc.
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	golog "log"
	"os"

	"github.com/cceckman/discoirc/log"
	"github.com/cceckman/discoirc/model"
	"github.com/cceckman/discoirc/view/channel"
	"github.com/marcusolsson/tui-go"
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

	logdir = flag.String("log-dir", "", "Directory to write debug logs to. Use a temporary directory if unset.")
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

	// TODO: Populate the initial view from something else.
	// TODO: Implement Client properly.
	mchan := model.NewMockChannel(logger, "testnet", "#testing", "We're all mad here")
	client := model.DumbClient(map[string]model.Connection{
		"testnet": model.DumbConnection(map[string]model.Channel{
			"#testing": mchan,
		}),
	})

	model.MessageGenerator(logger, 99, mchan)

	ctx := context.Background()
	v := channel.NewView()
	ui := tui.New(v)
	ui.SetTheme(Theme())
	channel.New(ctx, v, ui, client, "testnet", "#testing")

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		logger.Fatalf("unknown error: %v", err)
	}
}

func Theme() *tui.Theme {
	t := tui.NewTheme()
	instance := tui.Style{
		Fg: tui.ColorWhite,
		Bg: tui.ColorBlack,
	}
	t.SetStyle("default", instance)
	instance.Reverse = true
	t.SetStyle("reverse", instance)
	return t
}

func LoggerOrDie() *golog.Logger {
	// Initialize logger
	if *logdir == "" {
		tmp, err := ioutil.TempDir("", "")
		if err != nil {
			golog.Fatal("could not create directory for logging: ", err)
		}
		*logdir = tmp
	}
	logger, err := log.New(*logdir)
	if err != nil {
		golog.Fatal("could not open log file: ", err)
	}
	return logger
}
