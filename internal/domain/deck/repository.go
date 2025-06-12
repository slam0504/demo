package deck

import (
	"context"
	"github.com/google/uuid"
)

// Repository provides persistence for decks.
type Repository interface {
	Save(ctx context.Context, d *Deck) error
	Load(ctx context.Context, id uuid.UUID) (*Deck, error)
}
