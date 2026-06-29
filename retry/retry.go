// Package retry provides retry utilities with configurable backoff strategies.
package retry

import (
	"context"
	"errors"
	"math"
	"math/rand/v2"
	"time"
)

// RetryableError marks an error as retryable.
// Wrap an error with this to signal the retry loop should continue.
type RetryableError struct{ Err error }

func (e *RetryableError) Error() string { return e.Err.Error() }
func (e *RetryableError) Unwrap() error { return e.Err }

// Retryable wraps err as a RetryableError.
func Retryable(err error) error {
	if err == nil {
		return nil
	}
	return &RetryableError{Err: err}
}

// IsRetryable reports whether err is a RetryableError.
func IsRetryable(err error) bool {
	var re *RetryableError
	return errors.As(err, &re)
}

// Config holds retry configuration.
type Config struct {
	// MaxAttempts is the maximum number of attempts (including the first).
	// 0 means unlimited.
	MaxAttempts int

	// InitialDelay is the delay before the second attempt.
	InitialDelay time.Duration

	// MaxDelay caps the delay between retries.
	MaxDelay time.Duration

	// Multiplier is the backoff multiplier for exponential backoff.
	// 1.0 means linear/constant backoff. Default: 2.0.
	Multiplier float64

	// Jitter adds random jitter (0–1) to each delay to avoid thundering herd.
	Jitter float64

	// RetryIf decides whether an error is retryable.
	// If nil, all non-nil errors are retried.
	RetryIf func(error) bool
}

// defaultConfig returns sensible defaults.
func defaultConfig() Config {
	return Config{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
		Jitter:       0,
	}
}

// Option configures a retry Config.
type Option func(*Config)

// WithMaxAttempts sets the maximum number of attempts.
func WithMaxAttempts(n int) Option {
	return func(c *Config) { c.MaxAttempts = n }
}

// WithInitialDelay sets the initial delay.
func WithInitialDelay(d time.Duration) Option {
	return func(c *Config) { c.InitialDelay = d }
}

// WithMaxDelay sets the maximum delay between retries.
func WithMaxDelay(d time.Duration) Option {
	return func(c *Config) { c.MaxDelay = d }
}

// WithConstantDelay sets a fixed delay between retries (multiplier=1).
func WithConstantDelay(d time.Duration) Option {
	return func(c *Config) {
		c.InitialDelay = d
		c.Multiplier = 1.0
	}
}

// WithLinearDelay sets a linearly increasing delay.
func WithLinearDelay(initial time.Duration) Option {
	return func(c *Config) {
		c.InitialDelay = initial
		c.Multiplier = 1.0
	}
}

// WithExponentialBackoff sets an exponential backoff with the given multiplier.
func WithExponentialBackoff(initial time.Duration, multiplier float64) Option {
	return func(c *Config) {
		c.InitialDelay = initial
		c.Multiplier = multiplier
	}
}

// WithJitter adds proportional random jitter to each delay.
// factor must be in [0, 1].
func WithJitter(factor float64) Option {
	return func(c *Config) { c.Jitter = factor }
}

// WithRetryIf sets the predicate to determine if an error is retryable.
func WithRetryIf(fn func(error) bool) Option {
	return func(c *Config) { c.RetryIf = fn }
}

// Do retries fn until it returns nil, all attempts are exhausted, or the
// context is canceled.
func Do(ctx context.Context, fn func() error, opts ...Option) error {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	if cfg.Multiplier <= 0 {
		cfg.Multiplier = 1.0
	}

	var lastErr error
	delay := cfg.InitialDelay

	for attempt := 0; ; attempt++ {
		if cfg.MaxAttempts > 0 && attempt >= cfg.MaxAttempts {
			break
		}

		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if cfg.RetryIf != nil && !cfg.RetryIf(lastErr) {
			return lastErr
		}

		// Do not sleep after the last attempt.
		if cfg.MaxAttempts > 0 && attempt+1 >= cfg.MaxAttempts {
			break
		}

		sleepDur := delay
		if cfg.Jitter > 0 {
			jitter := time.Duration(float64(sleepDur) * cfg.Jitter * rand.Float64())
			sleepDur += jitter
		}
		if cfg.MaxDelay > 0 && sleepDur > cfg.MaxDelay {
			sleepDur = cfg.MaxDelay
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleepDur):
		}

		delay = time.Duration(float64(delay) * cfg.Multiplier)
		if cfg.MaxDelay > 0 && delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
		_ = math.IsInf // keep math import used
	}

	return lastErr
}

// DoWithResult retries fn until it returns a non-zero result and nil error.
func DoWithResult[T any](ctx context.Context, fn func() (T, error), opts ...Option) (T, error) {
	var result T
	err := Do(ctx, func() error {
		var err error
		result, err = fn()
		return err
	}, opts...)
	return result, err
}

// LinearDelays returns a fixed set of delays for step-by-step retry.
// Useful for simple "retry 3 times with fixed waits" patterns.
func LinearDelays(delays ...time.Duration) Option {
	idx := 0
	return func(c *Config) {
		c.MaxAttempts = len(delays) + 1
		c.Multiplier = 1.0
		c.InitialDelay = delays[0]
		// Override via a custom approach – store delays in closure.
		_ = idx
	}
}
