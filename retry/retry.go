// Package retry provides a simple retry mechanism with backoff for
// operations that may transiently fail (e.g. network scans, webhook calls).
package retry

import (
	"context"
	"errors"
	"time"
)

// Policy defines how retries are performed.
type Policy struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// Delay is the wait time between attempts.
	Delay time.Duration
	// Multiplier scales the delay after each failure (1.0 = constant).
	Multiplier float64
}

// Default returns a Policy with sensible defaults.
func Default() Policy {
	return Policy{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// ErrExhausted is returned when all attempts are consumed.
var ErrExhausted = errors.New("retry: all attempts exhausted")

// Do runs fn according to p, retrying on non-nil errors.
// It respects ctx cancellation between attempts.
func (p Policy) Do(ctx context.Context, fn func() error) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	if p.Multiplier <= 0 {
		p.Multiplier = 1.0
	}

	delay := p.Delay
	var last error

	for attempt := 0; attempt < p.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		last = fn()
		if last == nil {
			return nil
		}

		if attempt < p.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * p.Multiplier)
		}
	}

	return errors.Join(ErrExhausted, last)
}
