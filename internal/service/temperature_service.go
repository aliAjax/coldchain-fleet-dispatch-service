package service

import (
	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

type TemperatureService struct {
	store *store.Store
	clock platform.Clock
}

func NewTemperatureService(st *store.Store, clk platform.Clock) *TemperatureService {
	return &TemperatureService{store: st, clock: clk}
}

func (t *TemperatureService) Record(shipmentID string, celsius float64, sensorID string) (domain.TemperatureReading, bool, error) {
	sh, err := t.store.GetShipment(shipmentID)
	if err != nil {
		return domain.TemperatureReading{}, false, err
	}
	if sh.Status != domain.ShipmentAssigned && sh.Status != domain.ShipmentInTransit {
		return domain.TemperatureReading{}, false, platform.ErrConflict
	}

	reading := domain.TemperatureReading{
		ID:         platform.NewID("reading"),
		ShipmentID: shipmentID,
		RecordedAt: t.clock.Now(),
		Celsius:    celsius,
		SensorID:   sensorID,
	}
	if err := t.store.SaveReading(reading); err != nil {
		return domain.TemperatureReading{}, false, err
	}

	inRange := sh.AcceptableTemperature(celsius)
	if !inRange {
		alert := domain.Alert{
			ID:         platform.NewID("alert"),
			ShipmentID: shipmentID,
			Type:       "temperature_excursion",
			Message:    "temperature outside shipment range",
			CreatedAt:  t.clock.Now(),
		}
		if err := t.store.SaveAlert(alert); err != nil {
			return domain.TemperatureReading{}, false, err
		}
	}
	return reading, inRange, nil
}

func (t *TemperatureService) AssignmentStatus(shipmentID string) (domain.AssignmentStatus, error) {
	for _, a := range t.store.ListAssignments() {
		if a.ShipmentID == shipmentID {
			return a.Status, nil
		}
	}
	return "", platform.ErrNotAssigned
}
