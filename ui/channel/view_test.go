package channel_test

import (
	"fmt"
	"testing"

	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/cceckman/discoirc/ui/channel/mocks"
	discomocks "github.com/cceckman/discoirc/ui/mocks"

	"github.com/marcusolsson/tui-go"
)

const wantDecor40x10 = `
1111111111111111111111111111111111111111
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000
1111111111111111111111111111111111111111
0000000000000000000000000000000000000000
`
const wantContents40x10 = `
Act I, Scene 1                          
vnfold your selfe                       
1,6 <barnardo> Long liue the King       
2,1 <claudius> Welcome, dear Rosencrantz
and Guildenstern!                       
2,2 <gertrude> Good gentlemen, he hath  
much talk'd of you;                     
2,3 <rosencrantz> Both your majesties   
HamNet: ✓ #hamlet: +v                   
<nick>                                  
`

var renderTests = []struct {
	test            string
	setup           func(c *channel.View)
	wantContents    string
	wantDecorations string
}{
	{
		test: "base render",
		setup: func(c *channel.View) {
			c.SetRenderer(testRenderer)
		},
		wantContents:    wantContents40x10,
		wantDecorations: wantDecor40x10,
	},
	// TODO: Resize
	// TODO: Underflow - not enough events
}

func TestRender(t *testing.T) {
	for _, tt := range renderTests {
		t.Run(tt.test, func(t *testing.T) {
			surface := tui.NewTestSurface(40, 10)
			theme := tui.NewTheme()
			p := tui.NewPainter(surface, theme)

			ui := discomocks.NewController()
			defer ui.Close()
			d := mocks.NewBackend()

			var w *channel.View
			// Root creation must happen in the main thread
			ui.RunSync(func() {
				w = channel.NewView("HamNet", "#hamlet", ui, d)
			})
			tt.setup(w)
			// Render in the UI thread so that the race detector works properly.
			ui.Update(func() {
				p.Repaint(w)
			})

			// And run tests
			ui.RunSync(func() {
				gotContents := surface.String()
				if tt.wantContents != "" && gotContents != tt.wantContents {
					t.Errorf("unexpected contents:\ngot = \n%s\n--\nwant = \n%s\n--", gotContents, tt.wantContents)
				}
				gotDecorations := surface.Decorations()
				if tt.wantDecorations != "" && gotDecorations != tt.wantDecorations {
					t.Errorf("unexpected decorations:\ngot = \n%s\n--\nwant = \n%s\n--", gotDecorations, tt.wantDecorations)
				}
			})
		})

	}
}

func testRenderer(e data.Event) tui.Widget {
	r := tui.NewLabel(fmt.Sprintf("%d,%d %s", e.Epoch, e.Seq, e.Contents))
	r.SetWordWrap(true)
	return r
}

// TypeString issues KeyEvents to the Widget as if the provided string had been typed.
func TypeString(w tui.Widget, s string) {
	for _, rn := range s {
		var ev tui.KeyEvent
		if rn != '\n' {
			ev = tui.KeyEvent{
				Key:  tui.KeyRune,
				Rune: rn,
			}
		} else {
			ev = tui.KeyEvent{
				Key: tui.KeyEnter,
			}
		}
		w.OnKeyEvent(ev)
	}
}

// TODO: Redo typing tests:
// - Send message
// - Quit by message, quit by ctrl+c
// - Go to client
