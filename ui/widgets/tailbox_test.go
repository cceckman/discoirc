package widgets_test

import (
	"github.com/cceckman/discoirc/ui/widgets"
	"github.com/marcusolsson/tui-go"
	"testing"
)

var TailBoxTests = []struct {
	Test  string
	Setup func() tui.Widget
	Want  string
}{
	{
		Test: "draw small labels",
		Setup: func() tui.Widget {
			return widgets.NewTailBox(
				tui.NewLabel("hello mom"),
				tui.NewLabel("hello dad"),
			)
		},
		Want: `
          
          
          
hello mom 
hello dad 
`,
	},
	{
		Test: "draw unwrapped labels",
		Setup: func() tui.Widget {
			l1, l2 := tui.NewLabel("hello muddah"), tui.NewLabel("hello faddah")
			return widgets.NewTailBox(l1, l2)
		},
		Want: `
          
          
          
hello mudd
hello fadd
`,
	},
	{
		Test: "draw wrapped labels",
		Setup: func() tui.Widget {
			l1, l2 := tui.NewLabel("hello muddah"), tui.NewLabel("hello faddah")
			l1.SetWordWrap(true)
			l2.SetWordWrap(true)
			return widgets.NewTailBox(l1, l2)
		},
		Want: `
          
hello     
muddah    
hello     
faddah    
`,
	},
}

func TestTailBox(t *testing.T) {
	for _, tt := range TailBoxTests {
		tt := tt
		t.Run(tt.Test, func(t *testing.T) {
			surface := tui.NewTestSurface(10, 5)
			p := tui.NewPainter(surface, tui.NewTheme())
			p.Repaint(tt.Setup())

			if surface.String() != tt.Want {
				t.Errorf("unexpected contents: got = \n%s\nwant = \n%s", surface.String(), tt.Want)
			}
		})
	}
}
