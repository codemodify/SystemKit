package Helpers

import (
	"math"
	"time"
)

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

// DeltaDays -
func DeltaDays(t1, t2 time.Time) int {
	return int(math.Abs(t1.Sub(t2).Hours() / 24))
}

// MillisecondsToDateTime -
func MillisecondsToDateTime(milliseconds int64) string {
	var result string
	if milliseconds > 0 {
		unixTimeUTC := time.Unix(0, milliseconds*int64(time.Millisecond))
		result = unixTimeUTC.Format("2006-01-02 15:04:05")
		//unitTimeInRFC3339 :=unixTimeUTC.Format(time.RFC3339)
	}
	return result
}

// MillisecondsToDate -
func MillisecondsToDate(milliseconds int64) string {
	var result string
	if milliseconds > 0 {
		unixTimeUTC := time.Unix(0, milliseconds*int64(time.Millisecond))
		result = unixTimeUTC.Format("2006-01-02")
	}
	return result
}
