package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"portwatch/retry"
)

var errTemp = errors.New("temporary error")

func TestSucceedsOnFirstAttempt(t *testing.T) {
	p := retry.Default()
	calls := 0
	err := p.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetriesUpToMaxAttempts(t *testing.T) {
	p := retry.Policy{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1.0}
	calls := 0
	err := p.Do(context.Background(), func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, retry.ErrExhausted) {
		t.Fatalf("expected ErrExhausted, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestSucceedsOnSecondAttempt(t *testing.T) {
	p := retry.Policy{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1.0}
	calls := 0
	err := p.Do(context.Background(), func() error {
		calls++
		if calls < 2 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestRespectsContextCancellation(t *testing.T) {
	p := retry.Policy{MaxAttempts: 5, Delay: 100 * time.Millisecond, Multiplier: 1.0}
	ctx, cancel := context.WithCancel(context.Background())

	calls := 0
	err := p.Do(ctx, func() error {
		calls++
		if calls == 1 {
			cancel()
		}
		return errTemp
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestWrapsPreviousError(t *testing.T) {
	p := retry.Policy{MaxAttempts: 2, Delay: time.Millisecond, Multiplier: 1.0}
	err := p.Do(context.Background(), func() error { return errTemp })
	if !errors.Is(err, errTemp) {
		t.Fatalf("expected wrapped errTemp, got %v", err)
	}
}
