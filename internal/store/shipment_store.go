package store

import (
	"sort"

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
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Shipment, 0, len(s.shipments))
	for _, sh := range s.shipments {
		out = append(out, sh)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
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
