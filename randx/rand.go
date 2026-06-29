// Package randx provides random utilities.
package randx

import (
	"math/rand/v2"

	"github.com/google/uuid"
)

const (
	alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	alpha    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lower    = "abcdefghijklmnopqrstuvwxyz"
	upper    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits   = "0123456789"
	hex      = "0123456789abcdef"
)

// String returns a random alphanumeric string of length n.
func String(n int) string {
	return StringWithCharset(n, alphanum)
}

// LowerString returns a random lowercase string of length n.
func LowerString(n int) string {
	return StringWithCharset(n, lower)
}

// UpperString returns a random uppercase string of length n.
func UpperString(n int) string {
	return StringWithCharset(n, upper)
}

// DigitString returns a random numeric string of length n.
func DigitString(n int) string {
	return StringWithCharset(n, digits)
}

// HexString returns a random lowercase hex string of length n.
func HexString(n int) string {
	return StringWithCharset(n, hex)
}

// StringWithCharset returns a random string of length n using the given charset.
func StringWithCharset(n int, charset string) string {
	if n <= 0 || len(charset) == 0 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

// UUID returns a new random UUID v4 string.
func UUID() string {
	return uuid.NewString()
}

// MustUUID returns a new random UUID v4 string, panicking on error.
// uuid.NewString() never errors in practice, so this is safe.
func MustUUID() string {
	return uuid.New().String()
}

// Int returns a random int in [min, max).
func Int(min, max int) int {
	if min >= max {
		return min
	}
	return min + rand.IntN(max-min)
}

// Int64 returns a random int64 in [min, max).
func Int64(min, max int64) int64 {
	if min >= max {
		return min
	}
	return min + rand.Int64N(max-min)
}

// Float64 returns a random float64 in [min, max).
func Float64(min, max float64) float64 {
	if min >= max {
		return min
	}
	return min + rand.Float64()*(max-min)
}

// Bool returns a random boolean.
func Bool() bool {
	return rand.IntN(2) == 0
}

// Bytes returns n random bytes.
func Bytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rand.IntN(256))
	}
	return b
}

// Shuffle shuffles a copy of s and returns it.
func Shuffle[T any](s []T) []T {
	result := make([]T, len(s))
	copy(result, s)
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return result
}

// Sample returns a random element from s.
func Sample[T any](s []T) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	return s[rand.IntN(len(s))], true
}

// Samples returns n distinct random elements from s (without replacement).
func Samples[T any](s []T, n int) []T {
	shuffled := Shuffle(s)
	if n > len(shuffled) {
		n = len(shuffled)
	}
	return shuffled[:n]
}

// WeightedChoice picks a random index based on weights.
// weights must have the same length as items; returns -1 if empty.
func WeightedChoice(weights []float64) int {
	if len(weights) == 0 {
		return -1
	}
	total := 0.0
	for _, w := range weights {
		total += w
	}
	r := rand.Float64() * total
	for i, w := range weights {
		r -= w
		if r <= 0 {
			return i
		}
	}
	return len(weights) - 1
}
