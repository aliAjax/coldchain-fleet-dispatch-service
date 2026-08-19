package domain

import "time"

type Assignment struct {
	ID         string
	ShipmentID string
	VehicleID  string
	DriverID   string
	Status     AssignmentStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
