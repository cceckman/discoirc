// 2017-12-24 cceckman <charles@cceckman.com>
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	gctl "github.com/cceckman/discoirc/ui"
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

	ui := tui.New(tui.NewHBox())
	ui.SetTheme(GetTheme())
	ctl := gctl.New(ctx, ui)

	go func() {
		time.Sleep(2 * time.Second)

		chanCtl := GetStubChannel(ctx, ctl, "Barnetic", "#discoirc", "discobot")
		toggle := &Toggle{
			Channel:  chanCtl,
			Duration: 2 * time.Second,
		}

		ctl.Update(func() {
			ui.SetKeybinding("Ctrl+Space", func() {
				glog.V(1).Info("toggling metadata cycling")
				toggle.Meta(ctx)
			})
			ui.SetKeybinding("Ctrl+A", func() {
				glog.V(1).Info("toggling message cycling")
				toggle.Messages(ctx)
			})
		})
		toggle.Messages(ctx)
		toggle.Meta(ctx)
	}()

	if err := ui.Run(); err != nil {
		panic(err)
	}
}
