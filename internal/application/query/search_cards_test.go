package query

import (
	"context"
	"demo/internal/domain/card"
	"errors"
	"testing"
)

type mockRepo struct {
	SearchFn func(ctx context.Context, name string, cost int, f, c, s string) ([]*card.Card, error)
}

func (m *mockRepo) Save(ctx context.Context, evts []interface{}) error      { return nil }
func (m *mockRepo) Load(ctx context.Context, id string) (*card.Card, error) { return nil, nil }
func (m *mockRepo) Search(ctx context.Context, name string, cost int, faction, category, sub string) ([]*card.Card, error) {
	if m.SearchFn != nil {
		return m.SearchFn(ctx, name, cost, faction, category, sub)
	}
	return nil, nil
}

func TestSearchCardsError(t *testing.T) {
	repo := &mockRepo{SearchFn: func(ctx context.Context, name string, cost int, f, c, s string) ([]*card.Card, error) {
		return nil, errors.New("err")
	}}
	h := &SearchCardsHandler{Repo: repo}
	_, err := h.Handle(context.Background(), SearchCardsQuery{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSearchCards(t *testing.T) {
	repo := &mockRepo{SearchFn: func(ctx context.Context, name string, cost int, f, c, s string) ([]*card.Card, error) {
		return []*card.Card{{Name: "N"}}, nil
	}}
	h := &SearchCardsHandler{Repo: repo}
	res, err := h.Handle(context.Background(), SearchCardsQuery{})
	if err != nil || len(res) != 1 {
		t.Fatalf("unexpected %v %v", res, err)
	}
}
