package widgets

import (
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

func NewSplash() tui.Widget {
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
