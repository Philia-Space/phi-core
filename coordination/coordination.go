package coordination

import (
	"context"
	"time"
)

// DistributedLock provides a distributed locking mechanism.
type DistributedLock interface {
	Acquire(ctx context.Context, key string, ttl time.Duration) (Lock, error)
}

// Lock represents an acquired distributed lock.
type Lock interface {
	Release(ctx context.Context) error
	Extend(ctx context.Context, ttl time.Duration) error
}

// LeaderElection manages leader election among service instances.
type LeaderElection interface {
	Campaign(ctx context.Context) (context.Context, error)
	IsLeader() bool
	Resign(ctx context.Context) error
}

// InMemoryLock is a simple in-memory lock for single-instance dev.
type InMemoryLock struct {
	locks map[string]*memoryLock
}

type memoryLock struct {
	key       string
	expiresAt time.Time
}

func NewInMemoryLock() *InMemoryLock {
	return &InMemoryLock{
		locks: make(map[string]*memoryLock),
	}
}

func (l *InMemoryLock) Acquire(ctx context.Context, key string, ttl time.Duration) (Lock, error) {
	ml := &memoryLock{
		key:       key,
		expiresAt: time.Now().Add(ttl),
	}
	l.locks[key] = ml
	return ml, nil
}

func (l *memoryLock) Release(ctx context.Context) error {
	return nil
}

func (l *memoryLock) Extend(ctx context.Context, ttl time.Duration) error {
	l.expiresAt = time.Now().Add(ttl)
	return nil
}
