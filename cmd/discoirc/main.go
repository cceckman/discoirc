// 2017-12-24 cceckman <charles@cceckman.com>
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/cceckman/discoirc/ui/widgets"
	"github.com/golang/glog"
	"github.com/marcusolsson/tui-go"
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
	defer glog.Flush()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ui := tui.New(widgets.NewSplash())
	ui.SetTheme(GetTheme())
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	go func() {
		ctl, newRoot := GetStubChannel(ctx, ui, "Barnetic", "discobot", "#discoirc")
		toggle := &Toggle{
			Channel:  ctl,
			Duration: 2 * time.Second,
		}

		time.Sleep(2 * time.Second)
		ui.Update(func() {
			ui.SetKeybinding("Ctrl+Space", func() {
				glog.V(1).Info("toggling metadata cycling")
				toggle.Meta(ctx)
			})
			ui.SetKeybinding("Ctrl+A", func() {
				glog.V(1).Info("toggling message cycling")
				toggle.Messages(ctx)
			})
			ui.SetWidget(newRoot)
			toggle.Messages(ctx)
			toggle.Messages(ctx)
		})
	}()

	if err := ui.Run(); err != nil {
		panic(err)
	}
}
