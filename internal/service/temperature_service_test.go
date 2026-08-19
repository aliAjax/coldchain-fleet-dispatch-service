package service

import (
	"errors"
	"testing"
	"time"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

func TestRecordInRangeDoesNotRaiseAlert(t *testing.T) {
	st := store.New()
	st.SeedFleet(nil, nil)
	sh := domain.Shipment{ID: "ship-1", MinTempC: -30, MaxTempC: -18, Status: domain.ShipmentAssigned}
	_ = st.SaveShipment(sh)
	svc := NewTemperatureService(st, platform.FrozenClock{At: time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)})

	_, inRange, err := svc.Record("ship-1", -20, "sensor-a")
	if err != nil {
		t.Fatalf("record: %v", err)
	}
	if !inRange {
		t.Fatalf("expected in range")
	}
	if got := len(st.ListAlerts()); got != 0 {
		t.Fatalf("expected no alerts, got %d", got)
	}
}

func TestRecordOutOfRangeRaisesAlert(t *testing.T) {
	st := store.New()
	sh := domain.Shipment{ID: "ship-1", MinTempC: -30, MaxTempC: -18, Status: domain.ShipmentInTransit}
	_ = st.SaveShipment(sh)
	svc := NewTemperatureService(st, platform.FrozenClock{At: time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)})

	_, inRange, err := svc.Record("ship-1", -5, "sensor-a")
	if err != nil {
		t.Fatalf("record: %v", err)
	}
	if inRange {
		t.Fatalf("expected out of range")
	}
	if got := len(st.ListAlerts()); got != 1 {
		t.Fatalf("expected 1 alert, got %d", got)
	}
}

func TestAssignmentStatusNotFoundIsWrapped(t *testing.T) {
	st := store.New()
	_ = st.SaveShipment(domain.Shipment{ID: "ship-1", MinTempC: -30, MaxTempC: -18, Status: domain.ShipmentAssigned})
	svc := NewTemperatureService(st, platform.FrozenClock{At: time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)})

	_, err := svc.AssignmentStatus("ship-1")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, platform.ErrNotAssigned) {
		t.Fatalf("expected not assigned chain, got %v", err)
	}
}
