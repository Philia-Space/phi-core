package event

import (
	"context"
	"encoding/json"
	"time"
)

// Envelope wraps a domain event with metadata for transport.
type Envelope struct {
	EventID     string          `json:"event_id"`
	EventType   string          `json:"event_type"`
	AggregateID string          `json:"aggregate_id"`
	ServiceName string          `json:"service_name"`
	Payload     json.RawMessage `json:"payload"`
	OccurredAt  time.Time       `json:"occurred_at"`
	Version     int             `json:"version"`
}

// Stream defines an event stream for Kafka or similar.
type Stream interface {
	Publish(ctx context.Context, envelope Envelope) error
	Subscribe(ctx context.Context, topic string, handler func(Envelope) error) error
	Close() error
}

// Store persists events for event sourcing or outbox pattern.
type Store interface {
	Save(ctx context.Context, events []Envelope) error
	Load(ctx context.Context, aggregateID string, afterVersion int) ([]Envelope, error)
	All(ctx context.Context, afterPosition int64) ([]Envelope, error)
}
