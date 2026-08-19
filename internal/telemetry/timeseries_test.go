package telemetry

import (
	"testing"
	"time"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
)

func TestMovingAverageKeepsWindow(t *testing.T) {
	readings := []domain.TemperatureReading{
		{Celsius: 1, RecordedAt: time.Unix(1, 0)},
		{Celsius: 2, RecordedAt: time.Unix(2, 0)},
		{Celsius: 3, RecordedAt: time.Unix(3, 0)},
	}
	got := MovingAverage(readings, 2)
	want := []float64{1, 1.5, 2.5}
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: expected %v, got %v", i, want[i], got[i])
		}
	}
}

func TestExcursionSpansGroupsRuns(t *testing.T) {
	readings := []domain.TemperatureReading{
		{Celsius: 0}, {Celsius: 5}, {Celsius: -10}, {Celsius: 6},
	}
	spans := ExcursionSpans(readings, -18, -1)
	if len(spans) != 2 {
		t.Fatalf("expected 2 spans, got %d", len(spans))
	}
	if spans[0].Readings != 2 || spans[1].Readings != 1 {
		t.Fatalf("unexpected spans: %+v", spans)
	}
}
