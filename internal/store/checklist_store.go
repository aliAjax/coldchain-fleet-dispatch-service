package store

import (
	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
)

func (s *Store) SaveChecklist(c domain.HandoverChecklist) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c.AssignmentID == "" {
		return platform.ErrInvalid
	}
	s.checklists[c.AssignmentID] = c
	return nil
}

func (s *Store) GetChecklist(assignmentID string) (domain.HandoverChecklist, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.checklists[assignmentID]
	if !ok {
		return domain.HandoverChecklist{}, platform.ErrNotFound
	}
	return c, nil
}
