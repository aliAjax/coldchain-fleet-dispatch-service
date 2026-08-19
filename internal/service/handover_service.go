package service

import (
	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

type HandoverService struct {
	store *store.Store
	clock platform.Clock
}

func NewHandoverService(st *store.Store, clk platform.Clock) *HandoverService {
	return &HandoverService{store: st, clock: clk}
}

func (h *HandoverService) StartTransit(shipmentID string) (domain.Assignment, error) {
	a, err := h.findAssignment(shipmentID)
	if err != nil {
		return domain.Assignment{}, err
	}
	if a.Status != domain.AssignmentAssigned {
		return domain.Assignment{}, platform.ErrConflict
	}
	if _, err := h.store.UpdateAssignmentStatus(a.ID, domain.AssignmentInTransit); err != nil {
		return domain.Assignment{}, err
	}
	if _, err := h.store.UpdateShipmentStatus(shipmentID, domain.ShipmentInTransit); err != nil {
		return domain.Assignment{}, err
	}
	return h.store.GetAssignment(a.ID)
}

func (h *HandoverService) CompleteDelivery(shipmentID string) (domain.Assignment, error) {
	a, err := h.findAssignment(shipmentID)
	if err != nil {
		return domain.Assignment{}, err
	}
	if a.Status != domain.AssignmentInTransit {
		return domain.Assignment{}, platform.ErrConflict
	}
	if _, err := h.store.UpdateAssignmentStatus(a.ID, domain.AssignmentDelivered); err != nil {
		return domain.Assignment{}, err
	}
	if _, err := h.store.UpdateShipmentStatus(shipmentID, domain.ShipmentDelivered); err != nil {
		return domain.Assignment{}, err
	}
	if _, err := h.store.UpdateVehicleStatus(a.VehicleID, domain.VehicleIdle); err != nil {
		return domain.Assignment{}, err
	}
	return h.store.GetAssignment(a.ID)
}

func (h *HandoverService) findAssignment(shipmentID string) (domain.Assignment, error) {
	for _, a := range h.store.ListAssignments() {
		if a.ShipmentID == shipmentID {
			return a, nil
		}
	}
	return domain.Assignment{}, platform.ErrNotAssigned
}

// FinalizeDeliveredBatch releases every vehicle whose assignment is delivered.
// It intentionally uses a helper so each release happens before the next one.
func (h *HandoverService) FinalizeDeliveredBatch() (int, error) {
	released := 0
	for _, a := range h.store.ListAssignments() {
		if a.Status != domain.AssignmentDelivered {
			continue
		}
		if err := h.releaseVehicle(a.VehicleID); err != nil {
			return released, err
		}
		released++
	}
	return released, nil
}

func (h *HandoverService) releaseVehicle(vehicleID string) error {
	if _, err := h.store.UpdateVehicleStatus(vehicleID, domain.VehicleIdle); err != nil {
		return err
	}
	return nil
}
