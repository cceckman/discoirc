package controller_test

import (
	"context"
	"time"

	_ "github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/cceckman/discoirc/ui/channel/controller"
	"github.com/cceckman/discoirc/ui/channel/mocks"
	_ "github.com/marcusolsson/tui-go"

	"testing"
)

// Updater exists to fulfil the controller.UIUpdater interface in a minimal manner.
type Updater struct{}

func (_ Updater) Update(f func()) {
	f()
}

func setup(ctx context.Context) (*mocks.Model, *mocks.View, channel.Controller) {
	m := mocks.NewModel()
	v := &mocks.View{}
	c := controller.New(ctx, Updater(struct{}{}), v, m)

	return m, v, c

}

func TestController_Resize(t *testing.T) {
	_, v, _ := setup(context.Background())

	// Resize when no events available; expect view to still have no events.
	v.Controller.Resize(10)
	// Allow other threads to schedule.
	// TODO: Way to implement this without a race in the assertion?
	// testing.AwaitAllOtherRoutinesToBeBlocked?
	time.Sleep(1 * time.Millisecond)
	if len(v.Events) > 0 {
		t.Errorf("wrong number of events: got: %v want: none", v.Events)
	}
}
