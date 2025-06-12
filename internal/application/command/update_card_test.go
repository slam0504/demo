package command

import (
	"context"
	"demo/internal/domain/card"
	"errors"
	"github.com/google/uuid"
	"testing"
)

func TestUpdateCardLoadError(t *testing.T) {
	repo := &mockRepo{}
	repo.LoadFn = func(ctx context.Context, id string) (*card.Card, error) { return nil, errors.New("load") }
	h := &UpdateCardHandler{Repo: repo}
	_, err := h.Handle(context.Background(), UpdateCardCommand{ID: uuid.New()})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateCardSuccess(t *testing.T) {
	c := card.NewCard("n", 1, "f", "c", "s", "d")
	repo := &mockRepo{}
	repo.LoadFn = func(ctx context.Context, id string) (*card.Card, error) { return c, nil }
	repo.SaveFn = func(ctx context.Context, evts []interface{}) error { return nil }
	h := &UpdateCardHandler{Repo: repo}
	updated, err := h.Handle(context.Background(), UpdateCardCommand{ID: c.ID, Name: "x"})
	if err != nil || updated == nil {
		t.Fatalf("unexpected %v %v", updated, err)
	}
}
