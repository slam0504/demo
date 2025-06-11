# Card Service

This project implements a simple card management service following principles from Clean Architecture, DDD, CQRS, event sourcing and event-driven design. Events are published via Kafka using the Watermill library and tracing is enabled with OpenTelemetry.

## Development

Use the provided Dev Container configuration with VS Code to start a development environment. The service exposes HTTP endpoints:

- `POST /cards` – create a card
- `PUT /cards/{id}` – update a card
- `GET /cards` – search for cards

Run tests with `go test ./...` and lint with `golangci-lint run`.
