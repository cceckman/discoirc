package view

import (
	"github.com/cceckman/discoirc/ui/widgets"
	"github.com/marcusolsson/tui-go"
	"image"
)

func newSelector() *selector {
	r := &selector{}
	r.SetFill(' ')
	return r
}

type selector struct {
	widgets.Filler
}

func (f *selector) SizeHint() image.Point {
	return image.Pt(1, 0)
}

func (f *selector) SizePolicy() (tui.SizePolicy, tui.SizePolicy) {
	return tui.Maximum, tui.Preferred
}
