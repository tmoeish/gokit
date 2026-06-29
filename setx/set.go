// Package setx provides a generic, hash-based Set collection.
// Inspired by Java Guava's Sets and duke-git/lancet, filling the gap left by
// Go's lack of a built-in set type.
//
// Set is not safe for concurrent use; guard it with a sync.Mutex or use
// syncx.SafeMap if concurrent access is required.
package setx

// Set is a collection of unique elements.
type Set[T comparable] struct {
	m map[T]struct{}
}

// New returns a new Set containing the given elements.
func New[T comparable](elems ...T) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{}, len(elems))}
	s.Add(elems...)
	return s
}

// FromSlice returns a new Set containing the elements of s.
func FromSlice[T comparable](s []T) *Set[T] {
	return New(s...)
}

// Add inserts elems into the set.
func (s *Set[T]) Add(elems ...T) {
	for _, e := range elems {
		s.m[e] = struct{}{}
	}
}

// Remove deletes elems from the set. Missing elements are ignored.
func (s *Set[T]) Remove(elems ...T) {
	for _, e := range elems {
		delete(s.m, e)
	}
}

// Contains reports whether e is in the set.
func (s *Set[T]) Contains(e T) bool {
	_, ok := s.m[e]
	return ok
}

// ContainsAll reports whether every element of elems is in the set.
func (s *Set[T]) ContainsAll(elems ...T) bool {
	for _, e := range elems {
		if !s.Contains(e) {
			return false
		}
	}
	return true
}

// ContainsAny reports whether at least one element of elems is in the set.
func (s *Set[T]) ContainsAny(elems ...T) bool {
	for _, e := range elems {
		if s.Contains(e) {
			return true
		}
	}
	return false
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int { return len(s.m) }

// IsEmpty reports whether the set has no elements.
func (s *Set[T]) IsEmpty() bool { return len(s.m) == 0 }

// Clear removes all elements from the set.
func (s *Set[T]) Clear() { s.m = make(map[T]struct{}) }

// Clone returns a shallow copy of the set.
func (s *Set[T]) Clone() *Set[T] {
	out := &Set[T]{m: make(map[T]struct{}, len(s.m))}
	for e := range s.m {
		out.m[e] = struct{}{}
	}
	return out
}

// ToSlice returns the set's elements as a slice in unspecified order.
func (s *Set[T]) ToSlice() []T {
	out := make([]T, 0, len(s.m))
	for e := range s.m {
		out = append(out, e)
	}
	return out
}

// ForEach calls fn for each element in unspecified order.
func (s *Set[T]) ForEach(fn func(T)) {
	for e := range s.m {
		fn(e)
	}
}

// Equal reports whether s and other contain exactly the same elements.
func (s *Set[T]) Equal(other *Set[T]) bool {
	if len(s.m) != len(other.m) {
		return false
	}
	return s.ContainsAll(other.ToSlice()...)
}

// Union returns a new set containing elements in either s or other.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	out := s.Clone()
	out.Add(other.ToSlice()...)
	return out
}

// Intersection returns a new set containing elements in both s and other.
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	out := New[T]()
	// Iterate over the smaller set for efficiency.
	a, b := s, other
	if b.Len() < a.Len() {
		a, b = b, a
	}
	for e := range a.m {
		if b.Contains(e) {
			out.Add(e)
		}
	}
	return out
}

// Difference returns a new set containing elements in s but not in other.
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	out := New[T]()
	for e := range s.m {
		if !other.Contains(e) {
			out.Add(e)
		}
	}
	return out
}

// SymmetricDifference returns a new set containing elements in exactly one of
// s or other.
func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T] {
	return s.Difference(other).Union(other.Difference(s))
}

// IsSubsetOf reports whether every element of s is in other.
func (s *Set[T]) IsSubsetOf(other *Set[T]) bool {
	if s.Len() > other.Len() {
		return false
	}
	for e := range s.m {
		if !other.Contains(e) {
			return false
		}
	}
	return true
}

// IsSupersetOf reports whether every element of other is in s.
func (s *Set[T]) IsSupersetOf(other *Set[T]) bool {
	return other.IsSubsetOf(s)
}
