package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pomodoro/internal/clock"
	"pomodoro/internal/domain"
	"pomodoro/internal/storage"
)

func TestStatsHistoryHandlerReturnsRequestedDays(t *testing.T) {
	location := time.UTC
	now := time.Date(2026, 9, 8, 10, 0, 0, 0, location)
	store := storage.NewMemorySessionStore()
	_ = store.Add(domain.Session{Type: domain.WorkSession, DurationSec: 1500, CompletedAt: now})
	recorder := httptest.NewRecorder()
	handler := StatsHistoryHandler{Store: store, Clock: clock.FixedClock{Current: now}, Location: location}

	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/stats/history?days=3", nil))

	var history []struct {
		Date string `json:"date"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&history); err != nil {
		t.Fatal(err)
	}
	if len(history) != 3 {
		t.Fatalf("history length = %d, want 3", len(history))
	}
}

func TestDeleteSessionHandlerDeletesByID(t *testing.T) {
	store := storage.NewMemorySessionStore()
	_ = store.Add(domain.Session{ID: "delete-me", Type: domain.WorkSession, DurationSec: 1, CompletedAt: time.Now()})
	recorder := httptest.NewRecorder()

	(DeleteSessionHandler{Store: store}).ServeHTTP(recorder, httptest.NewRequest(http.MethodDelete, "/api/sessions/delete?id=delete-me", nil))

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
	sessions, _ := store.List()
	if len(sessions) != 0 {
		t.Fatalf("sessions = %+v, want empty", sessions)
	}
}
