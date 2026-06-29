// Package timex provides time utilities.
package timex

import (
	"fmt"
	"strings"
	"time"
)

// BeginOfDay returns the start of the day (00:00:00.000000000) for t.
func BeginOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of the day (23:59:59.999999999) for t.
func EndOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, 999999999, t.Location())
}

// BeginOfWeek returns the start of the week (Monday) containing t.
func BeginOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday becomes 7
	}
	daysBack := weekday - 1
	return BeginOfDay(t.AddDate(0, 0, -daysBack))
}

// EndOfWeek returns the end of the week (Sunday) containing t.
func EndOfWeek(t time.Time) time.Time {
	return EndOfDay(BeginOfWeek(t).AddDate(0, 0, 6))
}

// BeginOfMonth returns the first moment of the month of t.
func BeginOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the last moment of the month of t.
func EndOfMonth(t time.Time) time.Time {
	return EndOfDay(BeginOfMonth(t).AddDate(0, 1, -1))
}

// BeginOfYear returns the first moment of the year of t.
func BeginOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear returns the last moment of the year of t.
func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.December, 31, 23, 59, 59, 999999999, t.Location())
}

// BeginOfHour returns the start of the hour of t.
func BeginOfHour(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, t.Hour(), 0, 0, 0, t.Location())
}

// EndOfHour returns the end of the hour of t.
func EndOfHour(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, t.Hour(), 59, 59, 999999999, t.Location())
}

// IsLeapYear reports whether year is a leap year.
func IsLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

// IsWeekend reports whether t falls on Saturday or Sunday.
func IsWeekend(t time.Time) bool {
	w := t.Weekday()
	return w == time.Saturday || w == time.Sunday
}

// IsWeekday reports whether t falls on a weekday (Mon–Fri).
func IsWeekday(t time.Time) bool {
	return !IsWeekend(t)
}

// DaysBetween returns the number of complete days between a and b (absolute value).
func DaysBetween(a, b time.Time) int {
	d := b.Sub(a)
	if d < 0 {
		d = -d
	}
	return int(d.Hours() / 24)
}

// HoursBetween returns the number of hours between a and b (absolute value).
func HoursBetween(a, b time.Time) float64 {
	d := b.Sub(a)
	if d < 0 {
		d = -d
	}
	return d.Hours()
}

// Age returns the age in years relative to now.
func Age(birthDate time.Time) int {
	return AgeAt(birthDate, time.Now())
}

// AgeAt returns the age in years at the given reference time.
func AgeAt(birthDate, at time.Time) int {
	years := at.Year() - birthDate.Year()
	if at.YearDay() < birthDate.YearDay() {
		years--
	}
	if years < 0 {
		return 0
	}
	return years
}

// Today returns the current date with time zeroed.
func Today() time.Time {
	return BeginOfDay(time.Now())
}

// Yesterday returns the previous day with time zeroed.
func Yesterday() time.Time {
	return BeginOfDay(time.Now().AddDate(0, 0, -1))
}

// Tomorrow returns the next day with time zeroed.
func Tomorrow() time.Time {
	return BeginOfDay(time.Now().AddDate(0, 0, 1))
}

// AddWorkDays adds n work days (Mon–Fri) to t, skipping weekends.
// Negative n moves backward.
func AddWorkDays(t time.Time, n int) time.Time {
	sign := 1
	if n < 0 {
		sign = -1
		n = -n
	}
	for n > 0 {
		t = t.AddDate(0, 0, sign)
		if IsWeekday(t) {
			n--
		}
	}
	return t
}

// FormatDuration formats d as a human-readable string.
// e.g., "2h 3m 4s", "45s", "1d 2h".
func FormatDuration(d time.Duration) string {
	if d < 0 {
		return "-" + FormatDuration(-d)
	}
	if d == 0 {
		return "0s"
	}

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	parts := make([]string, 0, 4)
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	return strings.Join(parts, " ")
}

// DaysInMonth returns the number of days in the given month and year.
func DaysInMonth(year int, month time.Month) int {
	return BeginOfMonth(time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)).AddDate(0, 0, -1).Day()
}

// IsSameDay reports whether a and b fall on the same calendar day.
func IsSameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

// IsSameMonth reports whether a and b fall in the same year-month.
func IsSameMonth(a, b time.Time) bool {
	ay, am, _ := a.Date()
	by, bm, _ := b.Date()
	return ay == by && am == bm
}

// ParseDate parses a date string in "2006-01-02" layout.
func ParseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// FormatDate formats t as "2006-01-02".
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDatetime formats t as "2006-01-02 15:04:05".
func FormatDatetime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// UnixMillis returns the Unix time in milliseconds.
func UnixMillis(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

// FromUnixMillis converts Unix milliseconds to time.Time.
func FromUnixMillis(ms int64) time.Time {
	return time.Unix(ms/1000, (ms%1000)*1e6)
}
