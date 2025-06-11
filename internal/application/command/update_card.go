package command

import (
	"context"

	"demo/internal/domain/card"
	"github.com/google/uuid"
)

type UpdateCardCommand struct {
	ID          uuid.UUID
	Name        string
	Cost        int
	Faction     string
	Category    string
	SubCategory string
	Description string
}

type UpdateCardHandler struct {
	Repo      card.Repository
	Publisher EventPublisher
}

func (h *UpdateCardHandler) Handle(ctx context.Context, cmd UpdateCardCommand) (*card.Card, error) {
	existing, err := h.Repo.Load(ctx, cmd.ID.String())
	if err != nil || existing == nil {
		return nil, err
	}
	evt := card.CardUpdated{
		ID:          cmd.ID,
		Name:        cmd.Name,
		Cost:        cmd.Cost,
		Faction:     cmd.Faction,
		Category:    cmd.Category,
		SubCategory: cmd.SubCategory,
		Description: cmd.Description,
	}
	if err := h.Repo.Save(ctx, []interface{}{evt}); err != nil {
		return nil, err
	}
	if h.Publisher != nil {
		_ = h.Publisher.Publish(ctx, "card_events", evt)
	}
	updated, _ := h.Repo.Load(ctx, cmd.ID.String())
	return updated, nil
}
