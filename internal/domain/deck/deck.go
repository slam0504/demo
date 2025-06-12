package deck

import "github.com/google/uuid"

// Deck represents a collection of cards owned by a user.
type Deck struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	Name    string
	CardIDs []uuid.UUID
}

// NewDeck creates a new deck for a user.
func NewDeck(userID uuid.UUID, name string, cardIDs []uuid.UUID) *Deck {
	return &Deck{
		ID:      uuid.New(),
		UserID:  userID,
		Name:    name,
		CardIDs: append([]uuid.UUID(nil), cardIDs...),
	}
}
