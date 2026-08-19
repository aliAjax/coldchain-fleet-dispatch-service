package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/service"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

func TestChecklistNotFoundReturns404(t *testing.T) {
	st := store.New()
	clk := platform.FrozenClock{At: time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)}
	dispatch := service.NewDispatchService(st, clk)
	temperature := service.NewTemperatureService(st, clk)
	handover := service.NewHandoverService(st, clk)
	alerts := service.NewAlertService(st, clk)
	quotes := service.NewQuoteService(st, clk)
	routes := service.NewRouteService(st, clk)
	checklists := service.NewChecklistService(st, clk)
	api := NewAPI(st, dispatch, temperature, handover, alerts, quotes, routes, checklists, clk)

	req := httptest.NewRequest(http.MethodPost, "/v1/checklists/missing/items", strings.NewReader(`{"key":"seal","done":true}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	api.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
