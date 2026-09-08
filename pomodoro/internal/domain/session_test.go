package domain

import (
	"testing"
	"time"
)

func TestStatsForDayCountsOnlyWorkSessionsOnTheRequestedDay(t *testing.T) {
	location := time.FixedZone("JST", 9*60*60)
	day := time.Date(2026, 9, 8, 12, 0, 0, 0, location)
	sessions := []Session{
		{Type: WorkSession, DurationSec: 1500, CompletedAt: day.Add(-time.Hour)},
		{Type: ShortBreakSession, DurationSec: 300, CompletedAt: day.Add(-time.Minute)},
		{Type: WorkSession, DurationSec: 1500, CompletedAt: day.AddDate(0, 0, -1)},
	}

	stats := StatsForDay(sessions, day, location)

	if stats.CompletedWorkSessions != 1 || stats.FocusTimeSec != 1500 {
		t.Fatalf("stats = %+v, want one work session and 1500 seconds", stats)
	}
}

func TestSessionValidateRejectsInvalidValues(t *testing.T) {
	cases := []Session{
		{Type: "invalid", DurationSec: 1, CompletedAt: time.Now()},
		{Type: WorkSession, DurationSec: 0, CompletedAt: time.Now()},
		{Type: WorkSession, DurationSec: 1},
	}
	for _, session := range cases {
		if err := session.Validate(); err == nil {
			t.Errorf("Validate(%+v) returned nil", session)
		}
	}
}
