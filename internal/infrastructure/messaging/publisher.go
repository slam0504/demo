package messaging

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

// Publisher wraps a Watermill Kafka publisher.
type Publisher struct {
	pub *kafka.Publisher
}

// NewPublisher creates a new Kafka publisher.
func NewPublisher(brokers []string) (*Publisher, error) {
	pub, err := kafka.NewPublisher(kafka.PublisherConfig{Brokers: brokers}, nil)
	if err != nil {
		return nil, err
	}
	return &Publisher{pub: pub}, nil
}

// Publish encodes the event and sends it to Kafka.
func (p *Publisher) Publish(ctx context.Context, topic string, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := message.NewMessage(watermill.NewUUID(), payload)
	return p.pub.Publish(topic, msg)
}
