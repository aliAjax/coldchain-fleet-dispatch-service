package service

import (
	"testing"
	"time"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

func TestRoutePlanPreservesStopOrder(t *testing.T) {
	st := store.New()
	_ = st.SaveAssignment(domain.Assignment{ID: "asgn-1", ShipmentID: "ship-1"})
	svc := NewRouteService(st, platform.FrozenClock{At: time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)})

	start := time.Date(2026, 8, 19, 8, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 19, 10, 0, 0, 0, time.UTC)
	stops := []domain.RouteStop{
		{Seq: 2, Kind: domain.StopDelivery, LocationID: "loc-b", WindowStart: start, WindowEnd: end},
		{Seq: 1, Kind: domain.StopPickup, LocationID: "loc-a", WindowStart: start, WindowEnd: end},
	}
	route, err := svc.Plan("asgn-1", stops, 50)
	if err != nil {
		t.Fatalf("plan route: %v", err)
	}
	sorted := route.SortedStops()
	if sorted[0].Seq != 1 || sorted[1].Seq != 2 {
		t.Fatalf("expected sorted stops, got %+v", sorted)
	}
	if route.PlannedMinutes <= 0 {
		t.Fatalf("expected positive planned minutes, got %d", route.PlannedMinutes)
	}
}
