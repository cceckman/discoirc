package widgets

import (
	"github.com/marcusolsson/tui-go"
)

var _ tui.Widget = &Filler{}

// NewFiller returns a new Filler widget
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

// SetFill sets the character used for filling.
func (f *Filler) SetFill(fill rune) {
	f.fill = fill
}

// Draw renders the Widget using the provided Painter
func (f *Filler) Draw(p *tui.Painter) {
	sz := f.Size()
	for j := 0; j < sz.Y; j++ {
		for i := 0; i < sz.X; i++ {
			p.DrawRune(i, j, f.fill)
		}
	}
}
