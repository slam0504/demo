package http

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	appcmd "demo/internal/application/command"
	appquery "demo/internal/application/query"
	"demo/internal/domain/card"
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

func TestPostInvalidBody(t *testing.T) {
	repo := &mockRepo{}
	r := Router(&appcmd.CreateCardHandler{Repo: repo}, &appcmd.UpdateCardHandler{Repo: repo}, &appquery.SearchCardsHandler{Repo: repo})
	req := httptest.NewRequest("POST", "/cards", bytes.NewBufferString("{"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func TestPostRepoError(t *testing.T) {
	repo := &mockRepo{SaveFn: func(ctx context.Context, evts []interface{}) error { return errors.New("fail") }}
	r := Router(&appcmd.CreateCardHandler{Repo: repo}, &appcmd.UpdateCardHandler{Repo: repo}, &appquery.SearchCardsHandler{Repo: repo})
	body := `{"name":"n"}`
	req := httptest.NewRequest("POST", "/cards", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", w.Code)
	}
}

func TestPutInvalidID(t *testing.T) {
	repo := &mockRepo{}
	r := Router(&appcmd.CreateCardHandler{Repo: repo}, &appcmd.UpdateCardHandler{Repo: repo}, &appquery.SearchCardsHandler{Repo: repo})
	req := httptest.NewRequest("PUT", "/cards/bad", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func TestGetRepoError(t *testing.T) {
	repo := &mockRepo{SearchFn: func(ctx context.Context, name string, cost int, f, c, s string) ([]*card.Card, error) {
		return nil, errors.New("fail")
	}}
	r := Router(&appcmd.CreateCardHandler{Repo: repo}, &appcmd.UpdateCardHandler{Repo: repo}, &appquery.SearchCardsHandler{Repo: repo})
	req := httptest.NewRequest("GET", "/cards", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", w.Code)
	}
}

func TestGetSuccess(t *testing.T) {
	repo := &mockRepo{SearchFn: func(ctx context.Context, name string, cost int, f, c, s string) ([]*card.Card, error) {
		return []*card.Card{{Name: "N"}}, nil
	}}
	r := Router(&appcmd.CreateCardHandler{Repo: repo}, &appcmd.UpdateCardHandler{Repo: repo}, &appquery.SearchCardsHandler{Repo: repo})
	req := httptest.NewRequest("GET", "/cards", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}
