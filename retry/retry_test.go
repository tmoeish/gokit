package retry_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/tmoeish/gokit/retry"
)

func TestRetrySuccess(t *testing.T) {
	var count int32
	err := retry.Do(context.Background(), func() error {
		atomic.AddInt32(&count, 1)
		return nil
	}, retry.WithMaxAttempts(3))
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 attempt, got %d", count)
	}
}

func TestRetryExhausted(t *testing.T) {
	var count int32
	err := retry.Do(context.Background(), func() error {
		atomic.AddInt32(&count, 1)
		return errors.New("fail")
	}, retry.WithMaxAttempts(3), retry.WithConstantDelay(time.Millisecond))
	if err == nil {
		t.Fatal("expected error")
	}
	if count != 3 {
		t.Fatalf("expected 3 attempts, got %d", count)
	}
}

func TestRetrySucceedsOnThirdAttempt(t *testing.T) {
	var count int32
	err := retry.Do(context.Background(), func() error {
		n := atomic.AddInt32(&count, 1)
		if n < 3 {
			return errors.New("not yet")
		}
		return nil
	}, retry.WithMaxAttempts(5), retry.WithConstantDelay(time.Millisecond))
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if count != 3 {
		t.Fatalf("expected 3 attempts, got %d", count)
	}
}

func TestRetryContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := retry.Do(ctx, func() error {
		return errors.New("fail")
	}, retry.WithMaxAttempts(10))
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestRetryWithRetryIf(t *testing.T) {
	sentinel := errors.New("not retryable")
	var count int32
	err := retry.Do(context.Background(), func() error {
		atomic.AddInt32(&count, 1)
		return sentinel
	},
		retry.WithMaxAttempts(5),
		retry.WithConstantDelay(time.Millisecond),
		retry.WithRetryIf(func(err error) bool { return !errors.Is(err, sentinel) }),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if count != 1 {
		t.Fatalf("expected 1 attempt (non-retryable), got %d", count)
	}
}

func TestDoWithResult(t *testing.T) {
	var count int32
	result, err := retry.DoWithResult(context.Background(), func() (int, error) {
		n := atomic.AddInt32(&count, 1)
		if n < 3 {
			return 0, errors.New("not yet")
		}
		return 42, nil
	}, retry.WithMaxAttempts(5), retry.WithConstantDelay(time.Millisecond))
	if err != nil || result != 42 {
		t.Fatalf("DoWithResult: result=%d err=%v", result, err)
	}
}
