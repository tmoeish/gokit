package logx_test

import (
	"context"
	"testing"

	"github.com/tmoeish/gokit/logx"
)

func TestLogging(t *testing.T) {
	ctx := context.Background()
	// Just verify no panics.
	logx.Debug(ctx, "debug message", "key", "value")
	logx.Info(ctx, "info message")
	logx.Warn(ctx, "warn message")
	logx.Error(ctx, "error message")
	logx.Debugf(ctx, "debug %s", "formatted")
	logx.Infof(ctx, "info %d", 42)
	logx.Warnf(ctx, "warn %.2f", 3.14)
	logx.Errorf(ctx, "error %v", "test")
}

func TestLogError(t *testing.T) {
	ctx := context.Background()
	logx.LogError(ctx, nil, "should not log") // nil error → no-op
	// Non-nil error should not panic.
	logx.LogError(ctx, &testErr{"boom"}, "got error")
}

type testErr struct{ msg string }

func (e *testErr) Error() string { return e.msg }

func TestCallerLoc(t *testing.T) {
	fn, loc := logx.CallerLoc()
	if fn == "unknown" || loc == "unknown" {
		t.Fatalf("CallerLoc: fn=%q loc=%q", fn, loc)
	}
}

func TestStack(t *testing.T) {
	stack := logx.Stack(1)
	if len(stack) == 0 {
		t.Fatal("Stack should not be empty")
	}
}

func TestConfigInit(t *testing.T) {
	cfg := logx.Config{Level: "debug", JSON: false}
	cfg.Init() // should not panic
	cfg2 := logx.Config{Level: "warn", JSON: true}
	cfg2.Init()
}
