package domain

import "time"

type Vehicle struct {
	ID         string
	Plate      string
	Zone       TempZone
	CapacityKg float64
	CapacityL  float64
	Status     VehicleStatus
	UpdatedAt  time.Time
}

func (v Vehicle) CanCarry(s Shipment) bool {
	return v.Zone == s.Zone() &&
		v.CapacityKg >= s.WeightKg &&
		v.CapacityL >= s.VolumeL
}
