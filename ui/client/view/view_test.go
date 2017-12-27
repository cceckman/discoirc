package view_test

import (
	"testing"

	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/ui/client"
	"github.com/cceckman/discoirc/ui/client/view"
	discomocks "github.com/cceckman/discoirc/ui/mocks"
)

func TestNetwork_NoContents(t *testing.T) {
	w := view.NewNetwork(nil, "Barnetic")

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
	w := view.NewNetwork(nil, "Barnetic")
	w.SetNick("discobot")
	w.SetConnection("…")

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

var renderTests = []struct {
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
 AlphaNet:         edward
 Barnet:         barnacle
                         
                         
                         
                         
                         
                         
                         
                         
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
 AlphaNet:         edward
 Barnet:         barnacle
                         
                         
                         
                         
                         
                         
                         
                         
`,
	},
	{
		test: "Removed first network",
		setup: func() tui.Widget {
			w := view.New()
			w.GetNetwork("Charlienet").SetNick("charles")
			w.GetNetwork("Barnet").SetNick("barnacle")
			w.GetNetwork("AlphaNet").SetNick("edward")
			w.RemoveNetwork("AlphaNet")
			return w
		},
		want: `
 Barnet:         barnacle
 Charlienet:      charles
                         
                         
                         
                         
                         
                         
                         
                         
`,
	},
	{
		test: "Removed middle network",
		setup: func() tui.Widget {
			w := view.New()
			w.GetNetwork("Charlienet").SetNick("charles")
			w.GetNetwork("Barnet").SetNick("barnacle")
			w.GetNetwork("AlphaNet").SetNick("edward")
			w.RemoveNetwork("Barnet")
			return w
		},
		want: `
 AlphaNet:         edward
 Charlienet:      charles
                         
                         
                         
                         
                         
                         
                         
                         
`,
	},
	{
		test: "empty channel",
		setup: func() tui.Widget {
			return tui.NewVBox(
				view.NewChannel(nil, "#discoirc"),
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
			c := view.NewChannel(nil, "#discoirc")
			c.SetMode("+foobar")
			c.SetUnread(99)
			c.SetMembers(48)
			return tui.NewVBox(c, tui.NewSpacer())
		},
		want: `
 #discoirc        +foobar
 ✉ 99                48 ☺
                         
                         
                         
                         
                         
                         
                         
                         
`,
	},
	{
		test: "networks with channels",
		setup: func() tui.Widget {
			c := view.New()

			alpha := c.GetNetwork("AlphaNet")
			alpha.SetNick("edward")

			discoirc := alpha.GetChannel("#discoirc")
			discoirc.SetMode("+foobar")
			discoirc.SetUnread(99)
			discoirc.SetMembers(48)

			tuigo := alpha.GetChannel("#tui-go")
			tuigo.SetMode("+v")
			tuigo.SetUnread(0)
			tuigo.SetMembers(3)

			charlie := c.GetNetwork("Charlienet")
			charlie.SetNick("charles")
			charlie.SetConnection("✓")

			badpuns := charlie.GetChannel("#badpuns")
			badpuns.SetMode("+v")

			return c
		},
		want: `
 AlphaNet:         edward
 #discoirc        +foobar
 ✉ 99                48 ☺
 #tui-go               +v
 ✉ 0                  3 ☺
 Charlienet: ✓    charles
 #badpuns              +v
 ✉ ?                  ? ☺
                         
                         
`,
	},
	{
		test: "channel removal",
		setup: func() tui.Widget {
			c := view.New()

			alpha := c.GetNetwork("AlphaNet")
			alpha.SetNick("edward")

			discoirc := alpha.GetChannel("#discoirc")
			discoirc.SetMode("+foobar")
			discoirc.SetUnread(99)
			discoirc.SetMembers(48)

			tuigo := alpha.GetChannel("#tui-go")
			tuigo.SetMode("+v")
			tuigo.SetUnread(0)
			tuigo.SetMembers(3)

			alpha.RemoveChannel("#tui-go")

			return c
		},
		want: `
 AlphaNet:         edward
 #discoirc        +foobar
 ✉ 99                48 ☺
                         
                         
                         
                         
                         
                         
                         
`,
	},
	{
		test: "selected channel, deselected channel",
		setup: func() tui.Widget {
			c := view.New()

			alpha := c.GetNetwork("AlphaNet")
			alpha.SetNick("edward")

			discoirc := alpha.GetChannel("#discoirc")
			discoirc.SetMode("+foobar")
			discoirc.SetUnread(99)
			discoirc.SetMembers(48)

			charlie := c.GetNetwork("Charlienet")
			charlie.SetNick("charles")
			charlie.SetConnection("✓")

			badpuns := charlie.GetChannel("#badpuns")
			badpuns.SetMode("+v")

			badpuns.SetFocused(true)
			discoirc.SetFocused(true)
			discoirc.SetFocused(false)

			return c
		},
		want: `
 AlphaNet:         edward
 #discoirc        +foobar
 ✉ 99                48 ☺
 Charlienet: ✓    charles
|#badpuns              +v
|✉ ?                  ? ☺
                         
                         
                         
                         
`,
	},

	{
		test: "selected network, deselected network",
		setup: func() tui.Widget {
			c := view.New()

			alpha := c.GetNetwork("AlphaNet")
			alpha.SetNick("edward")

			discoirc := alpha.GetChannel("#discoirc")
			discoirc.SetMode("+foobar")
			discoirc.SetUnread(99)
			discoirc.SetMembers(48)

			charlie := c.GetNetwork("Charlienet")
			charlie.SetNick("charles")
			charlie.SetConnection("✓")

			badpuns := charlie.GetChannel("#badpuns")
			badpuns.SetMode("+v")

			alpha.SetFocused(true)
			charlie.SetFocused(true)
			charlie.SetFocused(false)

			return c
		},
		want: `
>AlphaNet:         edward
 #discoirc        +foobar
 ✉ 99                48 ☺
 Charlienet: ✓    charles
 #badpuns              +v
 ✉ ?                  ? ☺
                         
                         
                         
                         
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
	c := view.New()
	if c.FocusNext(c) != c {
		t.Errorf("unexpected next element for root: got: %v want: %v", c.FocusNext(c), c)
	}
	// TODO Test FocusPrev
}

var FocusTests = []struct {
	Test string
	Case func() (*view.Client, []namedWidget)
}{
	{
		Test: "no channels",
		Case: func() (*view.Client, []namedWidget) {
			c := view.New()
			gophernet := c.GetNetwork("gophernet")
			kubernet := c.GetNetwork("kubernet")

			return c, []namedWidget{gophernet, kubernet}
		},
	},
	{
		Test: "channel wraparound",
		Case: func() (*view.Client, []namedWidget) {
			c := view.New()
			gophernet := c.GetNetwork("gophernet")
			kubernet := c.GetNetwork("kubernet")
			metallb := kubernet.GetChannel("#metallb")

			return c, []namedWidget{gophernet, kubernet, metallb}
		},
	},
	{
		Test: "channel network traversal",
		Case: func() (*view.Client, []namedWidget) {
			c := view.New()
			gophernet := c.GetNetwork("gophernet")
			tuigo := gophernet.GetChannel("#tuigo")
			discoirc := gophernet.GetChannel("#discoirc")
			kubernet := c.GetNetwork("kubernet")
			metallb := kubernet.GetChannel("#metallb")

			return c, []namedWidget{
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
			c, want := tt.Case()
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
	Setup             func() client.View
	Input             []tui.KeyEvent
	WantView          discomocks.ActiveView
	WantNet, WantChan string
}{
	{
		Test: "hit Down, activate",
		Setup: func() client.View {
			root := view.New()
			root.GetNetwork("gonet").GetChannel("#discoirc")
			root.GetNetwork("zetanet").GetChannel("#bar")
			return root

		},
		Input: []tui.KeyEvent{
			tui.KeyEvent{Key: tui.KeyDown},
			tui.KeyEvent{Key: tui.KeyDown},
			tui.KeyEvent{Key: tui.KeyEnter},
		},
		WantView: discomocks.ChannelView,
		WantNet:  "gonet", WantChan: "#discoirc",
	},
	{
		Test: "hit J, activate",
		Setup: func() client.View {
			root := view.New()
			root.GetNetwork("gonet").GetChannel("#discoirc")
			root.GetNetwork("zetanet").GetChannel("#bar")
			return root

		},
		Input: []tui.KeyEvent{
			tui.KeyEvent{Key: tui.KeyRune, Rune: 'j'},
			tui.KeyEvent{Key: tui.KeyRune, Rune: 'j'},
			tui.KeyEvent{Key: tui.KeyEnter},
		},
		WantView: discomocks.ChannelView,
		WantNet:  "gonet", WantChan: "#discoirc",
	},
	{
		Test: "hit K, activate",
		Setup: func() client.View {
			root := view.New()
			root.GetNetwork("gonet").GetChannel("#discoirc")
			root.GetNetwork("zetanet").GetChannel("#bar")
			return root

		},
		Input: []tui.KeyEvent{
			tui.KeyEvent{Key: tui.KeyRune, Rune: 'k'},
			tui.KeyEvent{Key: tui.KeyEnter},
		},
		WantView: discomocks.ChannelView,
		WantNet:  "zetanet", WantChan: "#bar",
	},
	{
		Test: "hit Up, activate",
		Setup: func() client.View {
			root := view.New()
			root.GetNetwork("gonet").GetChannel("#discoirc")
			root.GetNetwork("zetanet").GetChannel("#bar")
			return root

		},
		Input: []tui.KeyEvent{
			tui.KeyEvent{Key: tui.KeyUp},
			tui.KeyEvent{Key: tui.KeyEnter},
		},
		WantView: discomocks.ChannelView,
		WantNet:  "zetanet", WantChan: "#bar",
	},

	{
		Test: "no activation on root",
		Setup: func() client.View {
			root := view.New()
			root.GetNetwork("gonet").GetChannel("#discoirc")
			root.GetNetwork("zetanet").GetChannel("#bar")
			return root

		},
		Input: []tui.KeyEvent{
			tui.KeyEvent{Key: tui.KeyEnter},
		},
		WantView: discomocks.ClientView,
	},

	{
		Test: "no activation on network",
		Setup: func() client.View {
			root := view.New()
			root.GetNetwork("gonet").GetChannel("#discoirc")
			root.GetNetwork("zetanet").GetChannel("#bar")
			return root

		},
		Input: []tui.KeyEvent{
			tui.KeyEvent{Key: tui.KeyDown},
			tui.KeyEvent{Key: tui.KeyEnter},
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

			root := tt.Setup()
			root.Attach(ui)
			ui.SetWidget(root)

			if tt.WantView != discomocks.ClientView {
				// Expect an Update
				ui.Add(1)
			}
			for _, ev := range tt.Input {
				root.OnKeyEvent(ev)
			}

			ui.RunSync(func(){
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

		})
	}

}

func TestNetwork_Quit(t *testing.T) {
	root := view.New()
	ui := discomocks.NewController()
	root.Attach(ui)

	ui.Add(1)
	// The below update itself.
	// It's ok for handlers to run in the main loop.
	ui.Update(func() {
		root.OnKeyEvent(tui.KeyEvent{
			Key: tui.KeyCtrlC,
		})
	})

	ui.RunSync(func() {
		if ! ui.HasQuit {
			t.Errorf("client hasn't quit")
		}
	})
}
