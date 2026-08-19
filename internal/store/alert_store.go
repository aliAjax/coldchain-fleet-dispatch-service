package store

import (
	"sort"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
)

func (s *Store) SaveReading(r domain.TemperatureReading) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r.ID == "" || r.ShipmentID == "" {
		return platform.ErrInvalid
	}
	s.readings[r.ShipmentID] = append(s.readings[r.ShipmentID], r)
	return nil
}

func (s *Store) ListReadings(shipmentID string) []domain.TemperatureReading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.TemperatureReading, 0, len(s.readings[shipmentID]))
	for _, r := range s.readings[shipmentID] {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].RecordedAt.Before(out[j].RecordedAt) })
	return out
}

func (s *Store) SaveAlert(a domain.Alert) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if a.ID == "" || a.ShipmentID == "" {
		return platform.ErrInvalid
	}
	s.alerts[a.ID] = a
	return nil
}

func (s *Store) ListAlerts() []domain.Alert {
	snapshot := make([]domain.Alert, 0, len(s.alerts))
	for _, a := range s.alerts {
		snapshot = append(snapshot, a)
	}
	return snapshot
}

func (s *Store) ResolveAlert(id string) (domain.Alert, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.alerts[id]
	if !ok {
		return domain.Alert{}, platform.ErrNotFound
	}
	a.Resolved = true
	s.alerts[id] = a
	return a, nil
}
