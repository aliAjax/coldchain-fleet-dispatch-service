package service

import (
	"time"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

type RouteService struct {
	store        *store.Store
	clock        platform.Clock
	speedKmh     float64
	dwellMinutes int
}

func NewRouteService(st *store.Store, clk platform.Clock) *RouteService {
	return &RouteService{store: st, clock: clk, speedKmh: 55, dwellMinutes: 20}
}

func (s *RouteService) Plan(assignmentID string, stops []domain.RouteStop, estimatedKm float64) (domain.Route, error) {
	route := domain.Route{
		ID:             platform.NewID("route"),
		AssignmentID:   assignmentID,
		Stops:          stops,
		EstimatedKm:    estimatedKm,
		PlannedMinutes: s.estimateMinutes(stops, estimatedKm),
	}
	if err := route.Validate(); err != nil {
		return domain.Route{}, err
	}
	if err := s.store.SaveRoute(route); err != nil {
		return domain.Route{}, err
	}
	return route, nil
}

func (s *RouteService) EstimateArrival(route domain.Route, start time.Time) (time.Time, error) {
	if err := route.Validate(); err != nil {
		return time.Time{}, err
	}
	minutes := s.estimateMinutes(route.SortedStops(), route.EstimatedKm)
	return start.Add(time.Duration(minutes) * time.Minute), nil
}

func (s *RouteService) estimateMinutes(stops []domain.RouteStop, km float64) int {
	if km <= 0 {
		km = 1
	}
	driveMinutes := int((km / s.speedKmh) * 60)
	return driveMinutes + len(stops)*s.dwellMinutes
}
