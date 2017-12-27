package widgets_test

import (
	"github.com/cceckman/discoirc/ui/widgets"
	"github.com/marcusolsson/tui-go"
	"testing"
)

var FillerTests = []struct {
	Test         string
	Setup        func() tui.Widget
	WantContents string
}{
	{
		Test: "Expand X and Y",
		Setup: func() tui.Widget {
			f := widgets.NewFiller('?')
			f.SetSizePolicy(tui.Expanding, tui.Expanding)
			return tui.NewVBox(tui.NewHBox(f))
		},
		WantContents: `
??????????
??????????
??????????
??????????
??????????
`,
	},
	{
		Test: "Expand Y only",
		Setup: func() tui.Widget {
			f := widgets.NewFiller('?')
			f.SetSizePolicy(tui.Maximum, tui.Expanding)
			l := tui.NewVBox(tui.NewLabel("text"))
			return tui.NewHBox(f, l)
		},
		WantContents: `
?text     
?         
?         
?         
?         
`,
	},
	{
		Test: "Expand X only",
		Setup: func() tui.Widget {
			f := widgets.NewFiller('?')
			f.SetSizePolicy(tui.Expanding, tui.Maximum)
			l := tui.NewHBox(tui.NewLabel("text"))
			return tui.NewVBox(f, l)
		},
		WantContents: `
??????????
text      
          
          
          
`,
	},
}

func TestFoo(t *testing.T) {
	for _, tt := range FillerTests {
		tt := tt
		t.Run(tt.Test, func(t *testing.T) {
			surface := tui.NewTestSurface(10, 5)
			painter := tui.NewPainter(surface, tui.NewTheme())

			w := tt.Setup()
			painter.Repaint(w)
			gotContents := surface.String()

			if gotContents != tt.WantContents {
				t.Errorf("unexpected contents: got: \n%s\n--\nwant: \n%s\n--", gotContents, tt.WantContents)
			}

		})
	}
}
