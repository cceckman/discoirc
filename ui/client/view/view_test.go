package view_test

import (
	"testing"

	"github.com/cceckman/discoirc/ui/client/view"
	"github.com/marcusolsson/tui-go"
)

func TestNetwork_NoContents(t *testing.T) {
	w := view.NewNetwork("Barnetic")

	surface := tui.NewTestSurface(20, 5)
	theme := tui.NewTheme()
	p := tui.NewPainter(surface, theme)
	p.Repaint(w)

	wantContents := `
Barnetic:           
                    
                    
                    
                    
`
	gotContents := surface.String()
	if gotContents != wantContents {
		t.Errorf("unexpected contents: got = \n%s\nwant = \n%s", gotContents, wantContents)
	}
}

func TestNetwork_NoChannels(t *testing.T) {
	w := view.NewNetwork("Barnetic")
	w.SetNick("discobot")
	w.SetConnection("…")

	surface := tui.NewTestSurface(25, 5)
	theme := tui.NewTheme()
	p := tui.NewPainter(surface, theme)
	p.Repaint(w)

	wantContents := `
Barnetic: …      discobot
                         
                         
                         
                         
`
	gotContents := surface.String()
	if gotContents != wantContents {
		t.Errorf("unexpected contents: got = \n%s\nwant = \n%s", gotContents, wantContents)
	}
}
