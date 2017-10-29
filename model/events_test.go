// 2017-10-28 cceckman <charles@cceckman.com>
package model_test

import (
	"strings"
	"testing"

	"github.com/cceckman/discoirc/model"
)

var testdata = []model.Event{
	model.Event{
		ID:       model.EventID{Epoch: 2, Seq: 2},
		Contents: "worldY",
	},
	model.Event{
		ID:       model.EventID{Epoch: 1, Seq: 1},
		Contents: "worldX",
	},
	model.Event{
		ID:       model.EventID{Epoch: 2, Seq: 3},
		Contents: "helloZ",
	},
	model.Event{
		ID:       model.EventID{Epoch: 3, Seq: 1},
		Contents: "worldZ",
	},
	model.Event{
		ID:       model.EventID{Epoch: 2, Seq: 1},
		Contents: "helloY",
	},
	model.Event{
		ID:       model.EventID{Epoch: -1, Seq: 1},
		Contents: "helloX",
	},
}

func joinEvents(e model.Events) string {
	var evstr []string
	for _, ev := range e {
		evstr = append(evstr, ev.String())
	}
	return strings.Join(evstr, ", ")
}

func TestEventSort(t *testing.T) {
	evs := model.NewEvents(testdata)
	want := "helloX, worldX, helloY, worldY, helloZ, worldZ"
	got := joinEvents(evs)
	if want != got {
		t.Errorf("want: '%s' got: '%s'", want, got)
	}
}

func TestSelect(t *testing.T) {
	evs := model.NewEvents(testdata)
	for i, c := range []struct {
		minE, maxE int
		minS, maxS uint
		want       string
	}{
		{minE: -2, minS: 0, maxE: -1, maxS: 0, want: ""},
		{minE: 1, minS: 1, maxE: 1, maxS: 1, want: "worldX"},
		{minE: 1, minS: 1, maxE: 2, maxS: 1, want: "worldX, helloY"},
	} {
		sel := model.EventRange{
			Min: model.EventID{Epoch: c.minE, Seq: c.minS},
			Max: model.EventID{Epoch: c.maxE, Seq: c.maxS},
		}
		got := joinEvents(evs.Select(sel))
		if c.want != got {
			t.Errorf("test case: %d want: '%s' got: '%s'", i, c.want, got)
		}
	}

}
