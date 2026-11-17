package main

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeStore struct {
	fail  bool
	notes []Note
}

func (f *fakeStore) Ping(context.Context) error {
	if f.fail {
		return errors.New("down")
	}
	return nil
}
func (f *fakeStore) List(context.Context) ([]Note, error) {
	if f.fail {
		return nil, errors.New("down")
	}
	return f.notes, nil
}
func (f *fakeStore) Add(_ context.Context, text string) (Note, error) {
	if f.fail {
		return Note{}, errors.New("down")
	}
	n := Note{ID: 1, Text: text}
	f.notes = append(f.notes, n)
	return n, nil
}
func TestCreateAndList(t *testing.T) {
	f := &fakeStore{notes: []Note{}}
	h := Handler(f, "test")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("POST", "/api/notes", strings.NewReader(`{"text":" hello "}`)))
	if w.Code != 201 || f.notes[0].Text != "hello" {
		t.Fatal(w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/api/notes", nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "hello") {
		t.Fatal(w.Body.String())
	}
}
func TestBadInput(t *testing.T) {
	for _, body := range []string{`{`, `{"text":" "}`, `{"text":"x","unknown":1}`, `{"text":"x"} {}`, `{"text":"` + strings.Repeat("x", 201) + `"}`} {
		w := httptest.NewRecorder()
		Handler(&fakeStore{}, "test").ServeHTTP(w, httptest.NewRequest("POST", "/api/notes", strings.NewReader(body)))
		if w.Code != 400 {
			t.Errorf("body %q: %d", body, w.Code)
		}
	}
}
func TestDatabaseFailure(t *testing.T) {
	h := Handler(&fakeStore{fail: true}, "test")
	for path, code := range map[string]int{"/healthz": 200, "/readyz": 503, "/api/notes": 503} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != code {
			t.Errorf("%s: %d", path, w.Code)
		}
	}
}
