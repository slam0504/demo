package eventstore

import (
	"context"
	"demo/internal/domain/card"
	"testing"
)

func TestInMemorySaveLoad(t *testing.T) {
	repo := NewInMemoryStore()
	c := card.NewCard("N", 1, "F", "C", "S", "D")
	evt := card.CardCreated{ID: c.ID, Name: c.Name, Cost: c.Cost, Faction: c.Faction, Category: c.Category, SubCategory: c.SubCategory, Description: c.Description}
	if err := repo.Save(context.Background(), []interface{}{evt}); err != nil {
		t.Fatal(err)
	}
	loaded, err := repo.Load(context.Background(), c.ID.String())
	if err != nil || loaded == nil || loaded.Name != "N" {
		t.Fatalf("load failed: %+v %v", loaded, err)
	}
}

func TestInMemoryUnknownEvent(t *testing.T) {
	repo := NewInMemoryStore()
	if err := repo.Save(context.Background(), []interface{}{struct{}{}}); err != nil {
		t.Fatal(err)
	}
	cards, err := repo.Search(context.Background(), "", 0, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 0 {
		t.Fatalf("expected 0 cards got %d", len(cards))
	}
}
