package handler

import (
	"net/http"
	"time"

	"pomodoro/internal/clock"
	"pomodoro/internal/domain"
	"pomodoro/internal/storage"
)

type StatsHandler struct {
	Store    storage.SessionStore
	Clock    clock.Clock
	Location *time.Location
}

func (handler StatsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	sessions, err := handler.Store.List()
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "could not load sessions")
		return
	}
	stats := domain.StatsForDay(sessions, handler.Clock.Now(), handler.Location)
	writeJSON(writer, http.StatusOK, stats)
}
