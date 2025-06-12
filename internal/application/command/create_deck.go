package command

import (
	"context"

	"demo/internal/domain/deck"
	"github.com/google/uuid"
)

// CreateDeckCommand contains info needed to create a deck.
type CreateDeckCommand struct {
	UserID  uuid.UUID
	Name    string
	CardIDs []uuid.UUID
}

// CreateDeckHandler handles deck creation.
type CreateDeckHandler struct{ Repo deck.Repository }

// Handle creates the deck and persists it.
func (h *CreateDeckHandler) Handle(ctx context.Context, cmd CreateDeckCommand) (*deck.Deck, error) {
	d := deck.NewDeck(cmd.UserID, cmd.Name, cmd.CardIDs)
	if err := h.Repo.Save(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}
