package domain

import "time"

type Alert struct {
	ID         string
	ShipmentID string
	Type       string
	Message    string
	CreatedAt  time.Time
	Resolved   bool
	ResolvedAt time.Time
}
