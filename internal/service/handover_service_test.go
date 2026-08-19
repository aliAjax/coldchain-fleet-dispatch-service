package service

import (
	"testing"
	"time"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

func TestCompleteDeliveryFreesVehicle(t *testing.T) {
	st := store.New()
	st.SeedFleet([]domain.Vehicle{{ID: "veh-1", Plate: "X", Zone: domain.ZoneFrozen, CapacityKg: 1000, CapacityL: 1000, Status: domain.VehicleBusy}}, nil)
	_ = st.SaveShipment(domain.Shipment{ID: "ship-1", MinTempC: -30, MaxTempC: -18, Status: domain.ShipmentInTransit})
	_ = st.SaveAssignment(domain.Assignment{ID: "asgn-1", ShipmentID: "ship-1", VehicleID: "veh-1", DriverID: "drv-1", Status: domain.AssignmentInTransit, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	svc := NewHandoverService(st, platform.FrozenClock{At: time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)})

	a, err := svc.CompleteDelivery("ship-1")
	if err != nil {
		t.Fatalf("complete delivery: %v", err)
	}
	if a.Status != domain.AssignmentDelivered {
		t.Fatalf("expected delivered, got %s", a.Status)
	}
	v, _ := st.GetVehicle("veh-1")
	if v.Status != domain.VehicleIdle {
		t.Fatalf("expected vehicle idle, got %s", v.Status)
	}
}

func TestFinalizeDeliveredBatch(t *testing.T) {
	st := store.New()
	st.SeedFleet([]domain.Vehicle{{ID: "veh-1", Plate: "X", Zone: domain.ZoneFrozen, CapacityKg: 1000, CapacityL: 1000, Status: domain.VehicleBusy}}, nil)
	_ = st.SaveShipment(domain.Shipment{ID: "ship-1", MinTempC: -30, MaxTempC: -18, Status: domain.ShipmentDelivered})
	_ = st.SaveAssignment(domain.Assignment{ID: "asgn-1", ShipmentID: "ship-1", VehicleID: "veh-1", DriverID: "drv-1", Status: domain.AssignmentDelivered, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	svc := NewHandoverService(st, platform.FrozenClock{At: time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)})

	released, err := svc.FinalizeDeliveredBatch()
	if err != nil {
		t.Fatalf("finalize batch: %v", err)
	}
	if released != 1 {
		t.Fatalf("expected 1 released, got %d", released)
	}
	v, _ := st.GetVehicle("veh-1")
	if v.Status != domain.VehicleIdle {
		t.Fatalf("expected vehicle idle, got %s", v.Status)
	}
}
