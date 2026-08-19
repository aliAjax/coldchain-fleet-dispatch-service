package service

import (
	"sync"
	"testing"
	"time"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

func newDispatchFixture(t *testing.T) (*store.Store, *DispatchService) {
	t.Helper()
	st := store.New()
	st.SeedFleet([]domain.Vehicle{
		{ID: "veh-frozen-1", Plate: "FROZEN-01", Zone: domain.ZoneFrozen, CapacityKg: 12000, CapacityL: 32000, Status: domain.VehicleIdle},
		{ID: "veh-frozen-2", Plate: "FROZEN-02", Zone: domain.ZoneFrozen, CapacityKg: 8000, CapacityL: 22000, Status: domain.VehicleIdle},
		{ID: "veh-chilled-1", Plate: "CHILLED-01", Zone: domain.ZoneChilled, CapacityKg: 10000, CapacityL: 26000, Status: domain.VehicleIdle},
	}, []domain.Driver{
		{ID: "drv-frozen-1", Name: "Ada Frost", Certification: domain.ZoneFrozen, Available: true},
		{ID: "drv-frozen-2", Name: "Bo Cold", Certification: domain.ZoneFrozen, Available: true},
		{ID: "drv-chilled-1", Name: "Cy Chill", Certification: domain.ZoneChilled, Available: true},
	})
	return st, NewDispatchService(st, platform.FrozenClock{At: time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)})
}

func seedShipment(t *testing.T, st *store.Store, id string, status domain.ShipmentStatus, maxC float64) {
	t.Helper()
	sh := domain.Shipment{
		ID:          id,
		Origin:      "warehouse-a",
		Destination: "store-b",
		MinTempC:    -30,
		MaxTempC:    maxC,
		WeightKg:    500,
		VolumeL:     900,
		Status:      status,
		CreatedAt:   time.Date(2026, 8, 19, 11, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2026, 8, 19, 11, 0, 0, 0, time.UTC),
	}
	if err := st.SaveShipment(sh); err != nil {
		t.Fatal(err)
	}
}

func TestAssignCompatibleVehicle(t *testing.T) {
	st, svc := newDispatchFixture(t)
	seedShipment(t, st, "ship-1", domain.ShipmentPending, -20)

	a, err := svc.Assign("ship-1")
	if err != nil {
		t.Fatalf("assign: %v", err)
	}
	if a.Status != domain.AssignmentAssigned {
		t.Fatalf("expected assigned, got %s", a.Status)
	}
	sh, _ := st.GetShipment("ship-1")
	if sh.Status != domain.ShipmentAssigned {
		t.Fatalf("expected shipment assigned, got %s", sh.Status)
	}
}

func TestAssignBatchWaitsForAllShards(t *testing.T) {
	st, svc := newDispatchFixture(t)
	seedShipment(t, st, "ship-1", domain.ShipmentPending, -20)
	seedShipment(t, st, "ship-2", domain.ShipmentPending, -20)
	seedShipment(t, st, "ship-3", domain.ShipmentPending, 2)

	ids := []string{"ship-1", "ship-2", "ship-3"}
	assignments, err := svc.AssignBatch(ids)
	if err != nil {
		t.Fatalf("assign batch: %v", err)
	}
	if len(assignments) != len(ids) {
		t.Fatalf("expected %d assignments, got %d", len(ids), len(assignments))
	}

	var wg sync.WaitGroup
	start := make(chan struct{})
	for _, id := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			<-start
			sh, err := st.GetShipment(id)
			if err != nil || sh.Status != domain.ShipmentAssigned {
				t.Errorf("shipment %s not assigned: %v", id, err)
			}
		}(id)
	}
	close(start)
	wg.Wait()
}
