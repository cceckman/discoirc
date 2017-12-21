// 2017-10-28 cceckman <charles@cceckman.com>
package data_test

import (
	"strings"
	"testing"

	"github.com/cceckman/discoirc/data"
)

var testdata = []data.Event{
	data.Event{
		EventID:  data.EventID{Epoch: 2, Seq: 2},
		Contents: "worldY",
	},
	data.Event{
		EventID:  data.EventID{Epoch: 1, Seq: 1},
		Contents: "worldX",
	},
	data.Event{
		EventID:  data.EventID{Epoch: 2, Seq: 3},
		Contents: "helloZ",
	},
	data.Event{
		EventID:  data.EventID{Epoch: 3, Seq: 1},
		Contents: "worldZ",
	},
	data.Event{
		EventID:  data.EventID{Epoch: 2, Seq: 1},
		Contents: "helloY",
	},
	data.Event{
		EventID:  data.EventID{Epoch: -1, Seq: 1},
		Contents: "helloX",
	},
}

func joinEvents(e []data.Event) string {
	var evstr []string
	for _, ev := range e {
		evstr = append(evstr, ev.String())
	}
	return strings.Join(evstr, ", ")
}

func TestSelectMinSize(t *testing.T) {
	evs := data.NewEvents(testdata)
	for i, c := range []struct {
		N     uint
		Epoch int
		Seq   uint
		Want  string
	}{
		{N: 0, Epoch: 2, Seq: 1, Want: ""},
		{N: 1, Epoch: 2, Seq: 1, Want: "helloY"},
		{N: 2, Epoch: 2, Seq: 0, Want: "helloY, worldY"},
		{N: 2, Epoch: 4, Seq: 1, Want: ""},
		{N: 10, Epoch: 2, Seq: 1, Want: "helloY, worldY, helloZ, worldZ"},
	} {
		id := data.EventID{Epoch: c.Epoch, Seq: c.Seq}
		got := joinEvents(evs.SelectMinSize(id, c.N))
		if c.Want != got {
			t.Errorf("test case: %d want: '%s' got: '%s'", i, c.Want, got)
		}
	}
}

func TestSelectSizeMax(t *testing.T) {
	evs := data.NewEvents(testdata)
	for i, c := range []struct {
		N     uint
		Epoch int
		Seq   uint
		Want  string
	}{
		{N: 0, Epoch: 2, Seq: 1, Want: ""},
		{N: 1, Epoch: 2, Seq: 1, Want: "helloY"},
		{N: 2, Epoch: 2, Seq: 0, Want: "helloX, worldX"},
		{N: 2, Epoch: 4, Seq: 1, Want: "helloZ, worldZ"},
		{N: 10, Epoch: 2, Seq: 1, Want: "helloX, worldX, helloY"},
	} {
		id := data.EventID{Epoch: c.Epoch, Seq: c.Seq}
		got := joinEvents(evs.SelectSizeMax(c.N, id))
		if c.Want != got {
			t.Errorf("test case: %d want: '%s' got: '%s'", i, c.Want, got)
		}
	}
}

func TestSelectSize(t *testing.T) {
	evs := data.NewEvents(testdata)
	for i, c := range []struct {
		N    uint
		Want string
	}{
		{N: 0, Want: ""},
		{N: 1, Want: "worldZ"},
		{N: 2, Want: "helloZ, worldZ"},
		{N: 10, Want: "helloX, worldX, helloY, worldY, helloZ, worldZ"},
	} {
		got := joinEvents(evs.SelectSize(c.N))
		if c.Want != got {
			t.Errorf("test case: %d want: '%s' got: '%s'", i, c.Want, got)
		}
	}
}

func TestEventSort(t *testing.T) {
	evs := data.NewEvents(testdata)
	want := "helloX, worldX, helloY, worldY, helloZ, worldZ"
	got := joinEvents(evs)
	if want != got {
		t.Errorf("want: '%s' got: '%s'", want, got)
	}
}

func TestSelectMinMax(t *testing.T) {
	evs := data.NewEvents(testdata)
	for i, c := range []struct {
		minE, maxE int
		minS, maxS uint
		want       string
	}{
		{minE: -2, minS: 0, maxE: -1, maxS: 0, want: ""},
		{minE: 1, minS: 1, maxE: 1, maxS: 1, want: "worldX"},
		{minE: 1, minS: 1, maxE: 2, maxS: 1, want: "worldX, helloY"},
		{minE: 2, minS: 2, maxE: 4, maxS: 1, want: "worldY, helloZ, worldZ"},
	} {
		min := data.EventID{Epoch: c.minE, Seq: c.minS}
		max := data.EventID{Epoch: c.maxE, Seq: c.maxS}

		got := joinEvents(evs.SelectMinMax(min, max))
		if c.want != got {
			t.Errorf("test case: %d want: '%s' got: '%s'", i, c.want, got)
		}
	}

}
