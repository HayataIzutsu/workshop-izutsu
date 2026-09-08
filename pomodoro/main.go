package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"pomodoro/internal/clock"
	"pomodoro/internal/handler"
	"pomodoro/internal/storage"
)

//go:embed static
var staticFiles embed.FS

func newHandler() (http.Handler, error) {
	return newHandlerWithStore(storage.NewMemorySessionStore())
}

func newHandlerWithStore(store storage.SessionStore) (http.Handler, error) {
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(staticFS)))
	systemClock := clock.SystemClock{}
	mux.Handle("/api/sessions", handler.SessionHandler{Store: store, Clock: systemClock})
	mux.Handle("/api/stats/today", handler.StatsHandler{Store: store, Clock: systemClock, Location: time.Local})
	mux.Handle("/api/history", handler.HistoryHandler{Store: store})
	mux.Handle("/api/tasks", handler.TaskHandler{Store: store})
	mux.Handle("/api/stats/history", handler.StatsHistoryHandler{Store: store, Clock: systemClock, Location: time.Local})
	mux.Handle("/api/stats/tasks", handler.TaskStatsHandler{Store: store})
	mux.Handle("/api/sessions/delete", handler.DeleteSessionHandler{Store: store})
	mux.Handle("/api/export", handler.ExportHandler(store))
	mux.HandleFunc("/api/", func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "API endpoint is not implemented", http.StatusNotImplemented)
	})
	return mux, nil
}

func main() {
	dataPath := os.Getenv("POMODORO_DATA_FILE")
	if dataPath == "" {
		dataPath = "data/sessions.json"
	}
	store, err := storage.NewFileSessionStore(dataPath)
	if err != nil {
		log.Fatal(err)
	}
	handler, err := newHandlerWithStore(store)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := ":" + port
	log.Printf("listening on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(address, handler))
}
