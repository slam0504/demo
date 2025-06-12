package messaging

import (
	"context"
	"testing"
)

func TestPublishMarshalError(t *testing.T) {
	p := &Publisher{}
	err := p.Publish(context.Background(), "t", map[interface{}]string{1: "x"})
	if err == nil {
		t.Fatal("expected error")
	}
}
