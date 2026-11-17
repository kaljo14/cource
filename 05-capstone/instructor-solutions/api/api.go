package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

type Note struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}
type Store interface {
	Ping(context.Context) error
	List(context.Context) ([]Note, error)
	Add(context.Context, string) (Note, error)
}

func Handler(store Store, version string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok\n")) })
	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(version + "\n")) })
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := store.Ping(ctx); err != nil {
			http.Error(w, "database unavailable", 503)
			return
		}
		w.Write([]byte("ready\n"))
	})
	mux.HandleFunc("GET /api/notes", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		notes, err := store.List(ctx)
		if err != nil {
			slog.Error("list notes", "error", err)
			http.Error(w, "database unavailable", 503)
			return
		}
		writeJSON(w, 200, notes)
	})
	mux.HandleFunc("POST /api/notes", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 4096)
		var input struct {
			Text string `json:"text"`
		}
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&input); err != nil {
			http.Error(w, "expected JSON with text", 400)
			return
		}
		if err := dec.Decode(&struct{}{}); err != io.EOF {
			http.Error(w, "expected one JSON object", 400)
			return
		}
		input.Text = strings.TrimSpace(input.Text)
		if n := utf8.RuneCountInString(input.Text); n < 1 || n > 200 {
			http.Error(w, "text must contain 1 to 200 characters", 400)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		note, err := store.Add(ctx, input.Text)
		if err != nil {
			slog.Error("add note", "error", err)
			http.Error(w, "database unavailable", 503)
			return
		}
		writeJSON(w, 201, note)
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request", "method", r.Method, "path", r.URL.Path)
		mux.ServeHTTP(w, r)
	})
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}
