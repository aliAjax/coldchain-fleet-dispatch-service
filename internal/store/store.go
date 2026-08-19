package store

import (
	"sync"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
)

type Store struct {
	mu sync.RWMutex

	shipments   map[string]domain.Shipment
	vehicles    map[string]domain.Vehicle
	drivers     map[string]domain.Driver
	assignments map[string]domain.Assignment
	readings    map[string][]domain.TemperatureReading
	alerts      map[string]domain.Alert
	routes      map[string]domain.Route
	quotes      map[string]domain.Quote
	checklists  map[string]domain.HandoverChecklist
	events      []domain.AuditEvent
}

func New() *Store {
	return &Store{
		shipments:   make(map[string]domain.Shipment),
		vehicles:    make(map[string]domain.Vehicle),
		drivers:     make(map[string]domain.Driver),
		assignments: make(map[string]domain.Assignment),
		readings:    make(map[string][]domain.TemperatureReading),
		alerts:      make(map[string]domain.Alert),
		routes:      make(map[string]domain.Route),
		quotes:      make(map[string]domain.Quote),
		checklists:  make(map[string]domain.HandoverChecklist),
		events:      make([]domain.AuditEvent, 0),
	}
}

func (s *Store) SeedFleet(vehicles []domain.Vehicle, drivers []domain.Driver) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range vehicles {
		s.vehicles[v.ID] = v
	}
	for _, d := range drivers {
		s.drivers[d.ID] = d
	}
}
