package domain

import "time"

type AuditEvent struct {
	ID        string
	EntityID  string
	Kind      string
	Payload   string
	CreatedAt time.Time
}
