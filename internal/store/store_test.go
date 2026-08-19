package store

import (
	"fmt"
	"sync"
	"testing"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
)

func TestStoreListUpdateNoRace(t *testing.T) {
	st := New()
	_ = st.SaveShipment(domain.Shipment{ID: "ship-1", Status: domain.ShipmentPending})
	_ = st.SaveAlert(domain.Alert{ID: "alert-1", ShipmentID: "ship-1"})

	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			if i%2 == 0 {
				for j := 0; j < 2000; j++ {
					_ = st.ListShipments()
					_ = st.ListAlerts()
				}
			} else {
				for j := 0; j < 2000; j++ {
					_, _ = st.UpdateShipmentStatus("ship-1", domain.ShipmentPending)
					_ = st.SaveAlert(domain.Alert{ID: fmt.Sprintf("alert-%d", j), ShipmentID: "ship-1"})
				}
			}
		}(i)
	}
	close(start)
	wg.Wait()
}

func TestSaveQuoteRoundtrip(t *testing.T) {
	st := New()
	quote := domain.Quote{ID: "quote-1", ShipmentID: "ship-1", TotalCents: 1200}
	if err := st.SaveQuote(quote); err != nil {
		t.Fatalf("save quote: %v", err)
	}
	got, err := st.GetQuote("quote-1")
	if err != nil {
		t.Fatalf("get quote: %v", err)
	}
	if got.ID != "quote-1" || got.TotalCents != 1200 {
		t.Fatalf("unexpected quote: %+v", got)
	}
}
