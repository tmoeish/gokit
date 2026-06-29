// Package logx provides structured logging utilities built on top of log/slog.
package logx

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

// Config holds logging configuration.
type Config struct {
	// Level is the minimum log level: "debug", "info", "warn", "error".
	Level string `json:"level" yaml:"level"`
	// JSON enables JSON format output; text format is used otherwise.
	JSON bool `json:"json" yaml:"json"`
}

// Init initializes the default slog logger from cfg.
func (c *Config) Init() {
	level := parseLevel(c.Level)
	opts := &slog.HandlerOptions{Level: level}
	var handler slog.Handler
	if c.JSON {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}
	slog.SetDefault(slog.New(handler))
}

// NewLogger returns a named logger derived from slog.Default().
func NewLogger(module string) *slog.Logger {
	return slog.Default().With("module", module)
}

// FromContext returns the slog.Logger stored in ctx (via contextx), or slog.Default().
// To avoid a direct dependency on contextx, we use a context-key approach here.
func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return slog.Default()
	}
	type loggerKey struct{}
	if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok && l != nil {
		return l
	}
	return slog.Default()
}

// Debug logs at DEBUG level using the logger from ctx.
func Debug(ctx context.Context, msg string, args ...any) {
	logAt(ctx, slog.LevelDebug, msg, args...)
}

// Info logs at INFO level using the logger from ctx.
func Info(ctx context.Context, msg string, args ...any) {
	logAt(ctx, slog.LevelInfo, msg, args...)
}

// Warn logs at WARN level using the logger from ctx.
func Warn(ctx context.Context, msg string, args ...any) {
	logAt(ctx, slog.LevelWarn, msg, args...)
}

// Error logs at ERROR level using the logger from ctx.
func Error(ctx context.Context, msg string, args ...any) {
	logAt(ctx, slog.LevelError, msg, args...)
}

// Debugf logs a formatted message at DEBUG level.
func Debugf(ctx context.Context, format string, args ...any) {
	Debug(ctx, fmt.Sprintf(format, args...))
}

// Infof logs a formatted message at INFO level.
func Infof(ctx context.Context, format string, args ...any) {
	Info(ctx, fmt.Sprintf(format, args...))
}

// Warnf logs a formatted message at WARN level.
func Warnf(ctx context.Context, format string, args ...any) {
	Warn(ctx, fmt.Sprintf(format, args...))
}

// Errorf logs a formatted message at ERROR level.
func Errorf(ctx context.Context, format string, args ...any) {
	Error(ctx, fmt.Sprintf(format, args...))
}

// LogError logs err at ERROR level with stack information.
// Does nothing if err is nil.
func LogError(ctx context.Context, err error, msg string, args ...any) {
	if err == nil {
		return
	}
	fn, loc := callerLoc()
	slog.Default().Error(
		msg,
		append([]any{"error", err, "func", fn, "loc", loc}, args...)...,
	)
}

// CallerLoc returns the function name and "file:line" of the caller
// (skipping logx-internal frames).
func CallerLoc() (funcName, location string) {
	return callerLoc()
}

// Stack returns the call stack as a slice of "func file:line" strings.
func Stack(depth int) []string {
	var stack []string
	for i := depth; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		name := "unknown"
		if fn != nil {
			name = fn.Name()
		}
		stack = append(stack, fmt.Sprintf("%s %s:%d", name, file, line))
	}
	return stack
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func logAt(ctx context.Context, level slog.Level, msg string, args ...any) {
	fn, loc := callerLoc()
	slog.Default().Log(ctx, level, msg, append([]any{"func", fn, "loc", loc}, args...)...)
}

func callerLoc() (string, string) {
	pcs := make([]uintptr, 16)
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.Function, "github.com/tmoeish/gokit/logx.") {
			fn := frame.Function
			if i := strings.LastIndex(fn, "/"); i >= 0 {
				fn = fn[i+1:]
			}
			return fn, fmt.Sprintf("%s:%d", frame.File, frame.Line)
		}
		if !more {
			break
		}
	}
	return "unknown", "unknown"
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
