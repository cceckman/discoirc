package view_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel/mocks"
	"github.com/cceckman/discoirc/ui/channel/view"
	"github.com/cceckman/discoirc/ui/channel"

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
topic                                   
vnfold your selfe                       
1,6 <barnardo> Long liue the King       
2,1 <claudius> Welcome, dear Rosencrantz
and Guildenstern!                       
2,2 <gertrude> Good gentlemen, he hath  
much talk'd of you;                     
2,3 <rosencrantz> Both your majesties   
network: connected channel: joined +v   
<nick>                                  
`

func testRenderer(e data.Event) tui.Widget {
	r := tui.NewLabel(fmt.Sprintf("%d,%d %s", e.Epoch, e.Seq, e.Contents))
	r.SetWordWrap(true)
	return r
}

func makeView() channel.View {
	v := view.New(&mocks.UI{})
	v.SetTopic("topic")
	v.SetNick("<nick>")
	v.SetConnection("network: connected")
	v.SetPresence("channel: joined")
	v.SetMode("+v")
	v.SetRenderer(testRenderer)
	v.SetEvents(mocks.Events)
	return v
}

func TestView_SimpleRender(t *testing.T) {
	surface := tui.NewTestSurface(40, 10)
	theme := tui.NewTheme()
	theme.SetStyle("reversed", tui.Style{Reverse: tui.DecorationOn})
	p := tui.NewPainter(surface, theme)

	v := makeView()
	p.Repaint(v)

	gotDecorations := surface.Decorations()
	if gotDecorations != wantDecor40x10 {
		t.Errorf("unexpected decorations: got = \n%s\nwant = \n%s", gotDecorations, wantDecor40x10)
	}

	gotContents := surface.String()
	if gotContents != wantContents40x10 {
		t.Errorf("unexpected contents: got = \n%s\nwant = \n%s", gotContents, wantContents40x10)
	}
}

func TestView_Resize(t *testing.T) {
	// Render with a different size.
	preSurface := tui.NewTestSurface(80, 20)
	theme := tui.NewTheme()
	theme.SetStyle("reversed", tui.Style{Reverse: tui.DecorationOn})
	prePainter := tui.NewPainter(preSurface, theme)

	surface := tui.NewTestSurface(40, 10)
	p := tui.NewPainter(surface, theme)

	v := makeView()

	// 1: Resize without controller set.
	prePainter.Repaint(v)

	// 2: Attach controller and assert size was set.
	c := &mocks.Controller{}
	v.Attach(c)
	wantSize := 17 // 20 - topic, status, input
	if c.Size != wantSize {
		t.Errorf("unexpected size: got = %d want = %d", c.Size, wantSize)
	}

	// 3: Resize; check that output is good, and that controller saw resize.

	p.Repaint(v)
	wantSize = 7 // 10 - topic, status, input
	if c.Size != wantSize {
		t.Errorf("unexpected size: got = %d want = %d", c.Size, wantSize)
	}

	gotDecorations := surface.Decorations()
	if gotDecorations != wantDecor40x10 {
		t.Errorf("unexpected decorations: got = \n%s\nwant = \n%s", gotDecorations, wantDecor40x10)
	}

	gotContents := surface.String()
	if gotContents != wantContents40x10 {
		t.Errorf("unexpected contents: got = \n%s\nwant = \n%s", gotContents, wantContents40x10)
	}
}

func TestView_Underfill(t *testing.T) {
	v := makeView()
	v.SetEvents(mocks.Events[len(mocks.Events)-2:])

	surface := tui.NewTestSurface(40, 10)
	p := tui.NewPainter(surface, tui.NewTheme())

	p.Repaint(v)
	wantContents := `
topic                                   
                                        
                                        
                                        
                                        
2,2 <gertrude> Good gentlemen, he hath  
much talk'd of you;                     
2,3 <rosencrantz> Both your majesties   
network: connected channel: joined +v   
<nick>                                  
`
		gotContents := surface.String()
	if gotContents != wantContents {
		t.Errorf("unexpected contents: got = \n%s\nwant = \n%s", gotContents, wantContents)
	}

}

func TestView_Input(t *testing.T) {
	v := makeView()
	c := &mocks.Controller{}
	v.Attach(c)
	want := []string{"message one", "/me sends a message"}
	inputs := strings.Join(want, "\n")

	for _, rn := range inputs {
		var ev tui.KeyEvent
		if rn != '\n' {
			ev = tui.KeyEvent{
				Key: tui.KeyRune,
				Rune: rn,
			}
		} else {
			ev = tui.KeyEvent{
				Key: tui.KeyEnter,
			}
		}
		v.OnKeyEvent(ev)
	}

	if len(c.Received) != len(want) {
		t.Errorf("unexpected messages: got = %v want %v",  c.Received, want)
	} else {
		for i, msg := range want {
			got := c.Received[i]
			if got != msg {
				t.Errorf("unexpected contents in message %d: got = %q want %q", i, got, msg)
			}
		}
	}
}
