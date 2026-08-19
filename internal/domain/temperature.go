package domain

import "time"

type TemperatureReading struct {
	ID         string
	ShipmentID string
	RecordedAt time.Time
	Celsius    float64
	SensorID   string
}

func (r TemperatureReading) Within(minC, maxC float64) bool {
	return r.Celsius >= minC && r.Celsius <= maxC
}
