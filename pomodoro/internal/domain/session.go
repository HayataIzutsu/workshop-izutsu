package domain

import (
	"errors"
	"time"
)

type SessionType string

const (
	WorkSession       SessionType = "work"
	ShortBreakSession SessionType = "shortBreak"
	LongBreakSession  SessionType = "longBreak"
)

type Session struct {
	ID          string      `json:"id"`
	Type        SessionType `json:"type"`
	DurationSec int         `json:"durationSec"`
	CompletedAt time.Time   `json:"completedAt"`
	Task        string      `json:"task,omitempty"`
	Note        string      `json:"note,omitempty"`
}

type Stats struct {
	CompletedWorkSessions int `json:"completedWorkSessions"`
	FocusTimeSec          int `json:"focusTimeSec"`
}

func (session Session) Validate() error {
	switch session.Type {
	case WorkSession, ShortBreakSession, LongBreakSession:
	default:
		return errors.New("invalid session type")
	}
	if session.DurationSec <= 0 || session.DurationSec > 24*60*60 {
		return errors.New("durationSec must be between 1 and 86400")
	}
	if session.CompletedAt.IsZero() {
		return errors.New("completedAt is required")
	}
	return nil
}

func StatsForDay(sessions []Session, day time.Time, location *time.Location) Stats {
	if location == nil {
		location = time.UTC
	}
	stats := Stats{}
	day = day.In(location)
	for _, session := range sessions {
		completedAt := session.CompletedAt.In(location)
		if completedAt.Year() != day.Year() || completedAt.YearDay() != day.YearDay() {
			continue
		}
		if session.Type == WorkSession {
			stats.CompletedWorkSessions++
			stats.FocusTimeSec += session.DurationSec
		}
	}
	return stats
}
