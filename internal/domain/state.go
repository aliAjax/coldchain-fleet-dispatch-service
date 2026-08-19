package domain

type ShipmentStatus string

const (
	ShipmentPending   ShipmentStatus = "pending"
	ShipmentAssigned  ShipmentStatus = "assigned"
	ShipmentInTransit ShipmentStatus = "in_transit"
	ShipmentDelivered ShipmentStatus = "delivered"
	ShipmentFailed    ShipmentStatus = "failed"
	ShipmentCancelled ShipmentStatus = "cancelled"
)

type AssignmentStatus string

const (
	AssignmentPending   AssignmentStatus = "pending"
	AssignmentAssigned  AssignmentStatus = "assigned"
	AssignmentInTransit AssignmentStatus = "in_transit"
	AssignmentDelivered AssignmentStatus = "delivered"
	AssignmentFailed    AssignmentStatus = "failed"
)

type TempZone string

const (
	ZoneFrozen  TempZone = "frozen"
	ZoneChilled TempZone = "chilled"
	ZoneAmbient TempZone = "ambient"
)

type VehicleStatus string

const (
	VehicleIdle        VehicleStatus = "idle"
	VehicleBusy        VehicleStatus = "busy"
	VehicleMaintenance VehicleStatus = "maintenance"
)
