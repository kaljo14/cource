package main

import (
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestEndpoints(t *testing.T) {
	var ready atomic.Bool
	ready.Store(true)
	h := handler(&ready)
	for _, p := range []string{"/", "/healthz", "/readyz", "/version"} {
		r := httptest.NewRecorder()
		h.ServeHTTP(r, httptest.NewRequest("GET", p, nil))
		if r.Code != 200 {
			t.Errorf("%s: %d", p, r.Code)
		}
	}
	ready.Store(false)
	r := httptest.NewRecorder()
	h.ServeHTTP(r, httptest.NewRequest("GET", "/readyz", nil))
	if r.Code != 503 {
		t.Fatal(r.Code)
	}
	t.Setenv("MESSAGE", "test message")
	r = httptest.NewRecorder()
	h.ServeHTTP(r, httptest.NewRequest("GET", "/", nil))
	if !strings.Contains(r.Body.String(), "test message") {
		t.Fatal(r.Body.String())
	}
}
