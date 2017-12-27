package widgets

import (
	"github.com/marcusolsson/tui-go"
)

var _ tui.Widget = &Filler{}

func NewFiller(fill rune) *Filler {
	return &Filler{
		fill: fill,
	}
}

// Filler is a Widget that fills itself with a given rune.
// It expands according to the usual SizePolicy mapping;
// note that "Maximum" may be necessary to restrict it in
// a single dimension.
type Filler struct {
	tui.WidgetBase

	fill rune
}

func (f *Filler) Draw(p *tui.Painter) {
	sz := f.Size()
	for j := 0; j < sz.Y; j++ {
		for i := 0; i < sz.X; i++ {
			p.DrawRune(i, j, f.fill)
		}
	}
}
