package api

import (
	"net/http"

	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/service"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

type API struct {
	store       *store.Store
	dispatch    *service.DispatchService
	temperature *service.TemperatureService
	handover    *service.HandoverService
	alerts      *service.AlertService
	quotes      *service.QuoteService
	routes      *service.RouteService
	checklists  *service.ChecklistService
	clock       platform.Clock
}

func NewAPI(st *store.Store, dispatch *service.DispatchService, temperature *service.TemperatureService, handover *service.HandoverService, alerts *service.AlertService, quotes *service.QuoteService, routes *service.RouteService, checklists *service.ChecklistService, clk platform.Clock) *API {
	return &API{store: st, dispatch: dispatch, temperature: temperature, handover: handover, alerts: alerts, quotes: quotes, routes: routes, checklists: checklists, clock: clk}
}

func (a *API) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", a.handleHealth)
	mux.HandleFunc("POST /v1/shipments", a.handleCreateShipment)
	mux.HandleFunc("GET /v1/shipments/{id}", a.handleGetShipment)
	mux.HandleFunc("POST /v1/dispatch/{id}", a.handleDispatch)
	mux.HandleFunc("POST /v1/shipments/{id}/readings", a.handleRecordReading)
	mux.HandleFunc("POST /v1/shipments/{id}/start-transit", a.handleStartTransit)
	mux.HandleFunc("POST /v1/shipments/{id}/complete-delivery", a.handleCompleteDelivery)
	mux.HandleFunc("GET /v1/alerts", a.handleListAlerts)
	mux.HandleFunc("POST /v1/quotes", a.handleQuote)
	mux.HandleFunc("GET /v1/report", a.handleReport)
	mux.HandleFunc("POST /v1/routes", a.handlePlanRoute)
	mux.HandleFunc("POST /v1/checklists", a.handleBuildChecklist)
	mux.HandleFunc("POST /v1/checklists/{assignmentID}/items", a.handleMarkChecklist)
	return recoverer(logging(requestID(mux)))
}
