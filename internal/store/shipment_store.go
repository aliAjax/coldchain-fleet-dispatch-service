package store

import (
	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
)

func (s *Store) SaveShipment(sh domain.Shipment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sh.ID == "" {
		return platform.ErrInvalid
	}
	s.shipments[sh.ID] = sh
	return nil
}

func (s *Store) GetShipment(id string) (domain.Shipment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sh, ok := s.shipments[id]
	if !ok {
		return domain.Shipment{}, platform.ErrNotFound
	}
	return sh, nil
}

func (s *Store) ListShipments() []domain.Shipment {
	snapshot := make([]domain.Shipment, 0, len(s.shipments))
	for _, sh := range s.shipments {
		snapshot = append(snapshot, sh)
	}
	return snapshot
}

func (s *Store) UpdateShipmentStatus(id string, status domain.ShipmentStatus) (domain.Shipment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sh, ok := s.shipments[id]
	if !ok {
		return domain.Shipment{}, platform.ErrNotFound
	}
	sh.Status = status
	s.shipments[id] = sh
	return sh, nil
}
