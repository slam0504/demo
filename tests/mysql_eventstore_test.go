package tests

import (
	"context"
	"testing"

	"demo/internal/application/command"
	"demo/internal/application/query"
	"demo/internal/infrastructure/eventstore"
)

func TestMySQLCreateAndSearchCard(t *testing.T) {
	repo, err := eventstore.NewMySQLStore("root@tcp(127.0.0.1:3306)/card_test?parseTime=true")
	if err != nil {
		t.Skipf("mysql not available: %v", err)
	}
	// clean table
	_ = repo.DB.Exec("TRUNCATE TABLE event_records")

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
	queryHandler := &query.SearchCardsHandler{Repo: repo}
	cards, err := queryHandler.Handle(context.Background(), query.SearchCardsQuery{Name: "Test"})
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 1 {
		t.Fatalf("expected 1 card got %d", len(cards))
	}
}
