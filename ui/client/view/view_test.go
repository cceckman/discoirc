package view_test

import (
	"testing"

	"github.com/cceckman/discoirc/ui/client/view"
	"github.com/marcusolsson/tui-go"
)

func TestNetwork_NoContents(t *testing.T) {
	w := view.NewNetwork("Barnetic")

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
	w := view.NewNetwork("Barnetic")
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
