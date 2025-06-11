package eventstore

import (
	"context"
	"sync"

	"demo/internal/domain/card"
)

type inMemoryStore struct {
	mu     sync.RWMutex
	events map[string][]interface{}
}

// NewInMemoryStore creates an in-memory event store.
func NewInMemoryStore() card.Repository {
	return &inMemoryStore{events: make(map[string][]interface{})}
}

func (s *inMemoryStore) Save(ctx context.Context, events []interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(events) == 0 {
		return nil
	}
	var id string
	switch e := events[0].(type) {
	case card.CardCreated:
		id = e.ID.String()
	case card.CardUpdated:
		id = e.ID.String()
	default:
		return nil
	}
	s.events[id] = append(s.events[id], events...)
	return nil
}

func (s *inMemoryStore) Load(ctx context.Context, id string) (*card.Card, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	evs := s.events[id]
	if len(evs) == 0 {
		return nil, nil
	}
	c := &card.Card{}
	for _, e := range evs {
		switch evt := e.(type) {
		case card.CardCreated:
			c.ID = evt.ID
			c.Name = evt.Name
			c.Cost = evt.Cost
			c.Faction = evt.Faction
			c.Category = evt.Category
			c.SubCategory = evt.SubCategory
			c.Description = evt.Description
		case card.CardUpdated:
			c.Name = evt.Name
			c.Cost = evt.Cost
			c.Faction = evt.Faction
			c.Category = evt.Category
			c.SubCategory = evt.SubCategory
			c.Description = evt.Description
		}
	}
	return c, nil
}

func (s *inMemoryStore) Search(ctx context.Context, name string, cost int, faction, category, sub string) ([]*card.Card, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var cards []*card.Card
	for _, evs := range s.events {
		c := &card.Card{}
		for _, e := range evs {
			switch evt := e.(type) {
			case card.CardCreated:
				c.ID = evt.ID
				c.Name = evt.Name
				c.Cost = evt.Cost
				c.Faction = evt.Faction
				c.Category = evt.Category
				c.SubCategory = evt.SubCategory
				c.Description = evt.Description
			case card.CardUpdated:
				c.Name = evt.Name
				c.Cost = evt.Cost
				c.Faction = evt.Faction
				c.Category = evt.Category
				c.SubCategory = evt.SubCategory
				c.Description = evt.Description
			}
		}
		if (name == "" || c.Name == name) &&
			(cost == 0 || c.Cost == cost) &&
			(faction == "" || c.Faction == faction) &&
			(category == "" || c.Category == category) &&
			(sub == "" || c.SubCategory == sub) {
			tmp := *c
			cards = append(cards, &tmp)
		}
	}
	return cards, nil
}
