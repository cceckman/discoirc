// 2017-12-24 cceckman <charles@cceckman.com>
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/marcusolsson/tui-go"

	gctl "github.com/cceckman/discoirc/ui"
	"github.com/cceckman/discoirc/ui/widgets"
	"github.com/cceckman/discoirc/backend/demo"
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

	ui := tui.New(tui.NewHBox())
	ui.SetTheme(GetTheme())
	// TODO: maybe put this in controller?
	ui.SetWidget(widgets.NewSplash(ui))

	backend := demo.New()

	ctl := gctl.New(ui, backend)


	go func() {
		time.Sleep(2 * time.Second)

		toggle := &Toggle{
			Demo: backend,
			Net: "Barnetic",
			Chan: "#discoirc",
			Duration: 2 * time.Second,
		}

		ctl.Update(func() {
			ui.SetKeybinding("Ctrl+R", func() {
				glog.V(1).Info("toggling network cycling")
				toggle.Network()
			})
			ui.SetKeybinding("Ctrl+F", func() {
				glog.V(1).Info("toggling channel cycling")
				toggle.Channel()
			})
		})
		toggle.Network()
		toggle.Channel()
	}()

	if err := ui.Run(); err != nil {
		panic(err)
	}
}
