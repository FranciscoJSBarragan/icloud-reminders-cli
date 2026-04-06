package utils

import (
	"testing"
	"time"
)

// TestStrToTs_NaiveUTC verifies that StrToTs treats input as naive wall-clock
// time (stored as UTC epoch ms) so CloudKit Reminders displays it correctly.
// CloudKit Reminders shows the stored ms value directly without UTC→local conversion.
func TestStrToTs_NaiveUTC(t *testing.T) {
	tsMs, err := StrToTs("2026-04-06T13:30")
	if err != nil {
		t.Fatalf("StrToTs error: %v", err)
	}

	// The stored ms should represent 13:30 UTC (naive), NOT 13:30 CST → 19:30 UTC.
	utc := time.UnixMilli(tsMs).UTC()
	t.Logf("Input: 13:30 → stored as UTC: %s", utc.Format("2006-01-02T15:04:05Z"))

	if utc.Hour() != 13 || utc.Minute() != 30 {
		t.Errorf(
			"expected 13:30 UTC (naive), got %02d:%02d UTC — timezone conversion should NOT be applied",
			utc.Hour(), utc.Minute(),
		)
	}
}

// TestStrToTs_DateOnly verifies date-only parsing stores midnight UTC.
func TestStrToTs_DateOnly(t *testing.T) {
	tsMs, err := StrToTs("2026-04-06")
	if err != nil {
		t.Fatalf("StrToTs error: %v", err)
	}

	utc := time.UnixMilli(tsMs).UTC()
	t.Logf("Input: 2026-04-06 → stored as UTC: %s", utc.Format("2006-01-02T15:04:05Z"))

	if utc.Hour() != 0 || utc.Minute() != 0 {
		t.Errorf("date-only should be midnight UTC, got %02d:%02d", utc.Hour(), utc.Minute())
	}
}

// TestTsToStr_Roundtrip verifies that StrToTs → TsToStr preserves the date.
func TestTsToStr_Roundtrip(t *testing.T) {
	tsMs, err := StrToTs("2026-04-06")
	if err != nil {
		t.Fatalf("StrToTs error: %v", err)
	}

	result := TsToStr(tsMs)
	if result != "2026-04-06" {
		t.Errorf("roundtrip failed: expected 2026-04-06, got %s", result)
	}
}

// TestStrToTs_TimezoneIndependent verifies the result is the same
// regardless of the TZ environment variable (naive behavior).
func TestStrToTs_TimezoneIndependent(t *testing.T) {
	// Parse without any TZ influence
	ts1, err := StrToTs("2026-04-06T09:50")
	if err != nil {
		t.Fatal(err)
	}

	// The result should always be 09:50 UTC regardless of system timezone
	utc := time.UnixMilli(ts1).UTC()
	if utc.Hour() != 9 || utc.Minute() != 50 {
		t.Errorf("expected 09:50 UTC, got %02d:%02d", utc.Hour(), utc.Minute())
	}
}

// TestHasTime verifies time component detection.
func TestHasTime(t *testing.T) {
	if !HasTime("2026-04-06T13:30") {
		t.Error("expected HasTime=true for YYYY-MM-DDTHH:MM")
	}
	if HasTime("2026-04-06") {
		t.Error("expected HasTime=false for YYYY-MM-DD")
	}
}
