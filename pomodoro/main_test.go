package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandlerReturnsHTML(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	rootHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	if got := rr.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Fatalf("unexpected content-type: %q", got)
	}

	body := rr.Body.String()
	checks := []string{
		"id=\"progressRing\"",
		"id=\"particleCanvas\"",
		"id=\"rippleLayer\"",
		"function getProgressColor",
		"requestAnimationFrame(tick)",
	}

	for _, check := range checks {
		if !strings.Contains(body, check) {
			t.Fatalf("response does not include %q", check)
		}
	}
}
