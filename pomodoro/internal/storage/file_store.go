package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"pomodoro/internal/domain"
)

type FileSessionStore struct {
	mu       sync.RWMutex
	path     string
	sessions []domain.Session
}

func NewFileSessionStore(path string) (*FileSessionStore, error) {
	store := &FileSessionStore{path: path}
	if err := store.load(); err != nil {
		return nil, err
	}
	return store, nil
}

func (store *FileSessionStore) Add(session domain.Session) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.sessions = append(store.sessions, session)
	return store.persistLocked()
}

func (store *FileSessionStore) List() ([]domain.Session, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	sessions := make([]domain.Session, len(store.sessions))
	copy(sessions, store.sessions)
	return sessions, nil
}

func (store *FileSessionStore) Delete(id string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	for index, session := range store.sessions {
		if session.ID == id {
			store.sessions = append(store.sessions[:index], store.sessions[index+1:]...)
			return store.persistLocked()
		}
	}
	return nil
}

func (store *FileSessionStore) load() error {
	data, err := os.ReadFile(store.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &store.sessions)
}

func (store *FileSessionStore) persistLocked() error {
	if err := os.MkdirAll(filepath.Dir(store.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(store.sessions, "", "  ")
	if err != nil {
		return err
	}
	temporaryPath := store.path + ".tmp"
	if err := os.WriteFile(temporaryPath, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(temporaryPath, store.path); err != nil {
		return fmt.Errorf("replace session file: %w", err)
	}
	return nil
}
