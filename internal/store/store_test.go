package store

import (
	"fmt"
	"sync"
	"testing"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
)

func TestConcurrentListAndUpdate(t *testing.T) {
	st := New()
	_ = st.SaveShipment(domain.Shipment{ID: "ship-1", Status: domain.ShipmentPending})
	_ = st.SaveAlert(domain.Alert{ID: "alert-1", ShipmentID: "ship-1"})

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 2000; i++ {
			_ = st.ListShipments()
			_ = st.ListAlerts()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 2000; i++ {
			_, _ = st.UpdateShipmentStatus("ship-1", domain.ShipmentPending)
			_ = st.SaveAlert(domain.Alert{ID: fmt.Sprintf("alert-%d", i), ShipmentID: "ship-1"})
		}
	}()
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
