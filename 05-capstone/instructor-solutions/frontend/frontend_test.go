package main

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func TestThemeReloadAndProxy(t *testing.T) {
	target, _ := url.Parse("http://api:8080")
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host != "api:8080" || r.URL.Path != "/api/notes" {
			t.Errorf("wrong upstream: %s", r.URL)
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("[]"))}, nil
	})
	file := filepath.Join(t.TempDir(), "BACKGROUND_COLOR")
	h := handlerWithTransport(target, file, transport)
	for _, color := range []string{"#e0f2fe", "#fef3c7"} {
		if err := os.WriteFile(file, []byte(color), 0600); err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", "/config", nil))
		if !strings.Contains(w.Body.String(), color) {
			t.Fatal(w.Body.String())
		}
		if w.Header().Get("Cache-Control") != "no-store" {
			t.Fatal("config can be cached")
		}
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/api/notes", nil))
	if w.Code != 200 || w.Body.String() != "[]" {
		t.Fatal(w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "Campus Board") {
		t.Fatal(w.Code)
	}
}
func TestInvalidColorAndUnavailableAPI(t *testing.T) {
	target, _ := url.Parse("http://api:8080")
	file := filepath.Join(t.TempDir(), "color")
	os.WriteFile(file, []byte("not a color"), 0600)
	h := handlerWithTransport(target, file, roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, errors.New("unavailable") }))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/config", nil))
	if !strings.Contains(w.Body.String(), "#e0f2fe") {
		t.Fatal(w.Body.String())
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/api/notes", nil))
	if w.Code != 503 {
		t.Fatal(w.Code)
	}
}
