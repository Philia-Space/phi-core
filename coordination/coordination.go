package coordination

import (
	"context"
	"sync"
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
	mu    sync.Mutex
	locks map[string]*memoryLock
}

type memoryLock struct {
	mu        sync.Mutex
	key       string
	expiresAt time.Time
}

func NewInMemoryLock() *InMemoryLock {
	return &InMemoryLock{
		locks: make(map[string]*memoryLock),
	}
}

func (l *InMemoryLock) Acquire(ctx context.Context, key string, ttl time.Duration) (Lock, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if existing, ok := l.locks[key]; ok {
		if time.Now().Before(existing.expiresAt) {
			return nil, ErrLockAlreadyHeld
		}
	}

	ml := &memoryLock{
		key:       key,
		expiresAt: time.Now().Add(ttl),
	}
	l.locks[key] = ml
	return ml, nil
}

func (l *memoryLock) Release(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.expiresAt = time.Time{}
	return nil
}

func (l *memoryLock) Extend(ctx context.Context, ttl time.Duration) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.expiresAt = time.Now().Add(ttl)
	return nil
}

// ErrLockAlreadyHeld is returned when a lock is already held.
var ErrLockAlreadyHeld = &LockError{message: "lock already held"}

// LockError represents a lock error.
type LockError struct {
	message string
}

func (e *LockError) Error() string { return e.message }
