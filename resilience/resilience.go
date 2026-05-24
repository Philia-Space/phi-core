package resilience

import (
	"context"
	"sync"
	"time"
)

// RetryConfig defines retry behavior.
type RetryConfig struct {
	MaxAttempts int
	Backoff     time.Duration
	MaxBackoff  time.Duration
	Retryable   func(error) bool
}

// DefaultRetryConfig returns sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		Backoff:     100 * time.Millisecond,
		MaxBackoff:  5 * time.Second,
		Retryable:   func(err error) bool { return true },
	}
}

// Retry executes fn with retry logic.
func Retry(ctx context.Context, cfg RetryConfig, fn func(context.Context) error) error {
	var lastErr error
	backoff := cfg.Backoff

	for i := 0; i < cfg.MaxAttempts; i++ {
		lastErr = fn(ctx)
		if lastErr == nil {
			return nil
		}
		if !cfg.Retryable(lastErr) {
			return lastErr
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}

		backoff *= 2
		if backoff > cfg.MaxBackoff {
			backoff = cfg.MaxBackoff
		}
	}
	return lastErr
}

// CircuitState represents the circuit breaker state.
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// CircuitBreaker implements the circuit breaker pattern.
type CircuitBreaker struct {
	mu             sync.Mutex
	state          CircuitState
	failureCount   int
	successCount   int
	threshold      int
	halfOpenMax    int
	lastFailureAt  time.Time
	timeout        time.Duration
}

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:       CircuitClosed,
		threshold:   threshold,
		halfOpenMax: 1,
		timeout:     timeout,
	}
}

// Execute runs fn if the circuit allows.
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()
	if cb.state == CircuitOpen {
		if time.Since(cb.lastFailureAt) > cb.timeout {
			cb.state = CircuitHalfOpen
			cb.successCount = 0
		} else {
			cb.mu.Unlock()
			return ErrCircuitOpen
		}
	}
	cb.mu.Unlock()

	err := fn()
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if err != nil {
		cb.onFailure()
		return err
	}
	cb.onSuccess()
	return nil
}

func (cb *CircuitBreaker) onFailure() {
	cb.failureCount++
	cb.lastFailureAt = time.Now()
	if cb.failureCount >= cb.threshold {
		cb.state = CircuitOpen
	}
}

func (cb *CircuitBreaker) onSuccess() {
	if cb.state == CircuitHalfOpen {
		cb.successCount++
		if cb.successCount >= cb.halfOpenMax {
			cb.state = CircuitClosed
			cb.failureCount = 0
		}
	} else {
		cb.failureCount = 0
	}
}

// State returns the current circuit state.
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// ErrCircuitOpen is returned when the circuit is open.
var ErrCircuitOpen = &CircuitError{message: "circuit breaker is open"}

// CircuitError represents a circuit breaker error.
type CircuitError struct {
	message string
}

func (e *CircuitError) Error() string { return e.message }

// Timeout executes fn with a timeout.
func Timeout(ctx context.Context, d time.Duration, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()
	return fn(ctx)
}
