package service

import (
	"fmt"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
)

type ChecklistService struct {
	store *store.Store
	clock platform.Clock
}

func NewChecklistService(st *store.Store, clk platform.Clock) *ChecklistService {
	return &ChecklistService{store: st, clock: clk}
}

func (c *ChecklistService) Build(assignmentID string) (domain.HandoverChecklist, error) {
	if _, err := c.store.GetAssignment(assignmentID); err != nil {
		return domain.HandoverChecklist{}, fmt.Errorf("build checklist: %w", err)
	}
	checklist := domain.HandoverChecklist{
		AssignmentID: assignmentID,
		Items: []domain.ChecklistItem{
			{Key: "seal", Label: "seal inspected"},
			{Key: "temperature", Label: "temperature confirmed"},
			{Key: "papers", Label: "papers signed"},
		},
	}
	if err := c.store.SaveChecklist(checklist); err != nil {
		return domain.HandoverChecklist{}, fmt.Errorf("build checklist: %w", err)
	}
	return checklist, nil
}

func (c *ChecklistService) Mark(assignmentID, key string, done bool) (domain.HandoverChecklist, error) {
	checklist, err := c.store.GetChecklist(assignmentID)
	if err != nil {
		return domain.HandoverChecklist{}, fmt.Errorf("mark checklist: %w", err)
	}
	found := false
	for i := range checklist.Items {
		if checklist.Items[i].Key == key {
			checklist.Items[i].Done = done
			found = true
		}
	}
	if !found {
		return domain.HandoverChecklist{}, fmt.Errorf("mark checklist: %w", platform.ErrNotFound)
	}
	if err := c.store.SaveChecklist(checklist); err != nil {
		return domain.HandoverChecklist{}, fmt.Errorf("mark checklist: %w", err)
	}
	return checklist, nil
}
