package event

import (
	"context"
	"sync"
)

// InMemoryBus is a simple in-memory event bus for testing and local dev.
type InMemoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]func(context.Context, Envelope) error
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers: make(map[string][]func(context.Context, Envelope) error),
	}
}

func (b *InMemoryBus) Publish(ctx context.Context, envelope Envelope) error {
	b.mu.RLock()
	handlers := b.handlers[envelope.EventType]
	b.mu.RUnlock()

	for _, h := range handlers {
		if err := h(ctx, envelope); err != nil {
			return err
		}
	}
	return nil
}

func (b *InMemoryBus) Subscribe(eventType string, handler func(context.Context, Envelope) error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *InMemoryBus) Close() error { return nil }
