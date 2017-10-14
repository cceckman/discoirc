package model

import (
	"fmt"
	"testing"
)

// TestMockChannelGetMessages makes sure I've done my indexing right.
func TestMockChannelGetMessages(t *testing.T) {
	m := &MockChannel{
		messages: []string{
			"0 hello muddah",
			"1 hello faddah",
			"2 here I am at",
			"3 Camp Grenada",
		},
	}

	for i, cs := range []struct {
		Size, Offset uint
		Want         []string
	}{
		{0, 100, m.messages},
		{0, 1, []string{"3 Camp Grenada"}},
		{1, 1, []string{"2 here I am at"}},
		{4, 1, []string{}},
		{2, 3, []string{"0 hello muddah", "1 hello faddah"}},
	} {
		got := m.GetMessages(cs.Size, cs.Offset)
		err := fmt.Sprintf("unexpected result for test case %d: got: %v want: %v", i, got, cs.Want)
		if len(got) != len(cs.Want) {
			t.Errorf(err)
		}
		for i := range got {
			if got[i] != cs.Want[i] {
				t.Errorf(err)
			}
		}
	}
}
