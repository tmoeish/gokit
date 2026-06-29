package conv_test

import (
	"testing"

	"github.com/tmoeish/gokit/conv"
)

func TestToString(t *testing.T) {
	tests := []struct {
		in   any
		want string
	}{
		{"hello", "hello"},
		{42, "42"},
		{3.14, "3.14"},
		{true, "true"},
		{false, "false"},
		{nil, ""},
		{[]byte("bytes"), "bytes"},
	}
	for _, tc := range tests {
		if got := conv.ToString(tc.in); got != tc.want {
			t.Fatalf("ToString(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		in   any
		want int64
	}{
		{42, 42},
		{"123", 123},
		{3.7, 3},
		{true, 1},
		{false, 0},
		{nil, 0},
	}
	for _, tc := range tests {
		got, err := conv.ToInt64(tc.in)
		if err != nil || got != tc.want {
			t.Fatalf("ToInt64(%v) = %d, %v; want %d", tc.in, got, err, tc.want)
		}
	}
}

func TestToFloat64(t *testing.T) {
	got, err := conv.ToFloat64("3.14")
	if err != nil || got != 3.14 {
		t.Fatalf("ToFloat64: %f %v", got, err)
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		in   any
		want bool
	}{
		{true, true},
		{false, false},
		{"true", true},
		{"false", false},
		{"1", true},
		{"0", false},
		{1, true},
		{0, false},
	}
	for _, tc := range tests {
		got, err := conv.ToBool(tc.in)
		if err != nil || got != tc.want {
			t.Fatalf("ToBool(%v) = %v, %v; want %v", tc.in, got, err, tc.want)
		}
	}
}

func TestToStringSlice(t *testing.T) {
	got := conv.ToStringSlice([]any{1, "two", true})
	if len(got) != 3 || got[0] != "1" || got[1] != "two" || got[2] != "true" {
		t.Fatalf("ToStringSlice: %v", got)
	}
}
