package consistency

import (
	"context"
)

// OutboxStore persists events to an outbox table for reliable delivery.
type OutboxStore interface {
	Save(ctx context.Context, eventType, aggregateID string, payload []byte) error
	Pending(ctx context.Context, limit int) ([]OutboxEntry, error)
	MarkProcessed(ctx context.Context, ids []string) error
}

// OutboxEntry represents a pending outbox event.
type OutboxEntry struct {
	ID          string
	EventType   string
	AggregateID string
	Payload     []byte
	CreatedAt   string
}

// OutboxProcessor publishes outbox events and marks them processed.
type OutboxProcessor struct {
	store   OutboxStore
	publisher func(ctx context.Context, entry OutboxEntry) error
	batchSize int
}

// NewOutboxProcessor creates a new outbox processor.
func NewOutboxProcessor(store OutboxStore, publisher func(ctx context.Context, entry OutboxEntry) error, batchSize int) *OutboxProcessor {
	return &OutboxProcessor{
		store:     store,
		publisher: publisher,
		batchSize: batchSize,
	}
}

// ProcessPending processes pending outbox events.
func (p *OutboxProcessor) ProcessPending(ctx context.Context) error {
	entries, err := p.store.Pending(ctx, p.batchSize)
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(entries))
	for _, entry := range entries {
		if err := p.publisher(ctx, entry); err != nil {
			continue
		}
		ids = append(ids, entry.ID)
	}

	if len(ids) > 0 {
		return p.store.MarkProcessed(ctx, ids)
	}
	return nil
}

// SagaStep defines a step in a saga.
type SagaStep struct {
	Name   string
	Action func(context.Context) error
	Compensate func(context.Context) error
}

// Saga executes a series of steps with compensation on failure.
type Saga struct {
	steps []SagaStep
}

// NewSaga creates a new saga.
func NewSaga(steps []SagaStep) *Saga {
	return &Saga{steps: steps}
}

// Execute runs all steps, compensating on failure.
func (s *Saga) Execute(ctx context.Context) error {
	completed := 0

	for i, step := range s.steps {
		if err := step.Action(ctx); err != nil {
			// Compensate completed steps in reverse
			for j := i - 1; j >= 0; j-- {
				_ = s.steps[j].Compensate(ctx)
			}
			return err
		}
		completed = i + 1
	}

	return nil
}
