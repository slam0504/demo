package eventstore

import (
	"demo/internal/domain/card"
	"github.com/google/uuid"
	"testing"
)

func TestEventCardID(t *testing.T) {
	c := card.CardCreated{ID: uuid.New()}
	id, err := eventCardID(c)
	if err != nil || id == "" {
		t.Fatalf("unexpected result %s %v", id, err)
	}
	_, err = eventCardID(struct{}{})
	if err == nil {
		t.Fatal("expected error for unknown event")
	}
}
