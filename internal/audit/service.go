package audit

import (
	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

type Service struct {
	store *store.Store
	clock platform.Clock
}

func NewService(st *store.Store, clk platform.Clock) *Service {
	return &Service{store: st, clock: clk}
}

func (s *Service) Record(entityID, kind, payload string) (domain.AuditEvent, error) {
	if entityID == "" || kind == "" {
		return domain.AuditEvent{}, platform.ErrInvalid
	}
	event := domain.AuditEvent{
		ID:        platform.NewID("evt"),
		EntityID:  entityID,
		Kind:      kind,
		Payload:   payload,
		CreatedAt: s.clock.Now(),
	}
	if err := s.store.AppendEvent(event); err != nil {
		return domain.AuditEvent{}, err
	}
	return event, nil
}

func (s *Service) List(entityID string) []domain.AuditEvent {
	return s.store.ListEvents(entityID)
}

func (s *Service) HasEvent(entityID, kind string) bool {
	for _, event := range s.store.ListEvents(entityID) {
		if event.Kind == kind {
			return true
		}
	}
	return false
}
