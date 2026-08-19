package store

import (
	"sort"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
)

func (s *Store) SaveAssignment(a domain.Assignment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if a.ID == "" || a.ShipmentID == "" {
		return platform.ErrInvalid
	}
	s.assignments[a.ID] = a
	return nil
}

func (s *Store) GetAssignment(id string) (domain.Assignment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.assignments[id]
	if !ok {
		return domain.Assignment{}, platform.ErrNotFound
	}
	return a, nil
}

func (s *Store) ListAssignments() []domain.Assignment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Assignment, 0, len(s.assignments))
	for _, a := range s.assignments {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (s *Store) UpdateAssignmentStatus(id string, status domain.AssignmentStatus) (domain.Assignment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.assignments[id]
	if !ok {
		return domain.Assignment{}, platform.ErrNotFound
	}
	a.Status = status
	s.assignments[id] = a
	return a, nil
}

func (s *Store) GetVehicle(id string) (domain.Vehicle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.vehicles[id]
	if !ok {
		return domain.Vehicle{}, platform.ErrNotFound
	}
	return v, nil
}

func (s *Store) ListVehicles() []domain.Vehicle {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Vehicle, 0, len(s.vehicles))
	for _, v := range s.vehicles {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (s *Store) UpdateVehicleStatus(id string, status domain.VehicleStatus) (domain.Vehicle, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.vehicles[id]
	if !ok {
		return domain.Vehicle{}, platform.ErrNotFound
	}
	v.Status = status
	s.vehicles[id] = v
	return v, nil
}

func (s *Store) GetDriver(id string) (domain.Driver, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.drivers[id]
	if !ok {
		return domain.Driver{}, platform.ErrNotFound
	}
	return d, nil
}

func (s *Store) ListDrivers() []domain.Driver {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Driver, 0, len(s.drivers))
	for _, d := range s.drivers {
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
