package mocks

import (
	"github.com/marcusolsson/tui-go"
)

var _ tui.UI = &UI{}

// UI is a mock of tui.UI.
type UI struct{
	tui.UI
}

