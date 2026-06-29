// Package contextx provides context utilities for attaching and retrieving
// typed values from a context.Context.
package contextx

import (
	"context"
	"log/slog"
)

type key[T any] struct{ name string }

// WithValue stores a typed value in ctx under a type-safe key.
func WithValue[T any](ctx context.Context, name string, val T) context.Context {
	return context.WithValue(ctx, key[T]{name}, val)
}

// Value retrieves a typed value from ctx. Returns the zero value and false if absent.
func Value[T any](ctx context.Context, name string) (T, bool) {
	v, ok := ctx.Value(key[T]{name}).(T)
	return v, ok
}

// ValueOrDefault retrieves a typed value from ctx, or returns def if absent.
func ValueOrDefault[T any](ctx context.Context, name string, def T) T {
	if v, ok := Value[T](ctx, name); ok {
		return v
	}
	return def
}

// ────────────────────────────────────────────────
// Package-defined context keys for common request metadata
// ────────────────────────────────────────────────

type ctxKey string

const (
	keyReqID   ctxKey = "req_id"
	keyTraceID ctxKey = "trace_id"
	keyRealIP  ctxKey = "real_ip"
	keyLogger  ctxKey = "logger"
)

// WithReqID attaches a request ID to ctx and enriches the logger.
func WithReqID(ctx context.Context, reqID string) context.Context {
	ctx = WithLogAttrs(ctx, "req_id", reqID)
	return context.WithValue(ctx, keyReqID, reqID)
}

// ReqID retrieves the request ID from ctx.
func ReqID(ctx context.Context) string {
	if s, ok := ctx.Value(keyReqID).(string); ok {
		return s
	}
	return ""
}

// WithTraceID attaches a trace ID to ctx and enriches the logger.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	ctx = WithLogAttrs(ctx, "trace_id", traceID)
	return context.WithValue(ctx, keyTraceID, traceID)
}

// TraceID retrieves the trace ID from ctx.
func TraceID(ctx context.Context) string {
	if s, ok := ctx.Value(keyTraceID).(string); ok {
		return s
	}
	return ""
}

// WithRealIP attaches the real client IP to ctx.
func WithRealIP(ctx context.Context, ip string) context.Context {
	ctx = WithLogAttrs(ctx, "real_ip", ip)
	return context.WithValue(ctx, keyRealIP, ip)
}

// RealIP retrieves the real client IP from ctx.
func RealIP(ctx context.Context) string {
	if s, ok := ctx.Value(keyRealIP).(string); ok {
		return s
	}
	return ""
}

// WithLogger attaches a slog.Logger to ctx.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	if logger == nil {
		return ctx
	}
	return context.WithValue(ctx, keyLogger, logger)
}

// Logger retrieves the slog.Logger from ctx.
// Falls back to slog.Default() if none is set.
func Logger(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(keyLogger).(*slog.Logger); ok && l != nil {
		return l
	}
	return slog.Default()
}

// WithLogAttrs enriches the context logger with additional key-value attrs.
func WithLogAttrs(ctx context.Context, attrs ...any) context.Context {
	return context.WithValue(ctx, keyLogger, Logger(ctx).With(attrs...))
}
