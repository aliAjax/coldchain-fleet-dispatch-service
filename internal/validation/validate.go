package validation

import (
	"fmt"
)

type ShipmentInput struct {
	Origin      string
	Destination string
	MinTempC    float64
	MaxTempC    float64
	WeightKg    float64
	VolumeL     float64
}

func ValidateShipment(in ShipmentInput) error {
	if in.Origin == "" || in.Destination == "" {
		return fmt.Errorf("%w: origin and destination are required", errInvalidShipment)
	}
	if in.MinTempC >= in.MaxTempC {
		return fmt.Errorf("%w: min temperature must be below max", errInvalidShipment)
	}
	if in.WeightKg <= 0 {
		return fmt.Errorf("%w: weight must be positive", errInvalidShipment)
	}
	if in.VolumeL <= 0 {
		return fmt.Errorf("%w: volume must be positive", errInvalidShipment)
	}
	return nil
}
