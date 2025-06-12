package http

import (
	"fmt"
	"net/http"

	appcmd "demo/internal/application/command"
	appquery "demo/internal/application/query"
	"demo/internal/i18n"
	"demo/internal/infrastructure/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Router sets up HTTP routes using Gin.
func Router(authSvc *auth.Service, createHandler *appcmd.CreateCardHandler, updateHandler *appcmd.UpdateCardHandler, searchHandler *appquery.SearchCardsHandler, deckHandler *appcmd.CreateDeckHandler) http.Handler {
	r := gin.New()
	r.Use(otelgin.Middleware("card_service"))

	r.POST("/login", func(c *gin.Context) {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		if token, ok := authSvc.Login(body.Username, body.Password); ok {
			c.JSON(http.StatusOK, gin.H{"token": token})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		}
	})

	r.POST("/cards", func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		var cmd appcmd.CreateCardCommand
		if err := c.ShouldBindJSON(&cmd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(lang, "invalid_body")})
			return
		}
		card, err := createHandler.Handle(c.Request.Context(), cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(lang, "internal_error")})
			return
		}
		c.JSON(http.StatusOK, i18n.TranslateCard(lang, card))
	})

	r.PUT("/cards/:id", func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(lang, "invalid_id")})
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
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(lang, "invalid_body")})
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
		card, err := updateHandler.Handle(c.Request.Context(), cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(lang, "internal_error")})
			return
		}
		c.JSON(http.StatusOK, i18n.TranslateCard(lang, card))
	})

	r.GET("/cards", func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		q := appquery.SearchCardsQuery{
			Name:     c.Query("name"),
			Faction:  c.Query("faction"),
			Category: c.Query("category"),
			Sub:      c.Query("sub"),
		}
		if cost := c.Query("cost"); cost != "" {
			_, _ = fmt.Sscanf(cost, "%d", &q.Cost)
		}
		cards, err := searchHandler.Handle(c.Request.Context(), q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(lang, "internal_error")})
			return
		}
		c.JSON(http.StatusOK, i18n.TranslateCards(lang, cards))
	})

	r.POST("/decks", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		var token string
		fmt.Sscanf(authHeader, "Bearer %s", &token)
		userID, ok := authSvc.Authenticate(token)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var body struct {
			Name    string
			CardIDs []string
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		var ids []uuid.UUID
		for _, s := range body.CardIDs {
			id, err := uuid.Parse(s)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card id"})
				return
			}
			ids = append(ids, id)
		}
		cmd := appcmd.CreateDeckCommand{UserID: userID, Name: body.Name, CardIDs: ids}
		d, err := deckHandler.Handle(c.Request.Context(), cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": d.ID.String()})
	})

	return r
}
