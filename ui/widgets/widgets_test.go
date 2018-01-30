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
	{Test: "Change fill rune",
		Setup: func() tui.Widget {
			f := widgets.NewFiller('?')
			f.SetFill('!')
			return f
		},
		WantContents: `
!!!!!!!!!!
!!!!!!!!!!
!!!!!!!!!!
!!!!!!!!!!
!!!!!!!!!!
`,
	},
}

func TestFiller(t *testing.T) {
	for _, tt := range FillerTests {
		tt := tt
		t.Run(tt.Test, func(t *testing.T) {
			t.Parallel()
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

type Quitter struct {
	HasQuit bool
}

func (q *Quitter) Quit() {
	q.HasQuit = true
}

func TestSplash_Quit(t *testing.T) {
	t.Parallel()
	q := &Quitter{}
	s := widgets.NewSplash(q)
	// No-op event, just for coverage
	s.OnKeyEvent(tui.KeyEvent{
		Key: tui.KeyCtrlA,
	})
	if q.HasQuit {
		t.Errorf("unexpected state: has quit")
	}

	s.OnKeyEvent(tui.KeyEvent{
		Key: tui.KeyCtrlC,
	})
	if !q.HasQuit {
		t.Errorf("unexpected state: has not quit")
	}
}
