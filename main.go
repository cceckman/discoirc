// Binary discocli provides a console client for discoirc.
package main

import (
	"flag"
	"fmt"
	golog "log"
	"os"
	"time"

	"github.com/marcusolsson/discoirc/log"
	"github.com/marcusolsson/discoirc/model"
	"github.com/marcusolsson/discoirc/view"
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

	logpath = flag.String("logpath", "", "File to write debug logs to. Use a temporary file and directory if unset.")
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
	mchan := model.NewMockChannel(logger, "testnet", "#testing")
	client := model.DumbClient(map[string]model.Connection{
		"testnet": model.DumbConnection(map[string]model.Channel{
			"#testing": mchan,
		}),
	})
	model.EventGenerator(logger, mchan)

	session := view.NewConsoleSession(logger, client)
	session.SetTheme(Theme())
	go func() {
		time.Sleep(1 * time.Second)
		session.OpenChannel("testnet", "#testing")
	}()

	if err := session.Run(); err != nil {
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
	if *logpath == "" {
		logger, err := log.New()
		if err != nil {
			golog.Fatal("could not create log: ", err)
		}
		return logger
	}
	file, err := os.OpenFile(*logpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		golog.Fatal("could not open log: ", err)
	}
	result, err := log.NewFor(file)
	if err != nil {

		golog.Fatal("could not open log: ", err)
	}
	return result
}
