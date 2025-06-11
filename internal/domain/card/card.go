package card

import "github.com/google/uuid"

// Card represents a collectible card in the system.
type Card struct {
	ID          uuid.UUID
	Name        string
	Cost        int
	Faction     string
	Category    string
	SubCategory string
	Description string
}

func NewCard(name string, cost int, faction, category, subCategory, description string) *Card {
	return &Card{
		ID:          uuid.New(),
		Name:        name,
		Cost:        cost,
		Faction:     faction,
		Category:    category,
		SubCategory: subCategory,
		Description: description,
	}
}
