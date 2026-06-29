package mapx_test

import (
	"testing"

	"github.com/tmoeish/gokit/mapx"
)

func TestKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := mapx.Keys(m)
	if len(keys) != 3 {
		t.Fatalf("Keys: expected 3, got %d", len(keys))
	}
}

func TestValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	vals := mapx.Values(m)
	if len(vals) != 2 {
		t.Fatalf("Values: expected 2, got %d", len(vals))
	}
}

func TestMerge(t *testing.T) {
	a := map[string]int{"a": 1, "b": 2}
	b := map[string]int{"b": 3, "c": 4}
	got := mapx.Merge(a, b)
	if got["a"] != 1 || got["b"] != 3 || got["c"] != 4 {
		t.Fatalf("Merge: got %v", got)
	}
}

func TestFilter(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	got := mapx.Filter(m, func(_ string, v int) bool { return v%2 == 0 })
	if len(got) != 2 {
		t.Fatalf("Filter: expected 2 entries, got %v", got)
	}
}

func TestInvert(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	got := mapx.Invert(m)
	if got[1] != "a" || got[2] != "b" {
		t.Fatalf("Invert: got %v", got)
	}
}

func TestPick(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	got := mapx.Pick(m, "a", "c")
	if len(got) != 2 || got["a"] != 1 || got["c"] != 3 {
		t.Fatalf("Pick: got %v", got)
	}
}

func TestOmit(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	got := mapx.Omit(m, "b")
	if len(got) != 2 || mapx.Has(got, "b") {
		t.Fatalf("Omit: got %v", got)
	}
}

func TestGetOrDefault(t *testing.T) {
	m := map[string]int{"a": 1}
	if mapx.GetOrDefault(m, "a", 0) != 1 {
		t.Fatal("GetOrDefault existing key failed")
	}
	if mapx.GetOrDefault(m, "z", 99) != 99 {
		t.Fatal("GetOrDefault missing key failed")
	}
}

func TestSortedKeys(t *testing.T) {
	m := map[int]string{3: "c", 1: "a", 2: "b"}
	keys := mapx.SortedKeys(m)
	if keys[0] != 1 || keys[1] != 2 || keys[2] != 3 {
		t.Fatalf("SortedKeys: got %v", keys)
	}
}
