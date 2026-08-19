package telemetry

import "github.com/example/coldchain-fleet-dispatch-service/internal/domain"

type ExcursionSpan struct {
	StartC   int
	EndC     int
	Readings int
}

func MovingAverage(readings []domain.TemperatureReading, window int) []float64 {
	if window <= 0 || len(readings) == 0 {
		return []float64{}
	}
	out := make([]float64, 0, len(readings))
	for i := 0; i < len(readings); i++ {
		start := i - window + 1
		if start < 0 {
			start = 0
		}
		var sum float64
		for j := start; j <= i; j++ {
			sum += readings[j].Celsius
		}
		out = append(out, sum/float64(i-start+1))
	}
	return out
}

func ExcursionSpans(readings []domain.TemperatureReading, minC, maxC float64) []ExcursionSpan {
	var spans []ExcursionSpan
	for i := 0; i < len(readings); i++ {
		if readings[i].Celsius < minC || readings[i].Celsius > maxC {
			if len(spans) == 0 || spans[len(spans)-1].EndC != i-1 {
				spans = append(spans, ExcursionSpan{StartC: i, EndC: i})
			} else {
				spans[len(spans)-1].EndC = i
			}
		}
	}
	for i := range spans {
		spans[i].Readings = spans[i].EndC - spans[i].StartC + 1
	}
	return spans
}
