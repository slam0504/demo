package card

import "github.com/google/uuid"

// CardCreated is emitted when a new card is created.
type CardCreated struct {
	ID          uuid.UUID
	Name        string
	Cost        int
	Faction     string
	Category    string
	SubCategory string
	Description string
}

// CardUpdated is emitted when an existing card is updated.
type CardUpdated struct {
	ID          uuid.UUID
	Name        string
	Cost        int
	Faction     string
	Category    string
	SubCategory string
	Description string
}
