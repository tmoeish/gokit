// Package mathx provides generic math utilities.
package mathx

import (
	"cmp"
	"math"
)

// Number is a type constraint for all numeric types.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Signed is a type constraint for signed numeric types.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// Integer is a type constraint for integer types.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Clamp returns val clamped to [lo, hi].
func Clamp[T cmp.Ordered](val, lo, hi T) T {
	if val < lo {
		return lo
	}
	if val > hi {
		return hi
	}
	return val
}

// Abs returns the absolute value of v.
func Abs[T Signed](v T) T {
	if v < 0 {
		return -v
	}
	return v
}

// Sum returns the sum of all values.
func Sum[T Number](vals ...T) T {
	var total T
	for _, v := range vals {
		total += v
	}
	return total
}

// Average returns the average of all values as float64.
// Returns 0 for empty input.
func Average[T Number](vals ...T) float64 {
	if len(vals) == 0 {
		return 0
	}
	var total T
	for _, v := range vals {
		total += v
	}
	return float64(total) / float64(len(vals))
}

// Min2 returns the smaller of a and b.
func Min2[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max2 returns the larger of a and b.
func Max2[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// MinOf returns the smallest of all values.
func MinOf[T cmp.Ordered](vals ...T) T {
	if len(vals) == 0 {
		var zero T
		return zero
	}
	m := vals[0]
	for _, v := range vals[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

// MaxOf returns the largest of all values.
func MaxOf[T cmp.Ordered](vals ...T) T {
	if len(vals) == 0 {
		var zero T
		return zero
	}
	m := vals[0]
	for _, v := range vals[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

// InRange reports whether val is in [lo, hi].
func InRange[T cmp.Ordered](val, lo, hi T) bool {
	return val >= lo && val <= hi
}

// InRangeExclusive reports whether val is in (lo, hi).
func InRangeExclusive[T cmp.Ordered](val, lo, hi T) bool {
	return val > lo && val < hi
}

// IsBetween is an alias for InRange.
func IsBetween[T cmp.Ordered](val, lo, hi T) bool {
	return InRange(val, lo, hi)
}

// GCD returns the greatest common divisor of a and b.
func GCD[T Integer](a, b T) T {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM returns the least common multiple of a and b.
func LCM[T Integer](a, b T) T {
	return a / GCD(a, b) * b
}

// Pow returns base^exp for integer exponents (exp >= 0).
func Pow[T Integer](base, exp T) T {
	var result T = 1
	for exp > 0 {
		if exp%2 == 1 {
			result *= base
		}
		base *= base
		exp /= 2
	}
	return result
}

// Percent returns what percent n is of total.
// Returns 0 if total is 0.
func Percent[T Number](n, total T) float64 {
	if total == 0 {
		return 0
	}
	return float64(n) / float64(total) * 100
}

// RoundTo rounds f to decimalPlaces decimal places.
func RoundTo(f float64, decimalPlaces int) float64 {
	factor := math.Pow10(decimalPlaces)
	return math.Round(f*factor) / factor
}

// CeilTo rounds f up to decimalPlaces decimal places.
func CeilTo(f float64, decimalPlaces int) float64 {
	factor := math.Pow10(decimalPlaces)
	return math.Ceil(f*factor) / factor
}

// FloorTo rounds f down to decimalPlaces decimal places.
func FloorTo(f float64, decimalPlaces int) float64 {
	factor := math.Pow10(decimalPlaces)
	return math.Floor(f*factor) / factor
}

// IsEven reports whether n is even.
func IsEven[T Integer](n T) bool {
	return n%2 == 0
}

// IsOdd reports whether n is odd.
func IsOdd[T Integer](n T) bool {
	return n%2 != 0
}

// IsPrime reports whether n is a prime number.
func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// Fibonacci returns the nth Fibonacci number (0-indexed).
func Fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// SafeDivide returns a/b, or def if b is zero.
func SafeDivide[T Number](a, b T, def T) T {
	if b == 0 {
		return def
	}
	return a / b
}
