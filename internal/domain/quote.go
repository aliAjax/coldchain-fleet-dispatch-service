package domain

import "time"

type Quote struct {
	ID         string
	ShipmentID string
	DistanceKm float64
	BaseCents  int64
	TempSurcharge int64
	TotalCents int64
	CreatedAt  time.Time
}
