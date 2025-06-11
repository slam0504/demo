package card

import "context"

// Repository defines methods for persisting cards via event sourcing.
type Repository interface {
	Save(ctx context.Context, events []interface{}) error
	Load(ctx context.Context, id string) (*Card, error)
	Search(ctx context.Context, name string, cost int, faction, category, sub string) ([]*Card, error)
}
