package tests

import (
	"context"
	"testing"

	"demo/internal/application/command"
	"demo/internal/application/query"
	"demo/internal/infrastructure/eventstore"
)

func TestCreateAndSearchCard(t *testing.T) {
	repo := eventstore.NewInMemoryStore()
	handler := &command.CreateCardHandler{Repo: repo}

	card, err := handler.Handle(context.Background(), command.CreateCardCommand{
		Name:        "Test",
		Cost:        1,
		Faction:     "Human",
		Category:    "Soldier",
		SubCategory: "Infantry",
		Description: "test card",
	})
	if err != nil {
		t.Fatal(err)
	}
	if card.Name != "Test" {
		t.Fatalf("expected name Test got %s", card.Name)
	}

	// search
	queryHandler := &query.SearchCardsHandler{Repo: repo}
	cards, err := queryHandler.Handle(context.Background(), query.SearchCardsQuery{Name: "Test"})
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 1 {
		t.Fatalf("expected 1 card got %d", len(cards))
	}
}
