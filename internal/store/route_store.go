package store

import (
	"sort"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
)

func (s *Store) SaveRoute(route domain.Route) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if route.ID == "" {
		return platform.ErrInvalid
	}
	s.routes[route.ID] = route
	return nil
}

func (s *Store) GetRoute(id string) (domain.Route, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	route, ok := s.routes[id]
	if !ok {
		return domain.Route{}, platform.ErrNotFound
	}
	return route, nil
}

func (s *Store) ListRoutes() []domain.Route {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Route, 0, len(s.routes))
	for _, route := range s.routes {
		out = append(out, route)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
