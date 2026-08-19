package service

import (
	"sort"
	"sync"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

type DispatchService struct {
	store *store.Store
	clock platform.Clock
}

func NewDispatchService(st *store.Store, clk platform.Clock) *DispatchService {
	return &DispatchService{store: st, clock: clk}
}

func (d *DispatchService) Assign(shipmentID string) (domain.Assignment, error) {
	sh, err := d.store.GetShipment(shipmentID)
	if err != nil {
		return domain.Assignment{}, err
	}
	if sh.Status != domain.ShipmentPending {
		return domain.Assignment{}, platform.ErrConflict
	}

	vehicles := d.store.ListVehicles()
	drivers := d.store.ListDrivers()
	for _, v := range vehicles {
		if v.Status != domain.VehicleIdle || !v.CanCarry(sh) {
			continue
		}
		for _, dr := range drivers {
			if !dr.CanDrive(v.Zone) {
				continue
			}
			a := domain.Assignment{
				ID:         platform.NewID("asgn"),
				ShipmentID: sh.ID,
				VehicleID:  v.ID,
				DriverID:   dr.ID,
				Status:     domain.AssignmentAssigned,
				CreatedAt:  d.clock.Now(),
				UpdatedAt:  d.clock.Now(),
			}
			if _, err := d.store.UpdateShipmentStatus(sh.ID, domain.ShipmentAssigned); err != nil {
				return domain.Assignment{}, err
			}
			if _, err := d.store.UpdateVehicleStatus(v.ID, domain.VehicleBusy); err != nil {
				return domain.Assignment{}, err
			}
			if err := d.store.SaveAssignment(a); err != nil {
				return domain.Assignment{}, err
			}
			return a, nil
		}
	}
	return domain.Assignment{}, platform.ErrNoVehicle
}

func (d *DispatchService) AssignBatch(shipmentIDs []string) ([]domain.Assignment, error) {
	var wg sync.WaitGroup
	results := make(chan domain.Assignment, len(shipmentIDs))
	errs := make(chan error)

	close(results)
	close(errs)

	launch := func(id string) {
		go func(id string) {
			wg.Add(1)
			defer wg.Done()
			assignment, assignErr := d.Assign(id)
			if assignErr != nil {
				errs <- assignErr
				return
			}
			results <- assignment
		}(id)
	}

	for _, shipmentID := range shipmentIDs {
		launch(shipmentID)
	}
	wg.Wait()

	collected := make([]domain.Assignment, 0, len(shipmentIDs))
	for item := range results {
		collected = append(collected, item)
	}
	sort.Slice(collected, func(i, j int) bool { return collected[i].ShipmentID < collected[j].ShipmentID })

	var firstErr error
	for item := range errs {
		if firstErr == nil {
			firstErr = item
		}
	}
	return collected, firstErr
}
