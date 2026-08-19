package telemetry

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
)

// ParseReadings reads CSV lines in the form sensorID,celsius and returns
// readings keyed by shipment ID.
func ParseReadings(shipmentID string, at time.Time, r io.Reader) ([]domain.TemperatureReading, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = 2
	cr.TrimLeadingSpace = true
	var out []domain.TemperatureReading
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse readings: %w", err)
		}
		sensorID := strings.TrimSpace(rec[0])
		celsius, err := strconv.ParseFloat(strings.TrimSpace(rec[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("parse celsius for sensor %s: %w", sensorID, err)
		}
		out = append(out, domain.TemperatureReading{
			ID:         platform.NewID("reading"),
			ShipmentID: shipmentID,
			RecordedAt: at,
			Celsius:    celsius,
			SensorID:   sensorID,
		})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no readings parsed")
	}
	return out, nil
}
