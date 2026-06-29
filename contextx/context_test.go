package contextx_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/tmoeish/gokit/contextx"
)

func TestWithValue(t *testing.T) {
	ctx := context.Background()
	ctx = contextx.WithValue[int](ctx, "count", 42)
	v, ok := contextx.Value[int](ctx, "count")
	if !ok || v != 42 {
		t.Fatalf("WithValue/Value: %d %v", v, ok)
	}
}

func TestValueOrDefault(t *testing.T) {
	ctx := context.Background()
	v := contextx.ValueOrDefault[string](ctx, "missing", "default")
	if v != "default" {
		t.Fatalf("ValueOrDefault: %s", v)
	}
}

func TestReqID(t *testing.T) {
	ctx := contextx.WithReqID(context.Background(), "req-123")
	if contextx.ReqID(ctx) != "req-123" {
		t.Fatal("ReqID mismatch")
	}
}

func TestTraceID(t *testing.T) {
	ctx := contextx.WithTraceID(context.Background(), "trace-abc")
	if contextx.TraceID(ctx) != "trace-abc" {
		t.Fatal("TraceID mismatch")
	}
}

func TestRealIP(t *testing.T) {
	ctx := contextx.WithRealIP(context.Background(), "1.2.3.4")
	if contextx.RealIP(ctx) != "1.2.3.4" {
		t.Fatal("RealIP mismatch")
	}
}

func TestLogger(t *testing.T) {
	ctx := context.Background()
	l := contextx.Logger(ctx)
	if l == nil {
		t.Fatal("Logger should never be nil")
	}
	custom := slog.New(slog.Default().Handler())
	ctx = contextx.WithLogger(ctx, custom)
	if contextx.Logger(ctx) != custom {
		t.Fatal("WithLogger/Logger mismatch")
	}
}

func TestWithLogAttrs(t *testing.T) {
	ctx := context.Background()
	ctx = contextx.WithLogAttrs(ctx, "key", "value")
	// Should not panic; logger should have the attr attached.
	if contextx.Logger(ctx) == nil {
		t.Fatal("WithLogAttrs: nil logger")
	}
}
