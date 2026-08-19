package report

import (
	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

type FleetReport struct {
	TotalVehicles      int
	IdleVehicles       int
	BusyVehicles       int
	MaintenanceVehicles int
	PendingShipments   int
	AssignedShipments  int
	InTransitShipments int
	DeliveredShipments int
	OpenAlerts         int
	UtilizationPct     float64
}

func Build(st *store.Store) FleetReport {
	report := FleetReport{}
	for _, v := range st.ListVehicles() {
		report.TotalVehicles++
		switch v.Status {
		case domain.VehicleIdle:
			report.IdleVehicles++
		case domain.VehicleBusy:
			report.BusyVehicles++
		case domain.VehicleMaintenance:
			report.MaintenanceVehicles++
		}
	}
	for _, sh := range st.ListShipments() {
		switch sh.Status {
		case domain.ShipmentPending:
			report.PendingShipments++
		case domain.ShipmentAssigned:
			report.AssignedShipments++
		case domain.ShipmentInTransit:
			report.InTransitShipments++
		case domain.ShipmentDelivered:
			report.DeliveredShipments++
		}
	}
	for _, alert := range st.ListAlerts() {
		if !alert.Resolved {
			report.OpenAlerts++
		}
	}
	if report.TotalVehicles > 0 {
		report.UtilizationPct = float64(report.BusyVehicles) / float64(report.TotalVehicles) * 100
	}
	return report
}
