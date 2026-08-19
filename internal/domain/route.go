package domain

import (
	"errors"
	"sort"
	"time"
)

type StopKind string

const (
	StopPickup  StopKind = "pickup"
	StopDelivery StopKind = "delivery"
	StopRest     StopKind = "rest"
)

type RouteStop struct {
	Seq         int
	Kind        StopKind
	LocationID  string
	WindowStart time.Time
	WindowEnd   time.Time
}

type Route struct {
	ID             string
	AssignmentID   string
	Stops          []RouteStop
	EstimatedKm    float64
	PlannedMinutes int
}

func (r Route) Validate() error {
	if r.AssignmentID == "" {
		return errors.New("route assignment is required")
	}
	if len(r.Stops) < 2 {
		return errors.New("route needs at least two stops")
	}
	for _, stop := range r.Stops {
		if stop.Kind == "" || stop.LocationID == "" {
			return errors.New("route stop is incomplete")
		}
		if !stop.WindowEnd.After(stop.WindowStart) {
			return errors.New("route stop window is invalid")
		}
	}
	return nil
}

func (r Route) SortedStops() []RouteStop {
	out := make([]RouteStop, len(r.Stops))
	copy(out, r.Stops)
	sort.Slice(out, func(i, j int) bool { return out[i].Seq < out[j].Seq })
	return out
}
