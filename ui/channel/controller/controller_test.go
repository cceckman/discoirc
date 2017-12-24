package controller_test

import (
	"context"

	_ "github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel/controller"
	"github.com/cceckman/discoirc/ui/channel/mocks"
	_ "github.com/marcusolsson/tui-go"

	"testing"
)

func TestController_Resize(t *testing.T) {
	ui := mocks.NewUpdateCounter()

	m := &mocks.Model{}
	v := &mocks.View{}
	_ = controller.New(context.Background(), ui, v, m)

	ui.Add(2)
	// Resize when no events available.
	v.Controller.Resize(10)

	ui.RunSync(func() {
		if len(v.Events) > 0 {
			t.Errorf("wrong number of events: got: %v want: none", v.Events)
		}
	})
}
