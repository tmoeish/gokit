package setx_test

import (
	"sort"
	"testing"

	"github.com/tmoeish/gokit/setx"
)

func TestAddContainsRemove(t *testing.T) {
	s := setx.New(1, 2, 3)
	if !s.Contains(2) {
		t.Fatal("expected set to contain 2")
	}
	s.Add(2, 4)
	if s.Len() != 4 {
		t.Fatalf("expected len 4, got %d", s.Len())
	}
	s.Remove(2)
	if s.Contains(2) || s.Len() != 3 {
		t.Fatalf("Remove failed: %v", s.ToSlice())
	}
}

func TestContainsAllAny(t *testing.T) {
	s := setx.New("a", "b", "c")
	if !s.ContainsAll("a", "c") {
		t.Fatal("ContainsAll should be true")
	}
	if s.ContainsAll("a", "z") {
		t.Fatal("ContainsAll should be false")
	}
	if !s.ContainsAny("z", "b") {
		t.Fatal("ContainsAny should be true")
	}
	if s.ContainsAny("x", "y") {
		t.Fatal("ContainsAny should be false")
	}
}

func TestUnionIntersectionDifference(t *testing.T) {
	a := setx.New(1, 2, 3)
	b := setx.New(2, 3, 4)

	union := sortedInts(a.Union(b).ToSlice())
	if !equalInts(union, []int{1, 2, 3, 4}) {
		t.Fatalf("Union: %v", union)
	}

	inter := sortedInts(a.Intersection(b).ToSlice())
	if !equalInts(inter, []int{2, 3}) {
		t.Fatalf("Intersection: %v", inter)
	}

	diff := sortedInts(a.Difference(b).ToSlice())
	if !equalInts(diff, []int{1}) {
		t.Fatalf("Difference: %v", diff)
	}

	sym := sortedInts(a.SymmetricDifference(b).ToSlice())
	if !equalInts(sym, []int{1, 4}) {
		t.Fatalf("SymmetricDifference: %v", sym)
	}
}

func TestSubsetSupersetEqual(t *testing.T) {
	a := setx.New(1, 2)
	b := setx.New(1, 2, 3)
	if !a.IsSubsetOf(b) {
		t.Fatal("a should be subset of b")
	}
	if !b.IsSupersetOf(a) {
		t.Fatal("b should be superset of a")
	}
	if a.Equal(b) {
		t.Fatal("a should not equal b")
	}
	if !a.Equal(setx.New(2, 1)) {
		t.Fatal("a should equal {2,1}")
	}
}

func TestCloneClearEmpty(t *testing.T) {
	s := setx.New(1, 2, 3)
	c := s.Clone()
	c.Add(4)
	if s.Contains(4) {
		t.Fatal("Clone should be independent")
	}
	s.Clear()
	if !s.IsEmpty() {
		t.Fatal("Clear should empty the set")
	}
}

func sortedInts(s []int) []int {
	sort.Ints(s)
	return s
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
