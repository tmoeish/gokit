package ptrx_test

import (
	"testing"

	"github.com/tmoeish/gokit/ptrx"
)

func TestOf(t *testing.T) {
	p := ptrx.Of(42)
	if p == nil || *p != 42 {
		t.Fatal("Of failed")
	}
	s := ptrx.Of("hello")
	if s == nil || *s != "hello" {
		t.Fatal("Of string failed")
	}
}

func TestDeref(t *testing.T) {
	v := 99
	if ptrx.Deref(&v, 0) != 99 {
		t.Fatal("Deref non-nil failed")
	}
	if ptrx.Deref[int](nil, 42) != 42 {
		t.Fatal("Deref nil failed")
	}
}

func TestDerefZero(t *testing.T) {
	if ptrx.DerefZero[int](nil) != 0 {
		t.Fatal("DerefZero failed")
	}
}

func TestIsNil(t *testing.T) {
	if !ptrx.IsNil[int](nil) {
		t.Fatal("IsNil(nil) should be true")
	}
	v := 1
	if ptrx.IsNil(&v) {
		t.Fatal("IsNil non-nil should be false")
	}
}

func TestCoalescePtr(t *testing.T) {
	a := ptrx.Of(1)
	b := ptrx.Of(2)
	got := ptrx.CoalescePtr[int](nil, a, b)
	if got != a {
		t.Fatal("CoalescePtr should return first non-nil")
	}
}
