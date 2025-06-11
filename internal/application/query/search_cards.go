package query

import (
	"context"

	"demo/internal/domain/card"
)

// SearchCardsQuery defines search parameters.
type SearchCardsQuery struct {
	Name     string
	Cost     int
	Faction  string
	Category string
	Sub      string
}

// SearchCardsHandler handles searching for cards.
type SearchCardsHandler struct {
	Repo card.Repository
}

func (h *SearchCardsHandler) Handle(ctx context.Context, q SearchCardsQuery) ([]*card.Card, error) {
	return h.Repo.Search(ctx, q.Name, q.Cost, q.Faction, q.Category, q.Sub)
}
