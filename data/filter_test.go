package data_test

import (
	"testing"

	"github.com/cceckman/discoirc/data"
)

var anyCases = []data.Scope{
	{},
	{
		Net:  "foonet",
		Name: "#charnet",
	},
	{
		Net: "barnet",
	},
}

func TestFilter_Any(t *testing.T) {
	t.Parallel()
	filter := data.Filter{}

	for i, tt := range anyCases {
		got := filter.Match(tt)
		if !got {
			t.Errorf("error in case %d: got: %v want: %v", i, got, true)
		}
	}
}

var netCases = []struct {
	scope data.Scope
	want  bool
}{
	{
		scope: data.Scope{},
		want:  false,
	},
	{
		scope: data.Scope{
			Net:  "foonet",
			Name: "#charnet",
		},
		want: true,
	},
	{
		scope: data.Scope{
			Net:  "foonet",
			Name: "",
		},
		want: true,
	},
	{
		scope: data.Scope{
			Net:  "barnet",
			Name: "",
		},
		want: false,
	},
}

func TestFilter_Net(t *testing.T) {
	t.Parallel()
	filter := data.Filter{
		Scope: data.Scope{
			Net:  "foonet",
			Name: "anything!",
		},
		MatchNet: true,
	}

	for i, tt := range netCases {
		got := filter.Match(tt.scope)
		if got != tt.want {
			t.Errorf("error in case %d: got: %v want: %v", i, got, tt.want)
		}
	}
}

var chanCases = []struct {
	scope data.Scope
	want  bool
}{
	{
		scope: data.Scope{},
		want:  false,
	},
	{
		scope: data.Scope{
			Net:  "foonet",
			Name: "#charnet",
		},
		want: true,
	},
	{
		scope: data.Scope{
			Net:  "foonet",
			Name: "",
		},
		want: false,
	},
	{
		scope: data.Scope{
			Net:  "foonet",
			Name: "#bulbanet",
		},
		want: false,
	},
	{
		scope: data.Scope{
			Net:  "barnet",
			Name: "#bulbanet",
		},
		want: false,
	},
}

func TestFilter_Chan(t *testing.T) {
	t.Parallel()
	filter := data.Filter{
		Scope: data.Scope{
			Net:  "foonet",
			Name: "#charnet",
		},
		MatchNet:  true,
		MatchName: true,
	}

	for i, tt := range chanCases {
		got := filter.Match(tt.scope)
		if got != tt.want {
			t.Errorf("error in case %d: got: %v want: %v", i, got, tt.want)
		}
	}
}
