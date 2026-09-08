package storage

import (
	"sync"

	"pomodoro/internal/domain"
)

type SessionStore interface {
	Add(session domain.Session) error
	List() ([]domain.Session, error)
	Delete(id string) error
}

type MemorySessionStore struct {
	mu       sync.RWMutex
	sessions []domain.Session
}

func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{}
}

func (store *MemorySessionStore) Add(session domain.Session) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.sessions = append(store.sessions, session)
	return nil
}

func (store *MemorySessionStore) List() ([]domain.Session, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	sessions := make([]domain.Session, len(store.sessions))
	copy(sessions, store.sessions)
	return sessions, nil
}

func (store *MemorySessionStore) Delete(id string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	for index, session := range store.sessions {
		if session.ID == id {
			store.sessions = append(store.sessions[:index], store.sessions[index+1:]...)
			return nil
		}
	}
	return nil
}
