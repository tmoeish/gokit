// Package must provides must-succeed helpers that panic on error.
// Use only in initialization code, tests, or when errors are truly impossible.
package must

import "fmt"

// Must returns v or panics if err is non-nil.
//
//	f, err := os.Open(path)
//	// becomes:
//	f := must.Must(os.Open(path))
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// Must2 returns v1, v2 or panics if err is non-nil.
func Must2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}
	return v1, v2
}

// OK panics if err is non-nil (no return value).
func OK(err error) {
	if err != nil {
		panic(err)
	}
}

// Assert panics with msg if cond is false.
func Assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

// Assertf panics with a formatted message if cond is false.
func Assertf(cond bool, format string, args ...any) {
	if !cond {
		panic(fmt.Sprintf(format, args...))
	}
}
