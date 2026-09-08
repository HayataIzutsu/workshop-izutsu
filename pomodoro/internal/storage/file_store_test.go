package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"pomodoro/internal/domain"
)

func TestFileSessionStorePersistsAndReloadsSessions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sessions.json")
	session := domain.Session{Type: domain.WorkSession, DurationSec: 1500, CompletedAt: time.Now().UTC()}
	store, err := NewFileSessionStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Add(session); err != nil {
		t.Fatal(err)
	}

	reloaded, err := NewFileSessionStore(path)
	if err != nil {
		t.Fatal(err)
	}
	sessions, err := reloaded.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 || sessions[0].DurationSec != 1500 {
		t.Fatalf("sessions = %+v", sessions)
	}
}

func TestFileSessionStoreRejectsInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sessions.json")
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := NewFileSessionStore(path); err == nil {
		t.Fatal("NewFileSessionStore() returned nil error for invalid JSON")
	}
}
