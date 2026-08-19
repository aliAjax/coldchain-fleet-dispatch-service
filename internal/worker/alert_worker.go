package worker

import (
	"context"
	"time"

	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/service"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

type AlertWorker struct {
	store    *store.Store
	alerts   *service.AlertService
	clock    platform.Clock
	window   time.Duration
	interval time.Duration
}

func NewAlertWorker(st *store.Store, alerts *service.AlertService, clk platform.Clock, window, interval time.Duration) *AlertWorker {
	return &AlertWorker{store: st, alerts: alerts, clock: clk, window: window, interval: interval}
}

func (w *AlertWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = w.Reconcile()
		}
	}
}

func (w *AlertWorker) Reconcile() error {
	cut := w.clock.Now().Add(-w.window)
	for _, sh := range w.store.ListShipments() {
		for _, reading := range w.store.ListReadings(sh.ID) {
			if reading.RecordedAt.Before(cut) {
				continue
			}
			if sh.AcceptableTemperature(reading.Celsius) {
				continue
			}
			if w.hasUnresolvedExcursion(sh.ID) {
				continue
			}
			if _, err := w.alerts.Raise(sh.ID, "temperature_excursion", "temperature outside shipment range"); err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *AlertWorker) hasUnresolvedExcursion(shipmentID string) bool {
	for _, alert := range w.alerts.ListUnresolved() {
		if alert.ShipmentID == shipmentID && alert.Type == "temperature_excursion" {
			return true
		}
	}
	return false
}
