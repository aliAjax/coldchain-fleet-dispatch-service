package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/validation"
)

type createShipmentRequest struct {
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	MinTempC    float64 `json:"min_temp_c"`
	MaxTempC    float64 `json:"max_temp_c"`
	WeightKg    float64 `json:"weight_kg"`
	VolumeL     float64 `json:"volume_l"`
}

type recordReadingRequest struct {
	Celsius  float64 `json:"celsius"`
	SensorID string  `json:"sensor_id"`
}

func (a *API) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) handleCreateShipment(w http.ResponseWriter, r *http.Request) {
	var req createShipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	in := validation.ShipmentInput{
		Origin:      req.Origin,
		Destination: req.Destination,
		MinTempC:    req.MinTempC,
		MaxTempC:    req.MaxTempC,
		WeightKg:    req.WeightKg,
		VolumeL:     req.VolumeL,
	}
	if err := validation.ValidateShipment(in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	sh := domain.Shipment{
		ID:          platform.NewID("ship"),
		Origin:      req.Origin,
		Destination: req.Destination,
		MinTempC:    req.MinTempC,
		MaxTempC:    req.MaxTempC,
		WeightKg:    req.WeightKg,
		VolumeL:     req.VolumeL,
		Status:      domain.ShipmentPending,
		CreatedAt:   a.clock.Now(),
		UpdatedAt:   a.clock.Now(),
	}
	if err := a.store.SaveShipment(sh); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, sh)
}

func (a *API) handleGetShipment(w http.ResponseWriter, r *http.Request) {
	sh, err := a.store.GetShipment(r.PathValue("id"))
	if err != nil {
		writeAPIError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sh)
}

func (a *API) handleDispatch(w http.ResponseWriter, r *http.Request) {
	assignment, err := a.dispatch.Assign(r.PathValue("id"))
	if err != nil {
		writeAPIError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, assignment)
}

func (a *API) handleRecordReading(w http.ResponseWriter, r *http.Request) {
	var req recordReadingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	reading, inRange, err := a.temperature.Record(r.PathValue("id"), req.Celsius, req.SensorID)
	if err != nil {
		writeAPIError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"reading": reading, "in_range": inRange})
}

func (a *API) handleStartTransit(w http.ResponseWriter, r *http.Request) {
	assignment, err := a.handover.StartTransit(r.PathValue("id"))
	if err != nil {
		writeAPIError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, assignment)
}

func (a *API) handleCompleteDelivery(w http.ResponseWriter, r *http.Request) {
	assignment, err := a.handover.CompleteDelivery(r.PathValue("id"))
	if err != nil {
		writeAPIError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, assignment)
}

func (a *API) handleListAlerts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.alerts.ListUnresolved())
}

func writeAPIError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, platform.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, platform.ErrNotAssigned):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, platform.ErrConflict):
		writeError(w, http.StatusConflict, err)
	case errors.Is(err, platform.ErrNoVehicle):
		writeError(w, http.StatusUnprocessableEntity, err)
	case errors.Is(err, platform.ErrInvalid):
		writeError(w, http.StatusBadRequest, err)
	default:
		writeError(w, http.StatusInternalServerError, err)
	}
}
