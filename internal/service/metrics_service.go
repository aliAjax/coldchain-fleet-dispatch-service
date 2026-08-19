package service

import (
	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

type MetricsService struct {
	store *store.Store
}

func NewMetricsService(st *store.Store) *MetricsService {
	return &MetricsService{store: st}
}

type DispatchMetrics struct {
	TotalShipments   int
	PendingShipments int
	Delivered        int
	DeliveryRatePct  float64
}

func (m *MetricsService) ShipmentMetrics() DispatchMetrics {
	var metrics DispatchMetrics
	for _, sh := range m.store.ListShipments() {
		metrics.TotalShipments++
		switch sh.Status {
		case domain.ShipmentPending:
			metrics.PendingShipments++
		case domain.ShipmentDelivered:
			metrics.Delivered++
		}
	}
	if metrics.TotalShipments > 0 {
		metrics.DeliveryRatePct = float64(metrics.Delivered) / float64(metrics.TotalShipments) * 100
	}
	return metrics
}
