package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexServed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	newMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d but got %d", http.StatusOK, rr.Code)
	}
}

func TestGamificationScriptServed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/js/gamification.js", nil)
	rr := httptest.NewRecorder()

	newMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d but got %d", http.StatusOK, rr.Code)
	}
}
