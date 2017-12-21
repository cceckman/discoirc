package widgets

import (
	"image"
	"github.com/marcusolsson/tui-go"
)

// TailBox is a container Widget that may not show all its Widgets.
// While tui.Box attempts to show every contained Widget - sometimes shrinking
// those Widgets to do so- TailBox prioritizes completely displaying its last
// Widget, then the next-to-last widget, etc.
// It is vertically-aligned, i.e. all the contained Widgets have the same width.
type TailBox struct {
	tui.WidgetBase
	sz image.Point
	contents []tui.Widget
}

var _ tui.Widget = &TailBox{}

func NewTailBox(w ...tui.Widget) *TailBox {
	return &TailBox{
		contents: w,
	}
}

func (t *TailBox) Append(w tui.Widget) {
	t.contents = append(t.contents, w)
	t.doLayout(t.Size())
}

func (t *TailBox) SetContents(w ...tui.Widget) {
	t.contents = w
	t.doLayout(t.Size())
}

func (t *TailBox) Draw(p *tui.Painter) {
	p.WithMask(image.Rect(0, 0, t.sz.X, t.sz.Y), func(p *tui.Painter) {
		// Draw background
		p.FillRect(0, 0, t.sz.X, t.sz.Y)

		// Draw from the bottom up.
		space := t.sz.Y
		p.Translate(0, space)
		defer p.Restore()
		for i := len(t.contents) - 1; i >= 0 && space > 0; i-- {
			w := t.contents[i]
			space -= w.Size().Y
			p.Translate(0, -w.Size().Y)
			defer p.Restore()
			w.Draw(p)
		}
	})
}

// Resize recalculates the layout of the box's contents.
func (t *TailBox) Resize(size image.Point) {
	t.WidgetBase.Resize(size)
	defer func() {
		t.sz = size
	}()

	// If it's just a height change, Draw should do the right thing already.
	if size.X != t.sz.X {
		t.doLayout(size)
	}
}

func (t *TailBox) doLayout(size image.Point) {
	for _, w := range t.contents {
		hint := w.SizeHint()
		// Set the width to the container width, and height to the requested height
		w.Resize(image.Pt(size.X, hint.Y))
		// ...and then resize again, now that the Y-hint has been refreshed by the X-value.
		hint = w.SizeHint()
		w.Resize(image.Pt(size.X, hint.Y))
	}
}
