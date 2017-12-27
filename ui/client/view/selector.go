package view

import (
	"github.com/cceckman/discoirc/ui/widgets"
	"github.com/marcusolsson/tui-go"
	"image"
)

func newIndicator() *indicator {
	r := &indicator{}
	r.SetFill(' ')
	return r
}

type indicator struct {
	widgets.Filler
}

func (f *indicator) SizeHint() image.Point {
	return image.Pt(1, 0)
}

func (f *indicator) SizePolicy() (tui.SizePolicy, tui.SizePolicy) {
	return tui.Maximum, tui.Preferred
}
