package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"pomodoro/internal/clock"
	"pomodoro/internal/domain"
	"pomodoro/internal/storage"
)

type SessionHandler struct {
	Store storage.SessionStore
	Clock clock.Clock
}

func (handler SessionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var session domain.Session
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&session); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := session.Validate(); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	if session.ID == "" {
		session.ID = session.CompletedAt.UTC().Format("20060102T150405.000000000Z07:00") + "-" + time.Now().Format("150405.000000000")
	}
	if err := handler.Store.Add(session); err != nil {
		writeError(writer, http.StatusInternalServerError, "could not store session")
		return
	}
	writeJSON(writer, http.StatusCreated, session)
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

func writeError(writer http.ResponseWriter, status int, message string) {
	writeJSON(writer, status, map[string]string{"error": message})
}
