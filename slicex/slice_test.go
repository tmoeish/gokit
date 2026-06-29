package slicex_test

import (
	"testing"

	"github.com/tmoeish/gokit/slicex"
)

func TestMap(t *testing.T) {
	got := slicex.Map([]int{1, 2, 3}, func(v int) int { return v * 2 })
	want := []int{2, 4, 6}
	if !slicex.Equal(got, want) {
		t.Fatalf("Map: got %v, want %v", got, want)
	}
}

func TestFilter(t *testing.T) {
	got := slicex.Filter([]int{1, 2, 3, 4, 5}, func(v int) bool { return v%2 == 0 })
	want := []int{2, 4}
	if !slicex.Equal(got, want) {
		t.Fatalf("Filter: got %v, want %v", got, want)
	}
}

func TestReduce(t *testing.T) {
	got := slicex.Reduce([]int{1, 2, 3, 4, 5}, 0, func(acc, v int) int { return acc + v })
	if got != 15 {
		t.Fatalf("Reduce: got %d, want 15", got)
	}
}

func TestContains(t *testing.T) {
	s := []int{1, 2, 3}
	if !slicex.Contains(s, 2) {
		t.Fatal("expected Contains(2) to be true")
	}
	if slicex.Contains(s, 5) {
		t.Fatal("expected Contains(5) to be false")
	}
}

func TestUnique(t *testing.T) {
	got := slicex.Unique([]int{1, 2, 2, 3, 1, 4})
	want := []int{1, 2, 3, 4}
	if !slicex.Equal(got, want) {
		t.Fatalf("Unique: got %v, want %v", got, want)
	}
}

func TestFlatten(t *testing.T) {
	got := slicex.Flatten([][]int{{1, 2}, {3, 4}, {5}})
	want := []int{1, 2, 3, 4, 5}
	if !slicex.Equal(got, want) {
		t.Fatalf("Flatten: got %v, want %v", got, want)
	}
}

func TestChunk(t *testing.T) {
	got := slicex.Chunk([]int{1, 2, 3, 4, 5}, 2)
	if len(got) != 3 || !slicex.Equal(got[0], []int{1, 2}) || !slicex.Equal(got[2], []int{5}) {
		t.Fatalf("Chunk: unexpected result %v", got)
	}
}

func TestGroupBy(t *testing.T) {
	got := slicex.GroupBy([]int{1, 2, 3, 4, 5, 6}, func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	})
	if len(got["even"]) != 3 || len(got["odd"]) != 3 {
		t.Fatalf("GroupBy: unexpected result %v", got)
	}
}

func TestPartition(t *testing.T) {
	evens, odds := slicex.Partition([]int{1, 2, 3, 4, 5}, func(v int) bool { return v%2 == 0 })
	if !slicex.Equal(evens, []int{2, 4}) || !slicex.Equal(odds, []int{1, 3, 5}) {
		t.Fatalf("Partition: evens=%v odds=%v", evens, odds)
	}
}

func TestIntersection(t *testing.T) {
	got := slicex.Intersection([]int{1, 2, 3, 4}, []int{2, 4, 6})
	want := []int{2, 4}
	if !slicex.Equal(got, want) {
		t.Fatalf("Intersection: got %v, want %v", got, want)
	}
}

func TestDifference(t *testing.T) {
	got := slicex.Difference([]int{1, 2, 3, 4}, []int{2, 4})
	want := []int{1, 3}
	if !slicex.Equal(got, want) {
		t.Fatalf("Difference: got %v, want %v", got, want)
	}
}

func TestUnion(t *testing.T) {
	got := slicex.Union([]int{1, 2, 3}, []int{2, 3, 4}, []int{4, 5})
	want := []int{1, 2, 3, 4, 5}
	if !slicex.Equal(got, want) {
		t.Fatalf("Union: got %v, want %v", got, want)
	}
}

func TestSum(t *testing.T) {
	if slicex.Sum([]int{1, 2, 3, 4, 5}) != 15 {
		t.Fatal("Sum failed")
	}
}

func TestMinMax(t *testing.T) {
	mn, ok := slicex.Min([]int{3, 1, 4, 1, 5, 9, 2, 6})
	if !ok || mn != 1 {
		t.Fatalf("Min: got %d, ok=%v", mn, ok)
	}
	mx, ok := slicex.Max([]int{3, 1, 4, 1, 5, 9, 2, 6})
	if !ok || mx != 9 {
		t.Fatalf("Max: got %d, ok=%v", mx, ok)
	}
}

func TestReverse(t *testing.T) {
	got := slicex.Reverse([]int{1, 2, 3})
	want := []int{3, 2, 1}
	if !slicex.Equal(got, want) {
		t.Fatalf("Reverse: got %v, want %v", got, want)
	}
}

func TestTakeDrop(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	if !slicex.Equal(slicex.Take(s, 3), []int{1, 2, 3}) {
		t.Fatal("Take failed")
	}
	if !slicex.Equal(slicex.Drop(s, 2), []int{3, 4, 5}) {
		t.Fatal("Drop failed")
	}
}

func TestWindow(t *testing.T) {
	got := slicex.Window([]int{1, 2, 3, 4}, 2)
	if len(got) != 3 {
		t.Fatalf("Window: got %v", got)
	}
}

func TestShuffle(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	got := slicex.Shuffle(s)
	// original should be unchanged
	if !slicex.Equal(s, []int{1, 2, 3, 4, 5}) {
		t.Fatal("Shuffle modified original")
	}
	if len(got) != 5 {
		t.Fatalf("Shuffle: len=%d", len(got))
	}
}

func TestEveryNoneSome(t *testing.T) {
	s := []int{2, 4, 6, 8}
	if !slicex.Every(s, func(v int) bool { return v%2 == 0 }) {
		t.Fatal("Every failed")
	}
	if !slicex.None(s, func(v int) bool { return v%2 != 0 }) {
		t.Fatal("None failed")
	}
	if !slicex.Some(s, func(v int) bool { return v > 5 }) {
		t.Fatal("Some failed")
	}
}
