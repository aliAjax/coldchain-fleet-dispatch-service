package service

import (
	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

type AlertService struct {
	store *store.Store
	clock platform.Clock
}

func NewAlertService(st *store.Store, clk platform.Clock) *AlertService {
	return &AlertService{store: st, clock: clk}
}

func (a *AlertService) Raise(shipmentID, typ, message string) (domain.Alert, error) {
	if shipmentID == "" {
		return domain.Alert{}, platform.ErrInvalid
	}
	alert := domain.Alert{
		ID:         platform.NewID("alert"),
		ShipmentID: shipmentID,
		Type:       typ,
		Message:    message,
		CreatedAt:  a.clock.Now(),
	}
	if err := a.store.SaveAlert(alert); err != nil {
		return domain.Alert{}, err
	}
	return alert, nil
}

func (a *AlertService) Resolve(id string) (domain.Alert, error) {
	return a.store.ResolveAlert(id)
}

func (a *AlertService) ListUnresolved() []domain.Alert {
	out := make([]domain.Alert, 0)
	for _, alert := range a.store.ListAlerts() {
		if !alert.Resolved {
			out = append(out, alert)
		}
	}
	return out
}
