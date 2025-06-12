package http

import (
	"fmt"
	"net/http"

	appcmd "demo/internal/application/command"
	appquery "demo/internal/application/query"
	"demo/internal/i18n"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Router sets up HTTP routes using Gin.
func Router(createHandler *appcmd.CreateCardHandler, updateHandler *appcmd.UpdateCardHandler, searchHandler *appquery.SearchCardsHandler) http.Handler {
	r := gin.New()
	r.Use(otelgin.Middleware("card_service"))

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

	return r
}
