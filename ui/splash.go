package view

import (
	"context"
	"github.com/marcusolsson/tui-go"
)

const splash string = `
             ||
             ||
           <><><>
         <><><><><>
        <><><><><><>
        <><><><><><>
        <><><><><><>
         <><><><><>
           <><><>

          discoirc

github.com/cceckman/discoirc

`

type splashRequest int

func (_ splashRequest) String() string {
	return "splash screen view"
}
func (_ splashRequest) New(_ context.Context, _ *ConsoleSession) tui.Widget {
	return tui.NewHBox(
		tui.NewSpacer(),
		tui.NewVBox(
			tui.NewSpacer(),
			tui.NewLabel(splash),
			tui.NewSpacer(),
		),
		tui.NewSpacer(),
	)
}
