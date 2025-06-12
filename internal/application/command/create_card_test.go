package command

import (
	"context"
	"demo/internal/domain/card"
	"errors"
	"testing"
)

type mockRepo struct {
	SaveFn   func(ctx context.Context, evts []interface{}) error
	LoadFn   func(ctx context.Context, id string) (*card.Card, error)
	SearchFn func(ctx context.Context, name string, cost int, faction, category, sub string) ([]*card.Card, error)
}

func (m *mockRepo) Save(ctx context.Context, evts []interface{}) error {
	if m.SaveFn != nil {
		return m.SaveFn(ctx, evts)
	}
	return nil
}
func (m *mockRepo) Load(ctx context.Context, id string) (*card.Card, error) {
	if m.LoadFn != nil {
		return m.LoadFn(ctx, id)
	}
	return nil, nil
}
func (m *mockRepo) Search(ctx context.Context, name string, cost int, faction, category, sub string) ([]*card.Card, error) {
	if m.SearchFn != nil {
		return m.SearchFn(ctx, name, cost, faction, category, sub)
	}
	return nil, nil
}

func TestCreateCardHandlerError(t *testing.T) {
	repo := &mockRepo{SaveFn: func(ctx context.Context, evts []interface{}) error { return errors.New("saveerr") }}
	h := &CreateCardHandler{Repo: repo}
	_, err := h.Handle(context.Background(), CreateCardCommand{Name: "n"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateCardHandlerSuccess(t *testing.T) {
	repo := &mockRepo{}
	h := &CreateCardHandler{Repo: repo}
	c, err := h.Handle(context.Background(), CreateCardCommand{Name: "n"})
	if err != nil || c == nil || c.Name != "n" {
		t.Fatalf("unexpected result %v %v", c, err)
	}
}
