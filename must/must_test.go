package must_test

import (
	"errors"
	"testing"

	"github.com/tmoeish/gokit/must"
)

func TestMust(t *testing.T) {
	v := must.Must(42, nil)
	if v != 42 {
		t.Fatal("Must: wrong value")
	}
}

func TestMustPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Must should panic on error")
		}
	}()
	must.Must(0, errors.New("boom"))
}

func TestOK(t *testing.T) {
	must.OK(nil) // should not panic
}

func TestOKPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("OK should panic on error")
		}
	}()
	must.OK(errors.New("boom"))
}

func TestAssert(t *testing.T) {
	must.Assert(true, "should not panic")
}

func TestAssertPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Assert should panic")
		}
	}()
	must.Assert(false, "expected failure")
}
