package domain

import (
	"context"
	"time"
)

// DomainEvent represents a fact that happened in the domain.
type DomainEvent interface {
	AggregateID() string
	EventType() string
	OccurredAt() time.Time
}

// BaseDomainEvent provides common event fields.
type BaseDomainEvent struct {
	aggregateID string
	eventType   string
	occurredAt  time.Time
}

func NewBaseDomainEvent(aggregateID, eventType string) BaseDomainEvent {
	return BaseDomainEvent{
		aggregateID: aggregateID,
		eventType:   eventType,
		occurredAt:  time.Now().UTC(),
	}
}

func (e BaseDomainEvent) AggregateID() string  { return e.aggregateID }
func (e BaseDomainEvent) EventType() string    { return e.eventType }
func (e BaseDomainEvent) OccurredAt() time.Time { return e.occurredAt }

// EventHandler handles a specific domain event.
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
	SupportedEvents() []string
}

// EventBus publishes and subscribes to domain events.
type EventBus interface {
	Publish(ctx context.Context, events []DomainEvent) error
	Subscribe(handler EventHandler) error
}
