package command

import (
	"context"

	"demo/internal/domain/card"
)

// CreateCardCommand holds data for creating a card.
type CreateCardCommand struct {
	Name        string
	Cost        int
	Faction     string
	Category    string
	SubCategory string
	Description string
}

// CreateCardHandler handles card creation.
type CreateCardHandler struct {
	Repo      card.Repository
	Publisher EventPublisher
}

// EventPublisher defines interface for publishing events.
type EventPublisher interface {
	Publish(ctx context.Context, topic string, event interface{}) error
}

// Handle executes the command.
func (h *CreateCardHandler) Handle(ctx context.Context, cmd CreateCardCommand) (*card.Card, error) {
	c := card.NewCard(cmd.Name, cmd.Cost, cmd.Faction, cmd.Category, cmd.SubCategory, cmd.Description)
	evt := card.CardCreated{
		ID:          c.ID,
		Name:        c.Name,
		Cost:        c.Cost,
		Faction:     c.Faction,
		Category:    c.Category,
		SubCategory: c.SubCategory,
		Description: c.Description,
	}
	if err := h.Repo.Save(ctx, []interface{}{evt}); err != nil {
		return nil, err
	}
	if h.Publisher != nil {
		_ = h.Publisher.Publish(ctx, "card_events", evt)
	}
	return c, nil
}
