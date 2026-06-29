// Package ptrx provides pointer utilities for working with Go generics.
package ptrx

// Of returns a pointer to v.
// Useful when you need a pointer to a literal or computed value.
//
//	ptrx.Of(42)      // *int
//	ptrx.Of("hello") // *string
func Of[T any](v T) *T {
	return &v
}

// Deref dereferences p. If p is nil, def is returned.
func Deref[T any](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}

// DerefZero dereferences p. If p is nil, the zero value of T is returned.
func DerefZero[T any](p *T) T {
	var zero T
	if p == nil {
		return zero
	}
	return *p
}

// IsNil reports whether p is nil.
func IsNil[T any](p *T) bool {
	return p == nil
}

// IsNotNil reports whether p is not nil.
func IsNotNil[T any](p *T) bool {
	return p != nil
}

// IfNil returns alt if p is nil, otherwise p.
func IfNil[T any](p *T, alt *T) *T {
	if p == nil {
		return alt
	}
	return p
}

// CoalescePtr returns the first non-nil pointer.
func CoalescePtr[T any](ptrs ...*T) *T {
	for _, p := range ptrs {
		if p != nil {
			return p
		}
	}
	return nil
}

// ToSlice converts a pointer to a single-element slice, or nil if the pointer is nil.
func ToSlice[T any](p *T) []T {
	if p == nil {
		return nil
	}
	return []T{*p}
}

// FromSlice returns a pointer to the first element of s, or nil if s is empty.
func FromSlice[T any](s []T) *T {
	if len(s) == 0 {
		return nil
	}
	return &s[0]
}
