package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	appcmd "demo/internal/application/command"
	appquery "demo/internal/application/query"
	"demo/internal/domain/card"
	"demo/internal/infrastructure/auth"
	"demo/internal/infrastructure/deckstore"
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
	authSvc := auth.NewService()
	deckRepo := deckstore.NewInMemoryStore()
	r := Router(authSvc, &appcmd.CreateCardHandler{Repo: repo}, &appcmd.UpdateCardHandler{Repo: repo}, &appquery.SearchCardsHandler{Repo: repo}, &appcmd.CreateDeckHandler{Repo: deckRepo})
	req := httptest.NewRequest("POST", "/cards", bytes.NewBufferString("{"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func TestPostRepoError(t *testing.T) {
	repo := &mockRepo{SaveFn: func(ctx context.Context, evts []interface{}) error { return errors.New("fail") }}
	authSvc := auth.NewService()
	deckRepo := deckstore.NewInMemoryStore()
	r := Router(authSvc, &appcmd.CreateCardHandler{Repo: repo}, &appcmd.UpdateCardHandler{Repo: repo}, &appquery.SearchCardsHandler{Repo: repo}, &appcmd.CreateDeckHandler{Repo: deckRepo})
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
	authSvc := auth.NewService()
	deckRepo := deckstore.NewInMemoryStore()
	r := Router(authSvc, &appcmd.CreateCardHandler{Repo: repo}, &appcmd.UpdateCardHandler{Repo: repo}, &appquery.SearchCardsHandler{Repo: repo}, &appcmd.CreateDeckHandler{Repo: deckRepo})
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
	authSvc := auth.NewService()
	deckRepo := deckstore.NewInMemoryStore()
	r := Router(authSvc, &appcmd.CreateCardHandler{Repo: repo}, &appcmd.UpdateCardHandler{Repo: repo}, &appquery.SearchCardsHandler{Repo: repo}, &appcmd.CreateDeckHandler{Repo: deckRepo})
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
	authSvc := auth.NewService()
	deckRepo := deckstore.NewInMemoryStore()
	r := Router(authSvc, &appcmd.CreateCardHandler{Repo: repo}, &appcmd.UpdateCardHandler{Repo: repo}, &appquery.SearchCardsHandler{Repo: repo}, &appcmd.CreateDeckHandler{Repo: deckRepo})
	req := httptest.NewRequest("GET", "/cards", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}

func TestLoginAndCreateDeck(t *testing.T) {
	repo := &mockRepo{}
	authSvc := auth.NewService()
	deckRepo := deckstore.NewInMemoryStore()
	r := Router(authSvc, &appcmd.CreateCardHandler{Repo: repo}, &appcmd.UpdateCardHandler{Repo: repo}, &appquery.SearchCardsHandler{Repo: repo}, &appcmd.CreateDeckHandler{Repo: deckRepo})

	loginReq := httptest.NewRequest("POST", "/login", bytes.NewBufferString(`{"username":"user","password":"password"}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, loginReq)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	token := w.Body.String()
	// token value like {"token":"..."}
	var resp struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	deckBody := `{"name":"d","cardIDs":[]}`
	req := httptest.NewRequest("POST", "/decks", bytes.NewBufferString(deckBody))
	req.Header.Set("Authorization", "Bearer "+resp.Token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w2.Code)
	}
}
