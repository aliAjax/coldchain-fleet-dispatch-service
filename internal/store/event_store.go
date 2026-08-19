package store

import (
	"sort"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
)

func (s *Store) AppendEvent(ev domain.AuditEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ev.ID == "" {
		return platform.ErrInvalid
	}
	s.events = append(s.events, ev)
	return nil
}

func (s *Store) ListEvents(entityID string) []domain.AuditEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.AuditEvent, 0, len(s.events))
	for _, ev := range s.events {
		if entityID == "" || ev.EntityID == entityID {
			out = append(out, ev)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out
}
