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

var clientTests = []struct {
	test  string
	setup func() tui.Widget
	want  string
}{
	{
		test: "blank client",
		setup: func() tui.Widget {
			w := view.New()
			return w
		},
		want: `
                         
                         
                         
                         
                         
`,
	},
	{
		test: "ordered networks",
		setup: func() tui.Widget {
			w := view.New()
			w.GetNetwork("Barnet").SetNick("barnacle")
			w.GetNetwork("AlphaNet").SetNick("edward")
			return w
		},
		want: `
AlphaNet:          edward
Barnet:          barnacle
                         
                         
                         
`,
	},
	{
		test: "Removed last network",
		setup: func() tui.Widget {
			w := view.New()
			w.GetNetwork("Charlienet").SetNick("charles")
			w.GetNetwork("Barnet").SetNick("barnacle")
			w.GetNetwork("AlphaNet").SetNick("edward")
			w.RemoveNetwork("Charlienet")
			return w
		},
		want: `
AlphaNet:          edward
Barnet:          barnacle
                         
                         
                         
`,
	},

}

func TestClient(t *testing.T) {
	for _, tt := range clientTests {
		t.Run(tt.test, func(t *testing.T) {
			surface := tui.NewTestSurface(25, 5)
			theme := tui.NewTheme()
			p := tui.NewPainter(surface, theme)

			w := tt.setup()
			p.Repaint(w)

			got := surface.String()
			if got != tt.want {
				t.Errorf("unexpected contents: got = \n%s\nwant = \n%s", got, tt.want)
			}
		})
	}
}
