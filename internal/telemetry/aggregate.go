package telemetry

import "github.com/example/coldchain-fleet-dispatch-service/internal/domain"

type Stats struct {
	Count       int
	Min         float64
	Max         float64
	Average     float64
	Excursions  int
	SensorCount int
}

func Aggregate(readings []domain.TemperatureReading, minC, maxC float64) Stats {
	if len(readings) == 0 {
		return Stats{}
	}
	stats := Stats{
		Count: len(readings),
		Min:   readings[0].Celsius,
		Max:   readings[0].Celsius,
	}
	sensors := make(map[string]struct{})
	var sum float64
	for _, r := range readings {
		if r.Celsius < stats.Min {
			stats.Min = r.Celsius
		}
		if r.Celsius > stats.Max {
			stats.Max = r.Celsius
		}
		sum += r.Celsius
		if r.Celsius < minC || r.Celsius > maxC {
			stats.Excursions++
		}
		sensors[r.SensorID] = struct{}{}
	}
	stats.Average = sum / float64(stats.Count)
	stats.SensorCount = len(sensors)
	return stats
}
