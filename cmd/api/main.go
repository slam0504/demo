package main

import (
	"context"
	"log"
	"net/http"

	appcmd "demo/internal/application/command"
	appquery "demo/internal/application/query"
	"demo/internal/infrastructure/auth"
	"demo/internal/infrastructure/cache"
	"demo/internal/infrastructure/deckstore"
	"demo/internal/infrastructure/eventstore"
	"demo/internal/infrastructure/messaging"
	httpiface "demo/internal/interfaces/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() func(context.Context) error {
	exp, err := stdouttrace.New()
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.Default()),
	)
	otel.SetTracerProvider(tp)
	return tp.Shutdown
}

func main() {
	shutdown := initTracer()
	defer func() { _ = shutdown(context.Background()) }()

	es, err := eventstore.NewMySQLStore("root@tcp(127.0.0.1:3306)/card_service?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	// wrap repository with redis cache
	repo := cache.NewRedisRepository(es, "localhost:6379")
	deckRepo := deckstore.NewInMemoryStore()
	authSvc := auth.NewService()
	publisher, err := messaging.NewPublisher([]string{"localhost:9092"})
	if err != nil {
		log.Println("failed to create publisher", err)
	}

	createHandler := &appcmd.CreateCardHandler{Repo: repo, Publisher: publisher}
	updateHandler := &appcmd.UpdateCardHandler{Repo: repo, Publisher: publisher}
	searchHandler := &appquery.SearchCardsHandler{Repo: repo}
	deckHandler := &appcmd.CreateDeckHandler{Repo: deckRepo}

	r := httpiface.Router(authSvc, createHandler, updateHandler, searchHandler, deckHandler)
	log.Println("http server started on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
