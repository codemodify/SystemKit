package Helpers

import "time"

// ToNanoSeconds -
func ToNanoSeconds(t time.Time) int64 {
	return t.UTC().UnixNano() / int64(time.Nanosecond)
}

// ToMicroSeconds -
func ToMicroSeconds(t time.Time) int64 {
	return t.UTC().UnixNano() / int64(time.Microsecond)
}

// ToMilliSeconds -
func ToMilliSeconds(t time.Time) int64 {
	return t.UTC().UnixNano() / int64(time.Millisecond)
}

// StartOfMonth -
func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth -
func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, -1)
}
