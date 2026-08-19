package service

import (
	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

type QuoteService struct {
	store *store.Store
	clock platform.Clock
}

func NewQuoteService(st *store.Store, clk platform.Clock) *QuoteService {
	return &QuoteService{store: st, clock: clk}
}

func (q *QuoteService) Quote(shipmentID string, distanceKm float64) (domain.Quote, error) {
	sh, err := q.store.GetShipment(shipmentID)
	if err != nil {
		return domain.Quote{}, err
	}
	base := int64(distanceKm * 180)
	surcharge := tempSurcharge(sh.Zone())
	return domain.Quote{
		ID:            platform.NewID("quote"),
		ShipmentID:    shipmentID,
		DistanceKm:    distanceKm,
		BaseCents:     base,
		TempSurcharge: surcharge,
		TotalCents:    base + surcharge,
		CreatedAt:     q.clock.Now(),
	}, nil
}

func tempSurcharge(zone domain.TempZone) int64 {
	switch zone {
	case domain.ZoneFrozen:
		return 1200
	case domain.ZoneChilled:
		return 600
	default:
		return 0
	}
}
