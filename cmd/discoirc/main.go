// discoirc is a terminal IRC client.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/backend/demo"
	gctl "github.com/cceckman/discoirc/ui"
	"github.com/cceckman/discoirc/ui/widgets"
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

	ui, err := tui.New(tui.NewHBox())
	if err != nil {
		glog.Fatal("error intitializing UI: ", err)
	}

	ui.SetTheme(getTheme())
	// TODO: maybe put this in controller?
	ui.SetWidget(widgets.NewSplash(ui))

	backend := demo.New()

	ctl := gctl.New(ui, backend)

	go func() {
		time.Sleep(2 * time.Second)
		ctl.Update(func() {
			ctl.ActivateClient()
		})

		toggle := &Toggle{
			Demo:     backend,
			Net:      "Barnetic",
			Chan:     "#discoirc",
			Duration: 2 * time.Second,
		}

		ctl.Update(func() {
			ui.SetKeybinding("Ctrl+R", func() {
				glog.V(1).Info("toggling network cycling")
				toggle.network()
			})
			ui.SetKeybinding("Ctrl+F", func() {
				glog.V(1).Info("toggling channel cycling")
				toggle.channel()
			})
			ui.SetKeybinding("Ctrl+V", func() {
				glog.V(1).Info("toggling message cycling")
				toggle.messages()
			})
			toggle.network()
			toggle.channel()
			toggle.messages()
		})
	}()

	if err := ui.Run(); err != nil {
		panic(err)
	}
}
