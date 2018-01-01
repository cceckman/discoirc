package channel_test

import (
	"fmt"
	"image"
	"testing"

	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/cceckman/discoirc/ui/mocks"

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
<yorick>                                
`

var theme = func() *tui.Theme {
	t := tui.NewTheme()
	t.SetStyle("reversed", tui.Style{Reverse: tui.DecorationOn})
	return t
}()

var renderTests = []struct {
	test            string
	setup           func(c *channel.View) func()
	wantContents    string
	wantDecorations string
}{
	{
		test: "base render",
		setup: func(c *channel.View) func() {
			c.UpdateNetwork(data.NetworkState{
				Network: "HamNet",
				State:   data.Connected,
				Nick:    "yorick",
			})
			c.UpdateChannel(data.ChannelState{
				Network:     "HamNet",
				Channel:     "#hamlet",
				Presence:    data.Joined,
				ChannelMode: "+v",
				Topic:       "Act I, Scene 1",
				Unread:      3834, // depending on your editor, of course.
				Members:     12,   // in the company, more characters.
				LastMessage: mocks.Events[len(mocks.Events)-1],
			})
			return func() {}
		},
		wantContents:    wantContents40x10,
		wantDecorations: wantDecor40x10,
	},
	{
		test: "resize render",
		setup: func(c *channel.View) func() {
			c.UpdateNetwork(data.NetworkState{
				Network: "HamNet",
				State:   data.Connected,
				Nick:    "yorick",
			})
			c.UpdateChannel(data.ChannelState{
				Network:     "HamNet",
				Channel:     "#hamlet",
				Presence:    data.Joined,
				ChannelMode: "+v",
				Topic:       "Act I, Scene 1",
				Unread:      3834, // depending on your editor, of course.
				Members:     12,   // in the company, more characters.
				LastMessage: mocks.Events[len(mocks.Events)-1],
			})
			return func() {
				// Will reset to appropriate size when applied to the surface.
				c.Resize(image.Pt(80, 80))
			}
		},
		wantContents:    wantContents40x10,
		wantDecorations: wantDecor40x10,
	},
	{
		test: "underflow render",
		setup: func(c *channel.View) func() {
			c.UpdateNetwork(data.NetworkState{
				Network: "HamNet",
				State:   data.Connected,
				Nick:    "yorick",
			})
			c.UpdateChannel(data.ChannelState{
				Network:     "HamNet",
				Channel:     "#hamlet",
				Presence:    data.Joined,
				ChannelMode: "+v",
				Topic:       "Act I, Scene 1",
				Unread:      3834, // depending on your editor, of course.
				Members:     12,   // in the company, more characters.
				LastMessage: mocks.Events[3],
			})
			return func() {}
		},
		wantContents: `
Act I, Scene 1                          
                                        
                                        
                                        
1,1 TOPIC Act I, Scene 1                
1,2 JOIN barnardo                       
1,3 JOIN francisco                      
1,4 <barnardo> Who's there?             
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

			ui := mocks.NewController()
			defer ui.Close()
			d := mocks.NewBackend()

			var w *channel.View
			// Root creation must happen in the main thread
			ui.RunSync(func() {
				w = channel.New("HamNet", "#hamlet", ui, d)
				w.SetRenderer(testRenderer)
			})
			f := tt.setup(w)
			ui.Update(f)
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

func TestInput_Message(t *testing.T) {
	ui := mocks.NewController()
	defer ui.Close()
	d := mocks.NewBackend()
	// Root creation must happen in the main thread
	ui.RunSync(func() {
		_ = channel.New("HamNet", "#hamlet", ui, d)
	})

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
	ui := mocks.NewController()
	defer ui.Close()
	d := mocks.NewBackend()
	// Root creation must happen in the main thread
	ui.RunSync(func() {
		_ = channel.New("HamNet", "#hamlet", ui, d)
	})

	ui.Type("/quit nothing to see here\n")
	if len(d.Sent) != 0 {
		t.Errorf("message unexpectedly sent: got: %v want: none", d.Sent)
	}

	ui.RunSync(func() {
		if !ui.HasQuit {
			t.Errorf("unexpected state: have not quit")
		}
	})
}

func TestInput_QuitKeybind(t *testing.T) {
	ui := mocks.NewController()
	defer ui.Close()
	d := mocks.NewBackend()
	// Root creation must happen in the main thread
	ui.RunSync(func() {
		_ = channel.New("HamNet", "#hamlet", ui, d)
	})

	ui.Update(func() {
		ui.Root.OnKeyEvent(tui.KeyEvent{
			Key: tui.KeyCtrlC,
		})
	})

	if len(d.Sent) != 0 {
		t.Errorf("message unexpectedly sent: got: %v want: none", d.Sent)
	}

	ui.RunSync(func() {
		if !ui.HasQuit {
			t.Errorf("unexpected state: have not quit")
		}
	})
}

func TestInput_ActivateClient(t *testing.T) {
	ui := mocks.NewController()
	defer ui.Close()
	d := mocks.NewBackend()

	ui.RunSync(func() {
		ui.V = mocks.ChannelView
	})


	// Root creation must happen in the main thread
	ui.RunSync(func() {
		_ = channel.New("HamNet", "#hamlet", ui, d)
	})

	ui.Type("/client\n")

	ui.RunSync(func() {
		if ui.V != mocks.ClientView {
			t.Errorf("unexpected root state: got: %v want: %v", ui.V, mocks.ClientView)
		}
	})
}

