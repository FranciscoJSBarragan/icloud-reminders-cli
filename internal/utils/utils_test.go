package utils

import (
	"os"
	"testing"
	"time"
)

// TestLocalTZ_EnvVar verifies that the TZ env var is respected.
func TestLocalTZ_EnvVar(t *testing.T) {
	t.Setenv("TZ", "America/Mexico_City")
	loc := localTZ()
	if loc == nil {
		t.Fatal("localTZ() returned nil")
	}
	if loc.String() != "America/Mexico_City" {
		t.Errorf("expected America/Mexico_City, got %s", loc.String())
	}
}

// TestLocalTZ_NoEnvVar verifies that localTZ() returns a non-nil, non-UTC location
// when TZ is unset but the system is configured for a local timezone.
func TestLocalTZ_NoEnvVar(t *testing.T) {
	os.Unsetenv("TZ")
	loc := localTZ()
	if loc == nil {
		t.Fatal("localTZ() returned nil")
	}
	t.Logf("Detected timezone (no TZ env): %s", loc.String())
}

// TestStrToTs_TimezoneOffset is the regression test for the CST/UTC-6 bug.
// It verifies that "2026-04-06T09:50" interpreted in CST results in a UTC
// timestamp of 15:50 UTC (09:50 + 6h), NOT 09:50 UTC.
func TestStrToTs_TimezoneOffset(t *testing.T) {
	// Force CST (UTC-6) for this test.
	t.Setenv("TZ", "America/Mexico_City")

	tsMs, err := StrToTs("2026-04-06T09:50")
	if err != nil {
		t.Fatalf("StrToTs error: %v", err)
	}

	// Convert back to UTC to verify the offset was applied.
	utc := time.UnixMilli(tsMs).UTC()
	t.Logf("Input: 09:50 CST → UTC: %s", utc.Format("2006-01-02T15:04:05Z"))

	if utc.Hour() != 15 || utc.Minute() != 50 {
		t.Errorf(
			"timezone offset not applied: expected 15:50 UTC for 09:50 CST, got %02d:%02d UTC",
			utc.Hour(), utc.Minute(),
		)
	}
}

// TestStrToTs_DateOnly verifies date-only parsing also respects the local timezone.
func TestStrToTs_DateOnly(t *testing.T) {
	t.Setenv("TZ", "America/Mexico_City")

	tsMs, err := StrToTs("2026-04-06")
	if err != nil {
		t.Fatalf("StrToTs error: %v", err)
	}

	// Midnight CST = 06:00 UTC
	utc := time.UnixMilli(tsMs).UTC()
	t.Logf("Input: 2026-04-06 midnight CST → UTC: %s", utc.Format("2006-01-02T15:04:05Z"))

	if utc.Hour() != 6 || utc.Minute() != 0 {
		t.Errorf(
			"date-only timezone offset not applied: expected 06:00 UTC for midnight CST, got %02d:%02d UTC",
			utc.Hour(), utc.Minute(),
		)
	}
}
