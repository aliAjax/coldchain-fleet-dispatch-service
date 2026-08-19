package api

import (
	"encoding/json"
	"net/http"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/report"
)

type quoteRequest struct {
	ShipmentID string  `json:"shipment_id"`
	DistanceKm float64 `json:"distance_km"`
}

type routeRequest struct {
	AssignmentID string             `json:"assignment_id"`
	EstimatedKm  float64            `json:"estimated_km"`
	Stops        []domain.RouteStop `json:"stops"`
}

type checklistRequest struct {
	AssignmentID string `json:"assignment_id"`
}

type checklistItemRequest struct {
	Key  string `json:"key"`
	Done bool   `json:"done"`
}

func (a *API) handleQuote(w http.ResponseWriter, r *http.Request) {
	var req quoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	quote, err := a.quotes.Quote(req.ShipmentID, req.DistanceKm)
	if err != nil {
		writeAPIError(w, err)
		return
	}
	if err := a.store.SaveQuote(quote); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, quote)
}

func (a *API) handleReport(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, report.Build(a.store))
}

func (a *API) handlePlanRoute(w http.ResponseWriter, r *http.Request) {
	var req routeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	route, err := a.routes.Plan(req.AssignmentID, req.Stops, req.EstimatedKm)
	if err != nil {
		writeAPIError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, route)
}

func (a *API) handleBuildChecklist(w http.ResponseWriter, r *http.Request) {
	var req checklistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	checklist, err := a.checklists.Build(req.AssignmentID)
	if err != nil {
		writeAPIError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, checklist)
}

func (a *API) handleMarkChecklist(w http.ResponseWriter, r *http.Request) {
	var req checklistItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	checklist, err := a.checklists.Mark(r.PathValue("assignmentID"), req.Key, req.Done)
	if err != nil {
		writeAPIError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, checklist)
}
