package domain

import "time"

type Shipment struct {
	ID          string
	Origin      string
	Destination string
	MinTempC    float64
	MaxTempC    float64
	WeightKg    float64
	VolumeL     float64
	Status      ShipmentStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s Shipment) Zone() TempZone {
	if s.MaxTempC <= -18 {
		return ZoneFrozen
	}
	if s.MaxTempC <= 8 {
		return ZoneChilled
	}
	return ZoneAmbient
}

func (s Shipment) AcceptableTemperature(celsius float64) bool {
	return celsius >= s.MinTempC && celsius <= s.MaxTempC
}
