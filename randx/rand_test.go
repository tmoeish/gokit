package randx_test

import (
	"strings"
	"testing"

	"github.com/tmoeish/gokit/randx"
)

func TestString(t *testing.T) {
	s := randx.String(16)
	if len(s) != 16 {
		t.Fatalf("String: len=%d", len(s))
	}
}

func TestStringWithCharset(t *testing.T) {
	s := randx.StringWithCharset(10, "abc")
	if len(s) != 10 {
		t.Fatalf("StringWithCharset: len=%d", len(s))
	}
	for _, r := range s {
		if !strings.ContainsRune("abc", r) {
			t.Fatalf("StringWithCharset: unexpected char %c", r)
		}
	}
}

func TestUUID(t *testing.T) {
	u := randx.UUID()
	if len(u) != 36 {
		t.Fatalf("UUID: len=%d, got=%s", len(u), u)
	}
	u2 := randx.UUID()
	if u == u2 {
		t.Fatal("UUIDs should be different")
	}
}

func TestInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		n := randx.Int(0, 10)
		if n < 0 || n >= 10 {
			t.Fatalf("Int: out of range %d", n)
		}
	}
}

func TestFloat64(t *testing.T) {
	for i := 0; i < 100; i++ {
		f := randx.Float64(1.0, 2.0)
		if f < 1.0 || f >= 2.0 {
			t.Fatalf("Float64: out of range %f", f)
		}
	}
}

func TestBool(t *testing.T) {
	seenTrue, seenFalse := false, false
	for i := 0; i < 100; i++ {
		if randx.Bool() {
			seenTrue = true
		} else {
			seenFalse = true
		}
	}
	if !seenTrue || !seenFalse {
		t.Fatal("Bool should produce both true and false")
	}
}

func TestShuffle(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	result := randx.Shuffle(s)
	if len(result) != 5 {
		t.Fatal("Shuffle length mismatch")
	}
	// original should be unchanged
	if s[0] != 1 || s[4] != 5 {
		t.Fatal("Shuffle modified original")
	}
}

func TestSample(t *testing.T) {
	s := []int{1, 2, 3}
	v, ok := randx.Sample(s)
	if !ok || v < 1 || v > 3 {
		t.Fatalf("Sample: v=%d ok=%v", v, ok)
	}
	_, ok = randx.Sample([]int{})
	if ok {
		t.Fatal("Sample empty should return false")
	}
}

func TestWeightedChoice(t *testing.T) {
	// weight 0 for index 0, high weight for index 1
	counts := make([]int, 3)
	for i := 0; i < 1000; i++ {
		idx := randx.WeightedChoice([]float64{0, 10, 0})
		counts[idx]++
	}
	if counts[1] < 990 {
		t.Fatalf("WeightedChoice: expected index 1 to be chosen almost always, got %v", counts)
	}
}
