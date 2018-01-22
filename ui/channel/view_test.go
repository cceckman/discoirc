package channel_test

import (
	"fmt"
	"image"
	"testing"

	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/cceckman/discoirc/ui/testhelper"

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
6 <barnardo> Long liue the King         
7 <claudius> Welcome, dear Rosencrantz  
and Guildenstern!                       
8 <gertrude> Good gentlemen, he hath    
much talk'd of you;                     
9 <rosencrantz> Both your majesties     
HamNet: ✓ #hamlet: +v                   
<yorick>                                
`

var theme = func() *tui.Theme {
	t := tui.NewTheme()
	t.SetStyle("reversed", tui.Style{Reverse: tui.DecorationOn})
	return t
}()

var renderTests = []struct {
	test            string
	setup           func(c *channel.View)
	wantContents    string
	wantDecorations string
}{
	{
		test: "base render",
		setup: func(c *channel.View) {
			c.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "HamNet"},
				State: data.Connected,
				Nick:  "yorick",
			})
			c.UpdateChannel(data.ChannelState{
				Scope: data.Scope{
					Net:  "HamNet",
					Name: "#hamlet",
				},
				Presence:    data.Joined,
				ChannelMode: "+v",
				Topic:       "Act I, Scene 1",
				Unread:      3834, // depending on your editor, of course.
				Members:     12,   // in the company, more characters.
				LastMessage: testhelper.Events[len(testhelper.Events)-1].Seq(),
			})
		},
		wantContents:    wantContents40x10,
		wantDecorations: wantDecor40x10,
	},
	{
		test: "resize render",
		setup: func(c *channel.View) {
			c.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "HamNet"},
				State: data.Connected,
				Nick:  "yorick",
			})
			c.UpdateChannel(data.ChannelState{
				Scope: data.Scope{
					Net:  "HamNet",
					Name: "#hamlet",
				},
				Presence:    data.Joined,
				ChannelMode: "+v",
				Topic:       "Act I, Scene 1",
				Unread:      3834, // depending on your editor, of course.
				Members:     12,   // in the company, more characters.
				LastMessage: testhelper.Events[len(testhelper.Events)-1].Seq(),
			})

			c.Resize(image.Pt(80, 80))
		},
		wantContents:    wantContents40x10,
		wantDecorations: wantDecor40x10,
	},
	{
		test: "underflow render",
		setup: func(c *channel.View) {
			c.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "HamNet"},
				State: data.Connected,
				Nick:  "yorick",
			})
			c.UpdateChannel(data.ChannelState{
				Scope: data.Scope{
					Net:  "HamNet",
					Name: "#hamlet",
				},
				Presence:    data.Joined,
				ChannelMode: "+v",
				Topic:       "Act I, Scene 1",
				Unread:      3834, // depending on your editor, of course.
				Members:     12,   // in the company, more characters.
				LastMessage: testhelper.Events[3].Seq(),
			})
		},
		wantContents: `
Act I, Scene 1                          
                                        
                                        
                                        
1 TOPIC Act I, Scene 1                  
2 JOIN barnardo                         
3 JOIN francisco                        
4 <barnardo> Who's there?               
HamNet: ✓ #hamlet: +v                   
<yorick>                                
`,
		wantDecorations: wantDecor40x10,
	},
}

func TestRender(t *testing.T) {
	for _, tt := range renderTests {
		t.Run(tt.test, func(t *testing.T) {
			surface := tui.NewTestSurface(40, 10)
			p := tui.NewPainter(surface, theme)

			ui := testhelper.NewController()
			d := testhelper.NewBackend()

			// Root creation must happen in the main thread
			w := channel.New(data.Scope{Net: "HamNet", Name: "#hamlet"}, ui, d)
			w.SetRenderer(testRenderer)
			tt.setup(w)
			p.Repaint(w)

			gotContents := surface.String()
			if tt.wantContents != "" && gotContents != tt.wantContents {
				t.Errorf("unexpected contents:\ngot = \n%s\n--\nwant = \n%s\n--", gotContents, tt.wantContents)
			}
			gotDecorations := surface.Decorations()
			if tt.wantDecorations != "" && gotDecorations != tt.wantDecorations {
				t.Errorf("unexpected decorations:\ngot = \n%s\n--\nwant = \n%s\n--", gotDecorations, tt.wantDecorations)
			}
		})

	}
}

func testRenderer(e data.Event) tui.Widget {
	r := tui.NewLabel(fmt.Sprintf("%d %s", e.Seq(), e.String()))
	r.SetWordWrap(true)
	return r
}

func TestInput_Message(t *testing.T) {
	ui := testhelper.NewController()
	d := testhelper.NewBackend()
	_ = channel.New(data.Scope{Net: "HamNet", Name: "#hamlet"}, ui, d)

	ui.Type("hello everyone")

	if len(d.Sent) != 0 {
		t.Errorf("message unexpectedly sent: got: %v want: none", d.Sent)
	}

	ui.Type("!\nhow are you")
	if len(d.Sent) != 1 || d.Sent[0] != "hello everyone!" {
		t.Errorf("unexpected messages sent: got: %v want: %q", d.Sent, "hello everyone!")
	}
}

func TestInput_QuitMessage(t *testing.T) {
	ui := testhelper.NewController()
	d := testhelper.NewBackend()
	_ = channel.New(data.Scope{Net: "HamNet", Name: "#hamlet"}, ui, d)

	ui.Type("/quit nothing to see here\n")
	if len(d.Sent) != 0 {
		t.Errorf("message unexpectedly sent: got: %v want: none", d.Sent)
	}

	if !ui.HasQuit {
		t.Errorf("unexpected state: have not quit")
	}
}

func TestInput_QuitKeybind(t *testing.T) {
	ui := testhelper.NewController()
	d := testhelper.NewBackend()
	_ = channel.New(data.Scope{Net: "HamNet", Name: "#hamlet"}, ui, d)

	ui.Root.OnKeyEvent(tui.KeyEvent{
		Key: tui.KeyCtrlC,
	})

	if len(d.Sent) != 0 {
		t.Errorf("message unexpectedly sent: got: %v want: none", d.Sent)
	}

	if !ui.HasQuit {
		t.Errorf("unexpected state: have not quit")
	}
}

func TestInput_ActivateClient(t *testing.T) {
	ui := testhelper.NewController()
	d := testhelper.NewBackend()

	ui.V = testhelper.ChannelView

	// Root creation must happen in the main thread
	_ = channel.New(data.Scope{Net: "HamNet", Name: "#hamlet"}, ui, d)

	ui.Type("/client\n")

	if ui.V != testhelper.ClientView {
		t.Errorf("unexpected root state: got: %v want: %v", ui.V, testhelper.ClientView)
	}
}
