package view_test

import (
	"github.com/cceckman/discoirc/ui/channel/mocks"
	"github.com/cceckman/discoirc/ui/channel/view"
	"github.com/marcusolsson/tui-go"
	"testing"
)

func TestView(t *testing.T) {
	surface := tui.NewTestSurface(40, 10)
	theme := tui.NewTheme()
	theme.SetStyle("reversed", tui.Style{Reverse: tui.DecorationOn})
	p := tui.NewPainter(surface, theme)

	v := view.New(&mocks.UI{})
	v.SetTopic("topic")
	v.SetNick("<nick>")
	v.SetConnection("network: connected")
	v.SetPresence("channel: joined")
	v.SetMode("+v")
	v.SetEvents(mocks.Events)

	p.Repaint(v)

	wantDecorations := `
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
	gotDecorations := surface.Decorations()
	if gotDecorations != wantDecorations {
		t.Errorf("unexpected decorations: got = \n%s\nwant = \n%s", gotDecorations, wantDecorations)
	}

	wantContents := `
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
	gotContents := surface.String()
	if gotContents != wantContents {
		t.Errorf("unexpected contents: got = \n%s\nwant = \n%s", gotContents, wantContents)
	}
}
