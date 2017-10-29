// 2017-10-28 cceckman <charles@cceckman.com>
package model_test

import (
	"strings"
	"testing"

	"github.com/cceckman/discoirc/model"
)

func TestEventSort(t *testing.T) {
	list := []model.Event{
		model.Event{
			ID:       model.EventID{Epoch: 2, Seq: 2},
			Contents: "worldX",
		},
		model.Event{
			ID:       model.EventID{Epoch: 1, Seq: 1},
			Contents: "worldY",
		},
		model.Event{
			ID:       model.EventID{Epoch: 2, Seq: 1},
			Contents: "helloX",
		},
		model.Event{
			ID:       model.EventID{Epoch: -1, Seq: 1},
			Contents: "helloY",
		},
	}

	evs := model.NewEvents(list)
	want := "helloY, worldY, helloX, worldX"
	var evstr []string
	for _, ev := range evs {
		evstr = append(evstr, ev.String())
	}
	got := strings.Join(evstr, ", ")
	if want != got {
		t.Errorf("want: %s got: %s", want, got)
	}

}
