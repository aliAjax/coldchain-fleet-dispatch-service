package service

import (
	"errors"
	"testing"
	"time"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

func TestChecklistNotFoundIsWrapped(t *testing.T) {
	st := store.New()
	_ = st.SaveAssignment(domain.Assignment{ID: "asgn-1", ShipmentID: "ship-1", Status: domain.AssignmentAssigned})
	svc := NewChecklistService(st, platform.FrozenClock{At: time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)})

	_, err := svc.Mark("asgn-missing", "seal", true)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, platform.ErrNotFound) {
		t.Fatalf("expected not found chain, got %v", err)
	}
}
