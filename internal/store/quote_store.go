package store

import (
	"sort"

	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
)

func (s *Store) SaveQuote(q domain.Quote) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if q.ID == "" {
		return platform.ErrInvalid
	}
	s.quotes[q.ID] = q
	return nil
}

func (s *Store) GetQuote(id string) (domain.Quote, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q, ok := s.quotes[id]
	if !ok {
		return domain.Quote{}, platform.ErrNotFound
	}
	return q, nil
}

func (s *Store) ListQuotes() []domain.Quote {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Quote, 0, len(s.quotes))
	for _, q := range s.quotes {
		out = append(out, q)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out
}
