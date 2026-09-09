package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewHandlerServesIndexPage(t *testing.T) {
	handler, err := newHandler()
	if err != nil {
		t.Fatalf("newHandler() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusOK)
	}
	if contentType := recorder.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", contentType)
	}
	if body := recorder.Body.String(); !strings.Contains(body, "ポモドーロタイマー") {
		t.Fatalf("response body does not contain page title")
	}
}

func TestNewHandlerServesStaticAssets(t *testing.T) {
	handler, err := newHandler()
	if err != nil {
		t.Fatalf("newHandler() error = %v", err)
	}

	tests := []struct {
		name        string
		path        string
		contentType string
	}{
		{name: "stylesheet", path: "/style.css", contentType: "text/css"},
		{name: "javascript", path: "/js/app.js", contentType: "text/javascript"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			handler.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusOK {
				t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusOK)
			}
			if contentType := recorder.Header().Get("Content-Type"); !strings.HasPrefix(contentType, test.contentType) {
				t.Fatalf("Content-Type = %q, want %s", contentType, test.contentType)
			}
		})
	}
}

func TestNewHandlerReservesAPIPath(t *testing.T) {
	handler, err := newHandler()
	if err != nil {
		t.Fatalf("newHandler() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/", nil)
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotImplemented {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusNotImplemented)
	}
	body, err := io.ReadAll(recorder.Body)
	if err != nil {
		t.Fatalf("reading response body: %v", err)
	}
	if !strings.Contains(string(body), "API endpoint is not implemented") {
		t.Fatalf("response body does not describe the unimplemented API")
	}
}

func TestNewHandlerAcceptsSessionAndReturnsTodayStats(t *testing.T) {
	handler, err := newHandler()
	if err != nil {
		t.Fatalf("newHandler() error = %v", err)
	}

	sessionRequest := httptest.NewRequest(http.MethodPost, "/api/sessions", strings.NewReader(`{"type":"work","durationSec":1500,"completedAt":"2026-09-08T10:00:00Z"}`))
	sessionRecorder := httptest.NewRecorder()
	handler.ServeHTTP(sessionRecorder, sessionRequest)
	if sessionRecorder.Code != http.StatusCreated {
		t.Fatalf("session status = %d, want %d", sessionRecorder.Code, http.StatusCreated)
	}

	statsRequest := httptest.NewRequest(http.MethodGet, "/api/stats/today", nil)
	statsRecorder := httptest.NewRecorder()
	handler.ServeHTTP(statsRecorder, statsRequest)
	if statsRecorder.Code != http.StatusOK {
		t.Fatalf("stats status = %d, want %d", statsRecorder.Code, http.StatusOK)
	}
}

func TestNewHandlerReturnsNotFoundForUnknownPath(t *testing.T) {
	handler, err := newHandler()
	if err != nil {
		t.Fatalf("newHandler() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}
