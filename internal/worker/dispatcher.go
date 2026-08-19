package worker

import (
	"context"
	"time"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/service"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

type Dispatcher struct {
	store    *store.Store
	dispatch *service.DispatchService
	clock    platform.Clock
	interval time.Duration
}

func NewDispatcher(st *store.Store, svc *service.DispatchService, clk platform.Clock, interval time.Duration) *Dispatcher {
	return &Dispatcher{store: st, dispatch: svc, clock: clk, interval: interval}
}

func (d *Dispatcher) Run(ctx context.Context) {
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = d.DispatchPending()
		}
	}
}

func (d *Dispatcher) DispatchPending() error {
	var ids []string
	for _, sh := range d.store.ListShipments() {
		if sh.Status == domain.ShipmentPending {
			ids = append(ids, sh.ID)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	_, err := d.dispatch.AssignBatch(ids)
	return err
}
