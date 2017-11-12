// Binary discocli provides a console client for discoirc.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/model"
	"github.com/cceckman/discoirc/view"
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

	// use glog for standard logging as well as any explicitly-glogged stuff.
	glog.CopyStandardLogTo("INFO")

	// TODO: Populate the initial view from something else.
	// TODO: Implement Client properly.
	mchan := model.NewMockChannel("testnet", "#testing")
	client := model.DumbClient(map[string]model.Connection{
		"testnet": model.DumbConnection(map[string]model.Channel{
			"#testing": mchan,
		}),
	})
	model.EventGenerator(mchan)

	session := view.NewConsoleSession(client)
	session.SetTheme(Theme())
	go func() {
		time.Sleep(1 * time.Second)
		session.OpenChannel("testnet", "#testing")
	}()

	if err := session.Run(); err != nil {
		glog.Fatalf("unknown error: %v", err)
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
