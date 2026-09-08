package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"pomodoro/internal/clock"
	"pomodoro/internal/domain"
	"pomodoro/internal/storage"
)

func TestSessionHandlerStoresValidSession(t *testing.T) {
	store := storage.NewMemorySessionStore()
	handler := SessionHandler{Store: store, Clock: clock.FixedClock{Current: time.Now()}}
	request := httptest.NewRequest(http.MethodPost, "/api/sessions", strings.NewReader(`{"type":"work","durationSec":1500,"completedAt":"2026-09-08T10:00:00Z"}`))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}
	sessions, _ := store.List()
	if len(sessions) != 1 || sessions[0].Type != domain.WorkSession {
		t.Fatalf("sessions = %+v, want one work session", sessions)
	}
}

func TestSessionHandlerRejectsInvalidRequest(t *testing.T) {
	handler := SessionHandler{Store: storage.NewMemorySessionStore(), Clock: clock.FixedClock{Current: time.Now()}}
	request := httptest.NewRequest(http.MethodPost, "/api/sessions", strings.NewReader(`{"type":"unknown","durationSec":0}`))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	var response map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil || response["error"] == "" {
		t.Fatalf("response = %v, want JSON error", response)
	}
}

func TestStatsHandlerReturnsTodayStats(t *testing.T) {
	location := time.FixedZone("JST", 9*60*60)
	now := time.Date(2026, 9, 8, 10, 0, 0, 0, location)
	store := storage.NewMemorySessionStore()
	_ = store.Add(domain.Session{Type: domain.WorkSession, DurationSec: 1500, CompletedAt: now})
	handler := StatsHandler{Store: store, Clock: clock.FixedClock{Current: now}, Location: location}
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/stats/today", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var stats domain.Stats
	if err := json.NewDecoder(recorder.Body).Decode(&stats); err != nil {
		t.Fatal(err)
	}
	if stats.CompletedWorkSessions != 1 || stats.FocusTimeSec != 1500 {
		t.Fatalf("stats = %+v", stats)
	}
}
