package schedule_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/schedule"
)

func TestRunesImmediately(t *testing.T) {
	var count int32ask := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := schedule.New(10*time.Second, task, nil)
	go r.Start(ctx)

	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&count) < 1 {
		t.Error("expected task to run at least once immediately")
	}
}

func TestRunnerTicksMultipleTimes(t *testing.T) {
	var count int32
	task := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := schedule.New(30*time.Millisecond, task, nil)
	go r.Start(ctx)

	time.Sleep(120 * time.Millisecond)
	cancel()

	if atomic.LoadInt32(&count) < 3 {
		t.Errorf("expected at least 3 executions, got %d", atomic.LoadInt32(&count))
	}
}

func TestRunnerCallsOnError(t *testing.T) {
	var errCount int32
	task := func(ctx context.Context) error {
		return errors.New("scan failed")
	}
	onError := func(err error) {
		atomic.AddInt32(&errCount, 1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := schedule.New(30*time.Millisecond, task, onError)
	go r.Start(ctx)

	time.Sleep(80 * time.Millisecond)
	cancel()

	if atomic.LoadInt32(&errCount) < 1 {
		t.Error("expected onError to be called at least once")
	}
}

func TestRunnerStopsOnContextCancel(t *testing.T) {
	var count int32
	task := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	r := schedule.New(20*time.Millisecond, task, nil)
	go r.Start(ctx)

	time.Sleep(70 * time.Millisecond)
	cancel()
	time.Sleep(50 * time.Millisecond)

	snap := atomic.LoadInt32(&count)
	time.Sleep(60 * time.Millisecond)
	if atomic.LoadInt32(&count) != snap {
		t.Error("expected task to stop running after context cancel")
	}
}
