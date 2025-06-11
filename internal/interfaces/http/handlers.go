package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	appcmd "demo/internal/application/command"
	appquery "demo/internal/application/query"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Router sets up HTTP routes.
func Router(createHandler *appcmd.CreateCardHandler, updateHandler *appcmd.UpdateCardHandler, searchHandler *appquery.SearchCardsHandler) http.Handler {
	r := chi.NewRouter()
	r.Method("POST", "/cards", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var cmd appcmd.CreateCardCommand
		if err := json.NewDecoder(req.Body).Decode(&cmd); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		card, err := createHandler.Handle(req.Context(), cmd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(card)
	}), "create_card"))

	r.Method("PUT", "/cards/{id}", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		idStr := chi.URLParam(req, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		var body struct {
			Name        string
			Cost        int
			Faction     string
			Category    string
			SubCategory string
			Description string
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		cmd := appcmd.UpdateCardCommand{
			ID:          id,
			Name:        body.Name,
			Cost:        body.Cost,
			Faction:     body.Faction,
			Category:    body.Category,
			SubCategory: body.SubCategory,
			Description: body.Description,
		}
		card, err := updateHandler.Handle(req.Context(), cmd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(card)
	}), "update_card"))

	r.Method("GET", "/cards", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		q := appquery.SearchCardsQuery{
			Name:     req.URL.Query().Get("name"),
			Faction:  req.URL.Query().Get("faction"),
			Category: req.URL.Query().Get("category"),
			Sub:      req.URL.Query().Get("sub"),
		}
		if cost := req.URL.Query().Get("cost"); cost != "" {
			if _, err := fmt.Sscanf(cost, "%d", &q.Cost); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		cards, err := searchHandler.Handle(req.Context(), q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(cards)
	}), "search_cards"))

	return r
}
