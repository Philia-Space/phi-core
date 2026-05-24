package domain

import (
	"errors"
)

// ValueObject represents an immutable object defined by its attributes.
type ValueObject interface {
	Equals(other ValueObject) bool
}

// AggregateRoot is the base for aggregate roots that manage consistency boundaries.
//
// Version is used for optimistic concurrency control: repository implementations
// (e.g. Mongo) should check Version before persisting and increment it on save.
type AggregateRoot struct {
	*BaseEntity
	Version      int `json:"version"`
	domainEvents []DomainEvent
}

func NewAggregateRoot(id string) *AggregateRoot {
	return &AggregateRoot{
		BaseEntity:   NewBaseEntity(id),
		Version:      0,
		domainEvents: make([]DomainEvent, 0),
	}
}

// DomainEvents returns and clears uncommitted events.
func (a *AggregateRoot) DomainEvents() []DomainEvent {
	events := a.domainEvents
	a.domainEvents = nil
	return events
}

// PeekDomainEvents returns uncommitted events WITHOUT clearing them.
// Useful for tests or for repositories that want to inspect events before
// committing.
func (a *AggregateRoot) PeekDomainEvents() []DomainEvent {
	out := make([]DomainEvent, len(a.domainEvents))
	copy(out, a.domainEvents)
	return out
}

// Raise adds a domain event to the aggregate.
func (a *AggregateRoot) Raise(event DomainEvent) {
	a.domainEvents = append(a.domainEvents, event)
	a.Touch()
}

// IncrementVersion bumps the aggregate version. Called by repository
// implementations on successful save.
func (a *AggregateRoot) IncrementVersion() {
	a.Version++
}

// ErrAggregateNotFound is returned when an aggregate is not found.
var ErrAggregateNotFound = errors.New("aggregate not found")

// ErrAggregateConflict is returned when a concurrency conflict occurs.
var ErrAggregateConflict = errors.New("aggregate conflict: version mismatch")
