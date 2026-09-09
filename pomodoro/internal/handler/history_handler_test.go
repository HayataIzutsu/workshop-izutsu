package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"pomodoro/internal/domain"
	"pomodoro/internal/storage"
)

func TestHistoryHandlerFiltersAndSortsSessions(t *testing.T) {
	store := storage.NewMemorySessionStore()
	_ = store.Add(domain.Session{ID: "old", Type: domain.WorkSession, DurationSec: 1500, CompletedAt: time.Date(2026, 9, 8, 9, 0, 0, 0, time.UTC), Task: "読書"})
	_ = store.Add(domain.Session{ID: "new", Type: domain.WorkSession, DurationSec: 1500, CompletedAt: time.Date(2026, 9, 8, 10, 0, 0, 0, time.UTC), Task: "実装"})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/history?task=実装", nil)

	(HistoryHandler{Store: store}).ServeHTTP(recorder, request)

	var sessions []domain.Session
	if err := json.NewDecoder(recorder.Body).Decode(&sessions); err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 || sessions[0].ID != "new" {
		t.Fatalf("sessions = %+v", sessions)
	}
}

func TestTaskHandlerReturnsUniqueTasks(t *testing.T) {
	store := storage.NewMemorySessionStore()
	_ = store.Add(domain.Session{Type: domain.WorkSession, Task: "実装", DurationSec: 1, CompletedAt: time.Now()})
	_ = store.Add(domain.Session{Type: domain.WorkSession, Task: "実装", DurationSec: 1, CompletedAt: time.Now()})
	_ = store.Add(domain.Session{Type: domain.WorkSession, Task: "読書", DurationSec: 1, CompletedAt: time.Now()})
	recorder := httptest.NewRecorder()

	(TaskHandler{Store: store}).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/tasks", nil))

	if body := recorder.Body.String(); !strings.Contains(body, "実装") || !strings.Contains(body, "読書") {
		t.Fatalf("body = %s", body)
	}
}

func TestExportHandlerReturnsCSV(t *testing.T) {
	store := storage.NewMemorySessionStore()
	_ = store.Add(domain.Session{ID: "session-1", Type: domain.WorkSession, DurationSec: 1500, CompletedAt: time.Now(), Task: "実装"})
	recorder := httptest.NewRecorder()

	ExportHandler(store).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/export", nil))

	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Header().Get("Content-Type"), "text/csv") {
		t.Fatalf("status = %d, content type = %q", recorder.Code, recorder.Header().Get("Content-Type"))
	}
	if !strings.Contains(recorder.Body.String(), "session-1") {
		t.Fatalf("CSV = %s", recorder.Body.String())
	}
}
