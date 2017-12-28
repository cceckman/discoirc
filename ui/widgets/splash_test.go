package widgets_test

import (
	"testing"

	"github.com/cceckman/discoirc/ui/widgets"
	"github.com/marcusolsson/tui-go"
)

type Quitter struct {
	HasQuit bool
}

func (q *Quitter) Quit() {
	q.HasQuit = true
}

func TestQuit(t *testing.T) {
	q := &Quitter{}
	s := widgets.NewSplash(q)
	s.OnKeyEvent(tui.KeyEvent{
		Key: tui.KeyCtrlC,
	})
	if !q.HasQuit {
		t.Errorf("unexpected state: has not quit")
	}
}
