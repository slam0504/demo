package deckstore

import (
	"context"
	"sync"

	"demo/internal/domain/deck"
	"github.com/google/uuid"
)

// InMemoryStore is a simple in-memory deck repository.
type InMemoryStore struct {
	mu    sync.RWMutex
	decks map[uuid.UUID]*deck.Deck
}

// NewInMemoryStore creates the store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{decks: make(map[uuid.UUID]*deck.Deck)}
}

// Save stores or updates a deck.
func (s *InMemoryStore) Save(ctx context.Context, d *deck.Deck) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.decks[d.ID] = d
	return nil
}

// Load retrieves a deck by id.
func (s *InMemoryStore) Load(ctx context.Context, id uuid.UUID) (*deck.Deck, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if d, ok := s.decks[id]; ok {
		copy := *d
		return &copy, nil
	}
	return nil, nil
}

var _ deck.Repository = (*InMemoryStore)(nil)
