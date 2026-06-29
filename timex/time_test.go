package timex_test

import (
	"testing"
	"time"

	"github.com/tmoeish/gokit/timex"
)

func TestBeginEndOfDay(t *testing.T) {
	now := time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC)
	bod := timex.BeginOfDay(now)
	eod := timex.EndOfDay(now)
	if bod.Hour() != 0 || bod.Minute() != 0 || bod.Second() != 0 {
		t.Fatalf("BeginOfDay: %v", bod)
	}
	if eod.Hour() != 23 || eod.Minute() != 59 || eod.Second() != 59 {
		t.Fatalf("EndOfDay: %v", eod)
	}
}

func TestBeginEndOfMonth(t *testing.T) {
	now := time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC)
	bom := timex.BeginOfMonth(now)
	eom := timex.EndOfMonth(now)
	if bom.Day() != 1 {
		t.Fatalf("BeginOfMonth: %v", bom)
	}
	if eom.Day() != 31 {
		t.Fatalf("EndOfMonth March: %v", eom)
	}
}

func TestIsLeapYear(t *testing.T) {
	leaps := []int{2000, 2004, 2024}
	for _, y := range leaps {
		if !timex.IsLeapYear(y) {
			t.Fatalf("IsLeapYear(%d) should be true", y)
		}
	}
	notLeaps := []int{1900, 2001, 2023}
	for _, y := range notLeaps {
		if timex.IsLeapYear(y) {
			t.Fatalf("IsLeapYear(%d) should be false", y)
		}
	}
}

func TestIsWeekend(t *testing.T) {
	sat := time.Date(2024, 3, 16, 0, 0, 0, 0, time.UTC) // Saturday
	sun := time.Date(2024, 3, 17, 0, 0, 0, 0, time.UTC) // Sunday
	mon := time.Date(2024, 3, 18, 0, 0, 0, 0, time.UTC) // Monday
	if !timex.IsWeekend(sat) || !timex.IsWeekend(sun) {
		t.Fatal("IsWeekend")
	}
	if timex.IsWeekend(mon) {
		t.Fatal("Monday should not be weekend")
	}
}

func TestDaysBetween(t *testing.T) {
	a := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	b := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
	if timex.DaysBetween(a, b) != 10 {
		t.Fatal("DaysBetween")
	}
}

func TestAgeAt(t *testing.T) {
	birth := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	at := time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC)
	if timex.AgeAt(birth, at) != 34 {
		t.Fatalf("AgeAt: got %d", timex.AgeAt(birth, at))
	}
}

func TestAddWorkDays(t *testing.T) {
	// Monday 2024-03-18 + 5 work days = Monday 2024-03-25
	mon := time.Date(2024, 3, 18, 0, 0, 0, 0, time.UTC)
	result := timex.AddWorkDays(mon, 5)
	if result.Weekday() != time.Monday || result.Day() != 25 {
		t.Fatalf("AddWorkDays: got %v", result)
	}
}

func TestFormatDuration(t *testing.T) {
	d := 2*time.Hour + 3*time.Minute + 4*time.Second
	got := timex.FormatDuration(d)
	if got != "2h 3m 4s" {
		t.Fatalf("FormatDuration: got %q", got)
	}
}

func TestParseDate(t *testing.T) {
	d, err := timex.ParseDate("2024-03-15")
	if err != nil || d.Year() != 2024 || d.Month() != 3 || d.Day() != 15 {
		t.Fatalf("ParseDate: %v %v", d, err)
	}
}

func TestFormatDate(t *testing.T) {
	d := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	if timex.FormatDate(d) != "2024-03-15" {
		t.Fatal("FormatDate")
	}
}

func TestUnixMillis(t *testing.T) {
	d := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	ms := timex.UnixMillis(d)
	back := timex.FromUnixMillis(ms)
	if !d.Equal(back) {
		t.Fatalf("UnixMillis roundtrip: %v -> %d -> %v", d, ms, back)
	}
}

func TestIsSameDay(t *testing.T) {
	a := time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
	b := time.Date(2024, 3, 15, 22, 0, 0, 0, time.UTC)
	c := time.Date(2024, 3, 16, 0, 0, 0, 0, time.UTC)
	if !timex.IsSameDay(a, b) {
		t.Fatal("IsSameDay same day")
	}
	if timex.IsSameDay(a, c) {
		t.Fatal("IsSameDay different day")
	}
}
