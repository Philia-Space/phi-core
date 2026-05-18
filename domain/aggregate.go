package domain

import (
	"errors"
)

// ValueObject represents an immutable object defined by its attributes.
type ValueObject interface {
	Equals(other ValueObject) bool
}

// AggregateRoot is the base for aggregate roots that manage consistency boundaries.
type AggregateRoot struct {
	*BaseEntity
	domainEvents []DomainEvent
}

func NewAggregateRoot(id string) *AggregateRoot {
	return &AggregateRoot{
		BaseEntity:   NewBaseEntity(id),
		domainEvents: make([]DomainEvent, 0),
	}
}

// DomainEvent returns and clears uncommitted events.
func (a *AggregateRoot) DomainEvents() []DomainEvent {
	events := a.domainEvents
	a.domainEvents = nil
	return events
}

// Raise adds a domain event to the aggregate.
func (a *AggregateRoot) Raise(event DomainEvent) {
	a.domainEvents = append(a.domainEvents, event)
	a.Touch()
}

// ErrAggregateNotFound is returned when an aggregate is not found.
var ErrAggregateNotFound = errors.New("aggregate not found")

// ErrAggregateConflict is returned when a concurrency conflict occurs.
var ErrAggregateConflict = errors.New("aggregate conflict: version mismatch")
