package domain

type ChecklistItem struct {
	Key   string
	Label string
	Done  bool
}

type HandoverChecklist struct {
	AssignmentID string
	Items        []ChecklistItem
}

func (c HandoverChecklist) Complete() bool {
	for _, item := range c.Items {
		if !item.Done {
			return false
		}
	}
	return len(c.Items) > 0
}

func (c HandoverChecklist) PendingKeys() []string {
	out := make([]string, 0)
	for _, item := range c.Items {
		if !item.Done {
			out = append(out, item.Key)
		}
	}
	return out
}
