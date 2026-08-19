package platform

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrConflict    = errors.New("conflict")
	ErrInvalid     = errors.New("invalid input")
	ErrNoVehicle   = errors.New("no available vehicle")
	ErrNotAssigned = errors.New("shipment not assigned")
	ErrClosed      = errors.New("already closed")
)
