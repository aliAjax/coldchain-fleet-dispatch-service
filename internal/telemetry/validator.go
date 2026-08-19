package telemetry

import (
	"fmt"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
)

func ValidateReading(r domain.TemperatureReading) error {
	if r.ShipmentID == "" {
		return fmt.Errorf("shipment id is required")
	}
	if r.SensorID == "" {
		return fmt.Errorf("sensor id is required")
	}
	if r.Celsius < -100 || r.Celsius > 100 {
		return fmt.Errorf("celsius value out of range: %v", r.Celsius)
	}
	return nil
}
