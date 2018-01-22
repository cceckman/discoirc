package client_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/client"
	discomocks "github.com/cceckman/discoirc/ui/testhelper"
)

func TestNetwork_NoContents(t *testing.T) {
	w := client.NewNetwork(nil, "Barnetic")

	surface := tui.NewTestSurface(20, 5)
	theme := tui.NewTheme()
	p := tui.NewPainter(surface, theme)
	p.Repaint(w)

	wantContents := `
 Barnetic: ?        
                    
                    
                    
                    
`
	gotContents := surface.String()
	if gotContents != wantContents {
		t.Errorf("unexpected contents: got = \n%s\nwant = \n%s", gotContents, wantContents)
	}
}

func TestNetwork_NoChannels(t *testing.T) {
	w := client.NewNetwork(nil, "Barnetic")
	w.UpdateNetwork(data.NetworkState{
		Nick:  "discobot",
		State: data.Connecting,
	})

	surface := tui.NewTestSurface(25, 5)
	theme := tui.NewTheme()
	p := tui.NewPainter(surface, theme)
	p.Repaint(w)

	wantContents := `
 Barnetic: …     discobot
                         
                         
                         
                         
`
	gotContents := surface.String()
	if gotContents != wantContents {
		t.Errorf("unexpected contents: got = \n%s\nwant = \n%s", gotContents, wantContents)
	}
}

var clientTests = []struct {
	test  string
	setup func(*client.Client)
	want  string
}{
	{
		test:  "blank client",
		setup: func(_ *client.Client) {},
		want: `
                         
                         
                         
                         
                         
                         
                         
                         
                         
                         
`,
	},
	{
		test: "ordered networks",
		setup: func(w *client.Client) {
			w.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "Barnet"},
				Nick:  "barnacle",
			})
			w.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "AlphaNet"},
				Nick:  "edward",
			})
		},
		want: `
 AlphaNet: ∅       edward
 Barnet: ∅       barnacle
                         
                         
                         
                         
                         
                         
                         
                         
`,
	},
	{
		test: "Removed last network",
		setup: func(w *client.Client) {
			w.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "Charlienet"},
				Nick:  "charles",
			})
			w.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "Barnet"},
				Nick:  "barnacle",
			})
			w.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "AlphaNet"},
				Nick:  "edward",
			})
			w.RemoveNetwork("Charlienet")
		},
		want: `
 AlphaNet: ∅       edward
 Barnet: ∅       barnacle
                         
                         
                         
                         
                         
                         
                         
                         
`,
	},
	{
		test: "Removed first network",
		setup: func(w *client.Client) {
			w.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "Charlienet"},
				Nick:  "charles",
			})
			w.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "Barnet"},
				Nick:  "barnacle",
			})
			w.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "AlphaNet"},
				Nick:  "edward",
			})
			w.RemoveNetwork("AlphaNet")
		},
		want: `
 Barnet: ∅       barnacle
 Charlienet: ∅    charles
                         
                         
                         
                         
                         
                         
                         
                         
`,
	},
	{
		test: "Removed middle network",
		setup: func(w *client.Client) {
			w.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "Charlienet"},
				Nick:  "charles",
			})
			w.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "Barnet"},
				Nick:  "barnacle",
			})
			w.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "AlphaNet"},
				Nick:  "edward",
			})
			w.RemoveNetwork("Barnet")
		},
		want: `
 AlphaNet: ∅       edward
 Charlienet: ∅    charles
                         
                         
                         
                         
                         
                         
                         
                         
`,
	},
	{
		test: "networks with channels",
		setup: func(c *client.Client) {
			c.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "AlphaNet"},
				Nick:  "edward",
			})
			c.UpdateChannel(data.ChannelState{
				Scope: data.Scope{
					Net:  "AlphaNet",
					Name: "#discoirc",
				},
				ChannelMode: "+foobar",
				Unread:      99,
				Members:     48,
			})
			c.UpdateChannel(data.ChannelState{
				Scope: data.Scope{
					Net:  "AlphaNet",
					Name: "#tui-go",
				},
				ChannelMode: "+v",
				Unread:      0,
				Members:     3,
			})
			c.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "Charlienet"},
				Nick:  "charles",
				State: data.Connected,
			})
			c.UpdateChannel(data.ChannelState{

				Scope: data.Scope{
					Net:  "Charlienet",
					Name: "#badpuns",
				},
				ChannelMode: "+v",
			})
		},
		want: `
 AlphaNet: ∅       edward
 #discoirc        +foobar
 ✉ 99                48 ☺
 #tui-go               +v
 ✉ 0                  3 ☺
 Charlienet: ✓    charles
 #badpuns              +v
 ✉ 0                  0 ☺
                         
                         
`,
	},
	{
		test: "channel removal",
		setup: func(c *client.Client) {
			c.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "AlphaNet"},
				Nick:  "edward",
			})
			c.UpdateChannel(data.ChannelState{

				Scope: data.Scope{
					Net:  "AlphaNet",
					Name: "#discoirc",
				},
				ChannelMode: "+foobar",
				Unread:      99,
				Members:     48,
			})
			c.UpdateChannel(data.ChannelState{
				Scope: data.Scope{
					Net:  "AlphaNet",
					Name: "#tui-go",
				},
				ChannelMode: "+v",
				Unread:      0,
				Members:     3,
			})

			c.GetNetwork("AlphaNet").RemoveChannel("#tui-go")
		},
		want: `
 AlphaNet: ∅       edward
 #discoirc        +foobar
 ✉ 99                48 ☺
                         
                         
                         
                         
                         
                         
                         
`,
	},
	{
		test: "selected channel, deselected channel",
		setup: func(c *client.Client) {
			c.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "AlphaNet"},
				Nick:  "edward",
			})
			c.UpdateChannel(data.ChannelState{
				Scope: data.Scope{
					Net:  "AlphaNet",
					Name: "#discoirc",
				},
				ChannelMode: "+foobar",
				Unread:      99,
				Members:     48,
			})

			c.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "Charlienet"},
				Nick:  "charles",
				State: data.Connected,
			})
			c.UpdateChannel(data.ChannelState{
				Scope: data.Scope{
					Net:  "Charlienet",
					Name: "#badpuns",
				},
				ChannelMode: "+v",
			})

			c.GetNetwork("Charlienet").GetChannel("#badpuns").SetFocused(true)
			c.GetNetwork("AlphaNet").GetChannel("#discoirc").SetFocused(true)
			c.GetNetwork("AlphaNet").GetChannel("#discoirc").SetFocused(false)

		},
		want: `
 AlphaNet: ∅       edward
 #discoirc        +foobar
 ✉ 99                48 ☺
 Charlienet: ✓    charles
|#badpuns              +v
|✉ 0                  0 ☺
                         
                         
                         
                         
`,
	},
	{
		test: "selected network, deselected network",
		setup: func(c *client.Client) {
			c.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "AlphaNet"},
				Nick:  "edward",
			})
			c.UpdateChannel(data.ChannelState{
				Scope: data.Scope{
					Net:  "AlphaNet",
					Name: "#discoirc",
				},
				ChannelMode: "+foobar",
				Unread:      99,
				Members:     48,
			})
			c.UpdateNetwork(data.NetworkState{
				Scope: data.Scope{Net: "Charlienet"},
				Nick:  "charles",
				State: data.Connected,
			})
			c.UpdateChannel(data.ChannelState{
				Scope: data.Scope{
					Net:  "Charlienet",
					Name: "#badpuns",
				},
				ChannelMode: "+v",
			})

			c.GetNetwork("AlphaNet").SetFocused(true)
			c.GetNetwork("Charlienet").SetFocused(true)
			c.GetNetwork("Charlienet").SetFocused(false)
		},
		want: `
>AlphaNet: ∅       edward
 #discoirc        +foobar
 ✉ 99                48 ☺
 Charlienet: ✓    charles
 #badpuns              +v
 ✉ 0                  0 ☺
                         
                         
                         
                         
`,
	},
}

func TestRender_Client(t *testing.T) {
	for _, tt := range clientTests {
		t.Run(tt.test, func(t *testing.T) {
			surface := tui.NewTestSurface(25, 10)
			theme := tui.NewTheme()
			p := tui.NewPainter(surface, theme)

			ui := discomocks.NewController()

			// Root creation must happen in the main thread
			w := client.New(ui, nil)
			tt.setup(w)
			p.Repaint(w)

			// And run tests
			got := surface.String()
			if got != tt.want {
				t.Errorf("unexpected contents:\ngot = \n%s\n--\nwant = \n%s\n--", got, tt.want)
			}
		})
	}
}

var renderTests = []struct {
	test  string
	setup func() tui.Widget
	want  string
}{
	{
		test: "empty channel",
		setup: func() tui.Widget {
			return tui.NewVBox(
				client.NewChannel(nil, "#discoirc"),
				tui.NewSpacer(),
			)
		},
		want: `
 #discoirc               
 ✉ ?                  ? ☺
                         
                         
                         
                         
                         
                         
                         
                         
`,
	},
	{
		test: "populated channel",
		setup: func() tui.Widget {
			c := client.NewChannel(nil, "#discoirc")
			c.UpdateChannel(data.ChannelState{
				ChannelMode: "+foobar",
				Unread:      99,
				Members:     48,
			})
			return tui.NewVBox(c, tui.NewSpacer())
		},
		want: `
 #discoirc        +foobar
 ✉ 99                48 ☺
                         
                         
                         
                         
                         
                         
                         
                         
`,
	},
}

func TestRender(t *testing.T) {
	for _, tt := range renderTests {
		t.Run(tt.test, func(t *testing.T) {
			surface := tui.NewTestSurface(25, 10)
			theme := tui.NewTheme()
			p := tui.NewPainter(surface, theme)

			w := tt.setup()
			p.Repaint(w)

			got := surface.String()
			if got != tt.want {
				t.Errorf("unexpected contents:\ngot = \n%s\n--\nwant = \n%s\n--", got, tt.want)
			}
		})
	}
}

// namedWidget is a widget with a name.
type namedWidget interface {
	tui.Widget
	Name() string
}

func TestNetwork_FocusNoNetworks(t *testing.T) {
	c := client.New(nil, nil)
	if c.FocusNext(c) != c {
		t.Errorf("unexpected next element for root: got: %v want: %v", c.FocusNext(c), c)
	}
	// TODO Test FocusPrev
}

var FocusTests = []struct {
	Test string
	Case func(*client.Client) []namedWidget
}{
	{
		Test: "no channels",
		Case: func(c *client.Client) []namedWidget {
			gophernet := c.GetNetwork("gophernet")
			kubernet := c.GetNetwork("kubernet")

			return []namedWidget{gophernet, kubernet}
		},
	},
	{
		Test: "channel wraparound",
		Case: func(c *client.Client) []namedWidget {
			gophernet := c.GetNetwork("gophernet")
			kubernet := c.GetNetwork("kubernet")
			metallb := kubernet.GetChannel("#metallb")

			return []namedWidget{gophernet, kubernet, metallb}
		},
	},
	{
		Test: "channel network traversal",
		Case: func(c *client.Client) []namedWidget {
			gophernet := c.GetNetwork("gophernet")
			tuigo := gophernet.GetChannel("#tuigo")
			discoirc := gophernet.GetChannel("#discoirc")
			kubernet := c.GetNetwork("kubernet")
			metallb := kubernet.GetChannel("#metallb")

			return []namedWidget{
				gophernet,
				discoirc,
				tuigo,
				kubernet,
				metallb,
			}
		},
	},
}

func TestNetwork_Focus(t *testing.T) {
	for _, tt := range FocusTests {
		tt := tt
		t.Run(tt.Test, func(t *testing.T) {
			ui := discomocks.NewController()
			c := client.New(ui, nil)

			want := tt.Case(c)
			if len(want) == 0 {
				t.Fatalf("test needs at least one element in wanted list")
			}
			last := want[len(want)-1]

			// Test root
			rootNext := c.FocusNext(c)
			rootPrev := c.FocusPrev(c)
			if rootNext != want[0] {
				t.Errorf("unexpected next element for root: got: %v want: %q", rootNext, want[0].Name())
			}
			if rootPrev != last {
				t.Errorf("unexpected previous element for root: got: %v want: %q", rootPrev, last.Name())
			}

			// Test ordering by walking through
			for i := 0; i < len(want)-1; i++ {
				got := c.FocusNext(want[i]).(namedWidget)
				if got != want[i+1] {
					t.Errorf("unexpected next element for %q: got: %q want: %q", want[i].Name(), got.Name(), want[i+1].Name())
				}
			}
			// Test wrap-around
			got := c.FocusNext(last).(namedWidget)
			if got != want[0] {
				t.Errorf("unexpected next element for %q: got: %q want: %q", last.Name(), got.Name(), want[0].Name())
			}

			for i := len(want) - 1; i > 0; i-- {
				got := c.FocusPrev(want[i]).(namedWidget)
				if got != want[i-1] {
					t.Errorf("unexpected next element for %q: got: %q want: %q", want[i].Name(), got.Name(), want[i-1].Name())
				}
			}
			got = c.FocusPrev(want[0]).(namedWidget)
			if got != last {
				t.Errorf("unexpected next element for %q: got: %q want: %q", want[0].Name(), got.Name(), last.Name())
			}
		})
	}
}

// ActivationTests test response to keypress events.
var ActivationTests = []struct {
	Test              string
	Input             []tui.KeyEvent
	WantView          discomocks.ActiveView
	WantNet, WantChan string
}{
	{
		Test: "hit Down, activate",
		Input: []tui.KeyEvent{
			{Key: tui.KeyDown},
			{Key: tui.KeyDown},
			{Key: tui.KeyEnter},
		},
		WantView: discomocks.ChannelView,
		WantNet:  "gonet", WantChan: "#discoirc",
	},
	{
		Test: "hit J, activate",
		Input: []tui.KeyEvent{
			{Key: tui.KeyRune, Rune: 'j'},
			{Key: tui.KeyRune, Rune: 'j'},
			{Key: tui.KeyEnter},
		},
		WantView: discomocks.ChannelView,
		WantNet:  "gonet", WantChan: "#discoirc",
	},
	{
		Test: "hit K, activate",
		Input: []tui.KeyEvent{
			{Key: tui.KeyRune, Rune: 'k'},
			{Key: tui.KeyEnter},
		},
		WantView: discomocks.ChannelView,
		WantNet:  "zetanet", WantChan: "#bar",
	},
	{
		Test: "hit Up, activate",
		Input: []tui.KeyEvent{
			{Key: tui.KeyUp},
			{Key: tui.KeyEnter},
		},
		WantView: discomocks.ChannelView,
		WantNet:  "zetanet", WantChan: "#bar",
	},
	{
		Test: "no activation on root",
		Input: []tui.KeyEvent{
			{Key: tui.KeyEnter},
		},
		WantView: discomocks.ClientView,
	},

	{
		Test: "no activation on network",
		Input: []tui.KeyEvent{
			{Key: tui.KeyDown},
			{Key: tui.KeyEnter},
		},
		WantView: discomocks.ClientView,
	},
}

func TestNetwork_ActivateChannel(t *testing.T) {
	for _, tt := range ActivationTests {
		tt := tt
		t.Run(tt.Test, func(t *testing.T) {
			ui := discomocks.NewController()
			ui.V = discomocks.ClientView

			root := client.New(ui, nil)

			root.GetNetwork("gonet").GetChannel("#discoirc")
			root.GetNetwork("zetanet").GetChannel("#bar")

			for _, ev := range tt.Input {
				root.OnKeyEvent(ev)
			}

			if ui.V != tt.WantView {
				t.Errorf("unexpected active view: got: %v want: %v", ui.V, tt.WantView)
			}

			if tt.WantNet != "" && ui.Network != tt.WantNet {
				t.Errorf("unexpected active network: got: %q want: %q", ui.Network, tt.WantNet)
			}

			if tt.WantChan != "" && ui.Channel != tt.WantChan {
				t.Errorf("unexpected active channel: got: %q want: %q", ui.Channel, tt.WantChan)
			}

		})
	}

}

func TestNetwork_Quit(t *testing.T) {
	ui := discomocks.NewController()
	root := client.New(ui, nil)

	// The below update itself.
	// It's ok for handlers to run in the main loop.
	root.OnKeyEvent(tui.KeyEvent{
		Key: tui.KeyCtrlC,
	})

	if !ui.HasQuit {
		t.Errorf("client hasn't quit")
	}
}

func Test_Issue18(t *testing.T) {
	// https://github.com/cceckman/discoirc/issues/18
	// Try to reproduce in a test.

	surface := tui.NewTestSurface(132, 2)
	p := tui.NewPainter(surface, tui.NewTheme())

	w := client.NewChannel(nil, "#discoirc")
	w.UpdateChannel(data.ChannelState{
		ChannelMode: "l",
		Unread:      8,
		Members:     8,
	})
	p.Repaint(w)

	want := `
 #discoirc                                                                                                                         l
 ✉ 8                                                                                                                             8 ☺
`
	if got := surface.String(); !cmp.Equal(got, want) {
		t.Error("unexpected contents:")
		t.Error("got:")

		g := bufio.NewScanner(bytes.NewBufferString(got))
		for g.Scan() {
			t.Errorf("%q", g.Text())
		}

		t.Error("want:")
		w := bufio.NewScanner(bytes.NewBufferString(want))
		for w.Scan() {
			t.Errorf("%q", w.Text())
		}
	}

	surface = tui.NewTestSurface(140, 2)
	p = tui.NewPainter(surface, tui.NewTheme())
	p.Repaint(w)
	want = `
 #discoirc                                                                                                                                 l
 ✉ 8                                                                                                                                     8 ☺
`

	if got := surface.String(); !cmp.Equal(got, want) {
		t.Error("unexpected contents:")
		t.Error("got:")

		g := bufio.NewScanner(bytes.NewBufferString(got))
		for g.Scan() {
			t.Errorf("%q", g.Text())
		}

		t.Error("want:")
		w := bufio.NewScanner(bytes.NewBufferString(want))
		for w.Scan() {
			t.Errorf("%q", w.Text())
		}
	}

}
