package card

import (
	"github.com/google/uuid"
	"testing"
)

func TestNewCard(t *testing.T) {
	c := NewCard("Name", 1, "Faction", "Category", "Sub", "Desc")
	if c.Name != "Name" || c.Cost != 1 || c.Faction != "Faction" || c.Category != "Category" || c.SubCategory != "Sub" || c.Description != "Desc" {
		t.Fatalf("fields not set correctly: %+v", c)
	}
	if c.ID == (uuid.UUID{}) {
		t.Fatal("id not set")
	}
}
