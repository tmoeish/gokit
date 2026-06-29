// Package mapx provides generic map utilities.
package mapx

import "cmp"

// Keys returns all keys of m in unspecified order.
func Keys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// Values returns all values of m in unspecified order.
func Values[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

// Pair represents a key-value pair.
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

// Entries returns all key-value pairs of m in unspecified order.
func Entries[K comparable, V any](m map[K]V) []Pair[K, V] {
	result := make([]Pair[K, V], 0, len(m))
	for k, v := range m {
		result = append(result, Pair[K, V]{Key: k, Value: v})
	}
	return result
}

// FromEntries creates a map from a slice of key-value pairs.
func FromEntries[K comparable, V any](entries []Pair[K, V]) map[K]V {
	result := make(map[K]V, len(entries))
	for _, e := range entries {
		result[e.Key] = e.Value
	}
	return result
}

// Merge merges multiple maps into one. Later maps overwrite earlier ones on conflict.
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// MergeWith merges multiple maps into one, using fn to resolve conflicts.
func MergeWith[K comparable, V any](fn func(existing, incoming V) V, maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			if existing, ok := result[k]; ok {
				result[k] = fn(existing, v)
			} else {
				result[k] = v
			}
		}
	}
	return result
}

// Filter returns a new map containing only entries for which fn returns true.
func Filter[K comparable, V any](m map[K]V, fn func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if fn(k, v) {
			result[k] = v
		}
	}
	return result
}

// Map transforms each key-value pair using fn and returns a new map.
func Map[K1 comparable, V1 any, K2 comparable, V2 any](
	m map[K1]V1,
	fn func(K1, V1) (K2, V2),
) map[K2]V2 {
	result := make(map[K2]V2, len(m))
	for k, v := range m {
		k2, v2 := fn(k, v)
		result[k2] = v2
	}
	return result
}

// MapValues transforms each value using fn, keeping keys the same.
func MapValues[K comparable, V1, V2 any](m map[K]V1, fn func(K, V1) V2) map[K]V2 {
	result := make(map[K]V2, len(m))
	for k, v := range m {
		result[k] = fn(k, v)
	}
	return result
}

// MapKeys transforms each key using fn, keeping values the same.
// If fn produces duplicate keys, later entries overwrite earlier ones.
func MapKeys[K1, K2 comparable, V any](m map[K1]V, fn func(K1, V) K2) map[K2]V {
	result := make(map[K2]V, len(m))
	for k, v := range m {
		result[fn(k, v)] = v
	}
	return result
}

// Invert swaps keys and values. If multiple keys map to the same value,
// one is kept arbitrarily.
func Invert[K, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K, len(m))
	for k, v := range m {
		result[v] = k
	}
	return result
}

// Pick returns a new map containing only the specified keys.
func Pick[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	result := make(map[K]V, len(keys))
	for _, k := range keys {
		if v, ok := m[k]; ok {
			result[k] = v
		}
	}
	return result
}

// Omit returns a new map excluding the specified keys.
func Omit[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	exclude := make(map[K]struct{}, len(keys))
	for _, k := range keys {
		exclude[k] = struct{}{}
	}
	result := make(map[K]V)
	for k, v := range m {
		if _, ok := exclude[k]; !ok {
			result[k] = v
		}
	}
	return result
}

// Has reports whether m contains key.
func Has[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

// GetOrDefault returns the value for key, or defaultVal if the key is absent.
func GetOrDefault[K comparable, V any](m map[K]V, key K, defaultVal V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultVal
}

// GetOrSet returns the value for key. If absent, sets it to defaultVal and returns it.
func GetOrSet[K comparable, V any](m map[K]V, key K, defaultVal V) V {
	if v, ok := m[key]; ok {
		return v
	}
	m[key] = defaultVal
	return defaultVal
}

// ForEach calls fn for each key-value pair.
func ForEach[K comparable, V any](m map[K]V, fn func(K, V)) {
	for k, v := range m {
		fn(k, v)
	}
}

// Count returns the number of entries for which fn returns true.
func Count[K comparable, V any](m map[K]V, fn func(K, V) bool) int {
	n := 0
	for k, v := range m {
		if fn(k, v) {
			n++
		}
	}
	return n
}

// Every reports whether all entries satisfy fn.
func Every[K comparable, V any](m map[K]V, fn func(K, V) bool) bool {
	for k, v := range m {
		if !fn(k, v) {
			return false
		}
	}
	return true
}

// Some reports whether at least one entry satisfies fn.
func Some[K comparable, V any](m map[K]V, fn func(K, V) bool) bool {
	for k, v := range m {
		if fn(k, v) {
			return true
		}
	}
	return false
}

// None reports whether no entries satisfy fn.
func None[K comparable, V any](m map[K]V, fn func(K, V) bool) bool {
	return !Some(m, fn)
}

// Clone returns a shallow copy of m.
func Clone[K comparable, V any](m map[K]V) map[K]V {
	result := make(map[K]V, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

// GroupBy groups the slice elements by a key extracted by fn.
func GroupBy[T any, K comparable](s []T, fn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range s {
		k := fn(v)
		result[k] = append(result[k], v)
	}
	return result
}

// SortedKeys returns the keys of m in sorted order.
func SortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	keys := Keys(m)
	sortOrdered(keys)
	return keys
}

// sortOrdered sorts a slice of ordered values in place.
func sortOrdered[T cmp.Ordered](s []T) {
	// insertion sort for small slices; falls back to a simple approach
	n := len(s)
	for i := 1; i < n; i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
