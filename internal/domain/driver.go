package domain

type Driver struct {
	ID            string
	Name          string
	Certification TempZone
	Available     bool
}

func (d Driver) CanDrive(zone TempZone) bool {
	return d.Available && d.Certification == zone
}
