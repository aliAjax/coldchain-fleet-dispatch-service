package validation

import "errors"

var errInvalidShipment = errors.New("invalid shipment")

func IsInvalidShipment(err error) bool {
	return errors.Is(err, errInvalidShipment)
}
