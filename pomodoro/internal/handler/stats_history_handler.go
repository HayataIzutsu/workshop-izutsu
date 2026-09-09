package handler

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"pomodoro/internal/clock"
	"pomodoro/internal/domain"
	"pomodoro/internal/storage"
)

type StatsHistoryHandler struct {
	Store    storage.SessionStore
	Clock    clock.Clock
	Location *time.Location
}

func (handler StatsHistoryHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	days := 7
	if value, err := strconv.Atoi(request.URL.Query().Get("days")); err == nil && value > 0 && value <= 366 {
		days = value
	}
	sessions, err := handler.Store.List()
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "could not load sessions")
		return
	}
	location := handler.Location
	if location == nil {
		location = time.UTC
	}
	now := handler.Clock.Now().In(location)
	type dailyStats struct {
		Date string `json:"date"`
		domain.Stats
	}
	history := make([]dailyStats, 0, days)
	for offset := 0; offset < days; offset++ {
		day := now.AddDate(0, 0, -offset)
		history = append(history, dailyStats{Date: day.Format("2006-01-02"), Stats: domain.StatsForDay(sessions, day, location)})
	}
	sort.Slice(history, func(left, right int) bool { return history[left].Date < history[right].Date })
	writeJSON(writer, http.StatusOK, history)
}

type TaskStatsHandler struct {
	Store storage.SessionStore
}

func (handler TaskStatsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	sessions, err := handler.Store.List()
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "could not load sessions")
		return
	}
	type taskStats struct {
		Task         string `json:"task"`
		Sessions     int    `json:"sessions"`
		FocusTimeSec int    `json:"focusTimeSec"`
	}
	statsByTask := map[string]*taskStats{}
	for _, session := range sessions {
		if session.Type != domain.WorkSession || session.Task == "" {
			continue
		}
		stats := statsByTask[session.Task]
		if stats == nil {
			stats = &taskStats{Task: session.Task}
			statsByTask[session.Task] = stats
		}
		stats.Sessions++
		stats.FocusTimeSec += session.DurationSec
	}
	result := make([]taskStats, 0, len(statsByTask))
	for _, stats := range statsByTask {
		result = append(result, *stats)
	}
	sort.Slice(result, func(left, right int) bool { return result[left].Task < result[right].Task })
	writeJSON(writer, http.StatusOK, result)
}

type DeleteSessionHandler struct {
	Store storage.SessionStore
}

func (handler DeleteSessionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id := request.URL.Query().Get("id")
	if id == "" {
		writeError(writer, http.StatusBadRequest, "id is required")
		return
	}
	if err := handler.Store.Delete(id); err != nil {
		writeError(writer, http.StatusInternalServerError, "could not delete session")
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}
