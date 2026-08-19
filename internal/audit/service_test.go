package audit

import (
	"errors"
	"testing"
	"time"

	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

func TestRecordEmptyEntityIsInvalid(t *testing.T) {
	st := store.New()
	svc := NewService(st, platform.FrozenClock{At: time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)})

	_, err := svc.Record("", "kind", "payload")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, platform.ErrInvalid) {
		t.Fatalf("expected invalid chain, got %v", err)
	}
}
