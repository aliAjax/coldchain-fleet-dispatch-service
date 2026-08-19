package worker

import (
	"context"
	"testing"
	"time"

	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/service"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

func TestDispatcherStopsOnCancel(t *testing.T) {
	st := store.New()
	dispatch := service.NewDispatchService(st, platform.NewSystemClock())
	d := NewDispatcher(st, dispatch, platform.NewSystemClock(), time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		d.Run(ctx)
		close(done)
	}()
	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("dispatcher did not stop after cancellation")
	}
}
