package handler

import (
	"encoding/csv"
	"net/http"
	"sort"
	"strconv"
	"time"

	"pomodoro/internal/storage"
)

type HistoryHandler struct {
	Store storage.SessionStore
}

func (handler HistoryHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	sessions, err := handler.Store.List()
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "could not load sessions")
		return
	}
	if task := request.URL.Query().Get("task"); task != "" {
		filtered := sessions[:0]
		for _, session := range sessions {
			if session.Task == task {
				filtered = append(filtered, session)
			}
		}
		sessions = filtered
	}
	sort.Slice(sessions, func(left, right int) bool {
		return sessions[left].CompletedAt.After(sessions[right].CompletedAt)
	})
	writeJSON(writer, http.StatusOK, sessions)
}

type TaskHandler struct {
	Store storage.SessionStore
}

func (handler TaskHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	sessions, err := handler.Store.List()
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "could not load sessions")
		return
	}
	seen := map[string]bool{}
	tasks := make([]string, 0)
	for _, session := range sessions {
		if session.Task != "" && !seen[session.Task] {
			seen[session.Task] = true
			tasks = append(tasks, session.Task)
		}
	}
	sort.Strings(tasks)
	writeJSON(writer, http.StatusOK, tasks)
}

func ExportHandler(store storage.SessionStore) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		sessions, err := store.List()
		if err != nil {
			writeError(writer, http.StatusInternalServerError, "could not load sessions")
			return
		}
		writer.Header().Set("Content-Type", "text/csv; charset=utf-8")
		writer.Header().Set("Content-Disposition", "attachment; filename=sessions.csv")
		csvWriter := csv.NewWriter(writer)
		_ = csvWriter.Write([]string{"id", "type", "durationSec", "completedAt", "task", "note"})
		for _, session := range sessions {
			_ = csvWriter.Write([]string{session.ID, string(session.Type), strconv.Itoa(session.DurationSec), session.CompletedAt.Format(time.RFC3339), session.Task, session.Note})
		}
		csvWriter.Flush()
	})
}
