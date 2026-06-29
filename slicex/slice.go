// Package slicex provides generic slice utilities.
// Inspired by samber/lo, duke-git/lancet, and Java's Guava.
package slicex

import (
	"cmp"
	"math/rand/v2"
)

// Number is a type constraint for numeric types.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Map applies fn to each element and returns a new slice.
func Map[T, R any](s []T, fn func(T) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// MapWithIndex applies fn (with index) to each element and returns a new slice.
func MapWithIndex[T, R any](s []T, fn func(int, T) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = fn(i, v)
	}
	return result
}

// Filter returns elements for which fn returns true.
func Filter[T any](s []T, fn func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// FilterMap applies fn to each element; non-zero second returns are collected.
func FilterMap[T, R any](s []T, fn func(T) (R, bool)) []R {
	result := make([]R, 0)
	for _, v := range s {
		if r, ok := fn(v); ok {
			result = append(result, r)
		}
	}
	return result
}

// Reduce reduces a slice to a single value using fn.
func Reduce[T, R any](s []T, init R, fn func(R, T) R) R {
	acc := init
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

// ForEach calls fn for each element with its index.
func ForEach[T any](s []T, fn func(int, T)) {
	for i, v := range s {
		fn(i, v)
	}
}

// Contains reports whether s contains v.
func Contains[T comparable](s []T, v T) bool {
	for _, e := range s {
		if e == v {
			return true
		}
	}
	return false
}

// ContainsBy reports whether any element satisfies fn.
func ContainsBy[T any](s []T, fn func(T) bool) bool {
	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}

// Every reports whether all elements satisfy fn.
func Every[T any](s []T, fn func(T) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Some reports whether at least one element satisfies fn.
func Some[T any](s []T, fn func(T) bool) bool {
	return ContainsBy(s, fn)
}

// None reports whether no elements satisfy fn.
func None[T any](s []T, fn func(T) bool) bool {
	return !ContainsBy(s, fn)
}

// Count counts elements that satisfy fn.
func Count[T any](s []T, fn func(T) bool) int {
	n := 0
	for _, v := range s {
		if fn(v) {
			n++
		}
	}
	return n
}

// CountEqual counts occurrences of v.
func CountEqual[T comparable](s []T, v T) int {
	return Count(s, func(e T) bool { return e == v })
}

// Unique returns a slice with duplicates removed, preserving order.
func Unique[T comparable](s []T) []T {
	seen := make(map[T]struct{}, len(s))
	result := make([]T, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// UniqueBy returns a slice with duplicates removed based on the key returned by fn.
func UniqueBy[T any, K comparable](s []T, fn func(T) K) []T {
	seen := make(map[K]struct{}, len(s))
	result := make([]T, 0, len(s))
	for _, v := range s {
		k := fn(v)
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Flatten flattens a slice of slices into a single slice.
func Flatten[T any](s [][]T) []T {
	total := 0
	for _, sub := range s {
		total += len(sub)
	}
	result := make([]T, 0, total)
	for _, sub := range s {
		result = append(result, sub...)
	}
	return result
}

// FlatMap applies fn to each element and flattens the result.
func FlatMap[T, R any](s []T, fn func(T) []R) []R {
	result := make([]R, 0, len(s))
	for _, v := range s {
		result = append(result, fn(v)...)
	}
	return result
}

// Chunk splits s into chunks of size n.
// The last chunk may be smaller than n.
func Chunk[T any](s []T, n int) [][]T {
	if n <= 0 {
		return nil
	}
	result := make([][]T, 0, (len(s)+n-1)/n)
	for i := 0; i < len(s); i += n {
		end := i + n
		if end > len(s) {
			end = len(s)
		}
		chunk := make([]T, end-i)
		copy(chunk, s[i:end])
		result = append(result, chunk)
	}
	return result
}

// Zip combines two slices element-by-element into a slice of [2]any.
// The resulting slice has length min(len(a), len(b)).
// Use ZipWith for typed pairs.
func Zip[T, U any](a []T, b []U) [][2]any {
	n := min(len(a), len(b))
	result := make([][2]any, n)
	for i := 0; i < n; i++ {
		result[i] = [2]any{a[i], b[i]}
	}
	return result
}

// ZipWith combines two slices element-by-element using fn.
func ZipWith[T, U, R any](a []T, b []U, fn func(T, U) R) []R {
	n := min(len(a), len(b))
	result := make([]R, n)
	for i := 0; i < n; i++ {
		result[i] = fn(a[i], b[i])
	}
	return result
}

// GroupBy groups elements of s by the key returned by fn.
func GroupBy[T any, K comparable](s []T, fn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range s {
		k := fn(v)
		result[k] = append(result[k], v)
	}
	return result
}

// Partition splits s into two slices: elements matching fn and those not.
func Partition[T any](s []T, fn func(T) bool) (matching []T, rest []T) {
	for _, v := range s {
		if fn(v) {
			matching = append(matching, v)
		} else {
			rest = append(rest, v)
		}
	}
	return
}

// ToMap converts s to a map using fn to produce key-value pairs.
func ToMap[T any, K comparable, V any](s []T, fn func(T) (K, V)) map[K]V {
	result := make(map[K]V, len(s))
	for _, v := range s {
		k, val := fn(v)
		result[k] = val
	}
	return result
}

// IndexOf returns the index of the first occurrence of v, or -1.
func IndexOf[T comparable](s []T, v T) int {
	for i, e := range s {
		if e == v {
			return i
		}
	}
	return -1
}

// IndexBy returns the index of the first element matching fn, or -1.
func IndexBy[T any](s []T, fn func(T) bool) int {
	for i, v := range s {
		if fn(v) {
			return i
		}
	}
	return -1
}

// First returns the first element of s, or the zero value and false.
func First[T any](s []T) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	return s[0], true
}

// Last returns the last element of s, or the zero value and false.
func Last[T any](s []T) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	return s[len(s)-1], true
}

// FirstBy returns the first element matching fn, or the zero value and false.
func FirstBy[T any](s []T, fn func(T) bool) (T, bool) {
	for _, v := range s {
		if fn(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// LastBy returns the last element matching fn, or the zero value and false.
func LastBy[T any](s []T, fn func(T) bool) (T, bool) {
	for i := len(s) - 1; i >= 0; i-- {
		if fn(s[i]) {
			return s[i], true
		}
	}
	var zero T
	return zero, false
}

// Reverse returns a reversed copy of s.
func Reverse[T any](s []T) []T {
	result := make([]T, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}
	return result
}

// Clone returns a shallow copy of s.
func Clone[T any](s []T) []T {
	result := make([]T, len(s))
	copy(result, s)
	return result
}

// Equal reports whether a and b contain the same elements in the same order.
func Equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Compact removes zero values from s.
func Compact[T comparable](s []T) []T {
	var zero T
	return Filter(s, func(v T) bool { return v != zero })
}

// Intersection returns elements present in both a and b, preserving a's order.
func Intersection[T comparable](a, b []T) []T {
	set := make(map[T]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}
	result := make([]T, 0)
	seen := make(map[T]struct{})
	for _, v := range a {
		if _, inB := set[v]; inB {
			if _, already := seen[v]; !already {
				seen[v] = struct{}{}
				result = append(result, v)
			}
		}
	}
	return result
}

// Union returns all unique elements from all slices, in order of first appearance.
func Union[T comparable](slices ...[]T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0)
	for _, s := range slices {
		for _, v := range s {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				result = append(result, v)
			}
		}
	}
	return result
}

// Difference returns elements in a but not in b.
func Difference[T comparable](a, b []T) []T {
	set := make(map[T]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}
	result := make([]T, 0)
	seen := make(map[T]struct{})
	for _, v := range a {
		if _, inB := set[v]; !inB {
			if _, already := seen[v]; !already {
				seen[v] = struct{}{}
				result = append(result, v)
			}
		}
	}
	return result
}

// SymmetricDifference returns elements in a or b but not both.
func SymmetricDifference[T comparable](a, b []T) []T {
	return Union(Difference(a, b), Difference(b, a))
}

// Take returns the first n elements of s.
func Take[T any](s []T, n int) []T {
	if n <= 0 {
		return []T{}
	}
	if n >= len(s) {
		return Clone(s)
	}
	return Clone(s[:n])
}

// Drop removes the first n elements of s.
func Drop[T any](s []T, n int) []T {
	if n <= 0 {
		return Clone(s)
	}
	if n >= len(s) {
		return []T{}
	}
	return Clone(s[n:])
}

// TakeLast returns the last n elements of s.
func TakeLast[T any](s []T, n int) []T {
	if n <= 0 {
		return []T{}
	}
	if n >= len(s) {
		return Clone(s)
	}
	return Clone(s[len(s)-n:])
}

// DropLast removes the last n elements of s.
func DropLast[T any](s []T, n int) []T {
	if n <= 0 {
		return Clone(s)
	}
	if n >= len(s) {
		return []T{}
	}
	return Clone(s[:len(s)-n])
}

// TakeWhile takes elements from the beginning while fn returns true.
func TakeWhile[T any](s []T, fn func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range s {
		if !fn(v) {
			break
		}
		result = append(result, v)
	}
	return result
}

// DropWhile drops elements from the beginning while fn returns true.
func DropWhile[T any](s []T, fn func(T) bool) []T {
	i := 0
	for i < len(s) && fn(s[i]) {
		i++
	}
	return Clone(s[i:])
}

// Repeat returns a slice containing v repeated n times.
func Repeat[T any](v T, n int) []T {
	result := make([]T, n)
	for i := range result {
		result[i] = v
	}
	return result
}

// Fill fills s with val in place and returns s.
func Fill[T any](s []T, val T) []T {
	for i := range s {
		s[i] = val
	}
	return s
}

// Sum returns the sum of all elements.
func Sum[T Number](s []T) T {
	var total T
	for _, v := range s {
		total += v
	}
	return total
}

// Min returns the minimum element and true, or the zero value and false for empty slices.
func Min[T cmp.Ordered](s []T) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	m := s[0]
	for _, v := range s[1:] {
		if v < m {
			m = v
		}
	}
	return m, true
}

// Max returns the maximum element and true, or the zero value and false for empty slices.
func Max[T cmp.Ordered](s []T) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	m := s[0]
	for _, v := range s[1:] {
		if v > m {
			m = v
		}
	}
	return m, true
}

// MinBy returns the element for which fn returns the smallest value.
func MinBy[T any, K cmp.Ordered](s []T, fn func(T) K) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	best := s[0]
	bestKey := fn(best)
	for _, v := range s[1:] {
		if k := fn(v); k < bestKey {
			best = v
			bestKey = k
		}
	}
	return best, true
}

// MaxBy returns the element for which fn returns the largest value.
func MaxBy[T any, K cmp.Ordered](s []T, fn func(T) K) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	best := s[0]
	bestKey := fn(best)
	for _, v := range s[1:] {
		if k := fn(v); k > bestKey {
			best = v
			bestKey = k
		}
	}
	return best, true
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
// If n >= len(s), all elements are returned in random order.
func Samples[T any](s []T, n int) []T {
	if n <= 0 || len(s) == 0 {
		return []T{}
	}
	cp := Clone(s)
	rand.Shuffle(len(cp), func(i, j int) { cp[i], cp[j] = cp[j], cp[i] })
	if n > len(cp) {
		n = len(cp)
	}
	return cp[:n]
}

// Shuffle returns a randomly shuffled copy of s.
func Shuffle[T any](s []T) []T {
	result := Clone(s)
	rand.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })
	return result
}

// Pop removes and returns the last element of s.
func Pop[T any](s []T) (T, []T) {
	if len(s) == 0 {
		var zero T
		return zero, s
	}
	return s[len(s)-1], s[:len(s)-1]
}

// Keys extracts a key from each element using fn.
func Keys[T any, K comparable](s []T, fn func(T) K) []K {
	result := make([]K, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// Values extracts a value from each element using fn.
func Values[T, V any](s []T, fn func(T) V) []V {
	return Map(s, fn)
}

// Window returns all contiguous sub-slices of length n.
func Window[T any](s []T, n int) [][]T {
	if n <= 0 || n > len(s) {
		return nil
	}
	result := make([][]T, 0, len(s)-n+1)
	for i := 0; i <= len(s)-n; i++ {
		w := make([]T, n)
		copy(w, s[i:i+n])
		result = append(result, w)
	}
	return result
}

// Interleave merges multiple slices by taking one element from each in turn.
func Interleave[T any](slices ...[]T) []T {
	if len(slices) == 0 {
		return nil
	}
	maxLen := 0
	for _, s := range slices {
		if len(s) > maxLen {
			maxLen = len(s)
		}
	}
	result := make([]T, 0, maxLen*len(slices))
	for i := 0; i < maxLen; i++ {
		for _, s := range slices {
			if i < len(s) {
				result = append(result, s[i])
			}
		}
	}
	return result
}

// Associate converts s into a map using fn to produce key-value pairs.
// Same as ToMap; provided as an alias for readability.
func Associate[T any, K comparable, V any](s []T, fn func(T) (K, V)) map[K]V {
	return ToMap(s, fn)
}

// Pluck extracts a field from each element using fn (alias for Keys/Map).
func Pluck[T, V any](s []T, fn func(T) V) []V {
	return Map(s, fn)
}
