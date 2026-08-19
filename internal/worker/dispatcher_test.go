package worker

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
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

func TestDispatchPendingAssignsAllShards(t *testing.T) {
	st := store.New()
	st.SeedFleet([]domain.Vehicle{
		{ID: "veh-frozen-1", Zone: domain.ZoneFrozen, CapacityKg: 12000, CapacityL: 32000, Status: domain.VehicleIdle},
		{ID: "veh-frozen-2", Zone: domain.ZoneFrozen, CapacityKg: 8000, CapacityL: 22000, Status: domain.VehicleIdle},
		{ID: "veh-chilled-1", Zone: domain.ZoneChilled, CapacityKg: 10000, CapacityL: 26000, Status: domain.VehicleIdle},
	}, []domain.Driver{
		{ID: "drv-frozen-1", Certification: domain.ZoneFrozen, Available: true},
		{ID: "drv-frozen-2", Certification: domain.ZoneFrozen, Available: true},
		{ID: "drv-chilled-1", Certification: domain.ZoneChilled, Available: true},
	})
	for _, spec := range []struct {
		id   string
		maxC float64
	}{{"ship-1", -20}, {"ship-2", -20}, {"ship-3", 2}} {
		_ = st.SaveShipment(domain.Shipment{
			ID: spec.id, Origin: "warehouse-a", Destination: "store-b",
			MinTempC: -30, MaxTempC: spec.maxC, WeightKg: 500, VolumeL: 900,
			Status: domain.ShipmentPending,
		})
	}

	clk := platform.FrozenClock{At: time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)}
	d := NewDispatcher(st, service.NewDispatchService(st, clk), clk, time.Millisecond)

	var wg sync.WaitGroup
	start := make(chan struct{})
	ids := []string{"ship-1", "ship-2", "ship-3"}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			<-start
			_ = d.DispatchPending()
		}(ids[i])
	}
	close(start)
	wg.Wait()

	for _, id := range []string{"ship-1", "ship-2", "ship-3"} {
		sh, err := st.GetShipment(id)
		if err != nil || sh.Status != domain.ShipmentAssigned {
			t.Fatalf("shipment %s not assigned: %v", id, err)
		}
	}
}
