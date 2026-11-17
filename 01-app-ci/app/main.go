package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var version = "dev"

func env(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
func handler(ready *atomic.Bool) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"name": "go-app", "message": env("MESSAGE", "Hello students"), "version": version})
	})
	m.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(version + "\n")) })
	m.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok\n")) })
	m.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		if !ready.Load() {
			http.Error(w, "not ready", 503)
			return
		}
		w.Write([]byte("ready\n"))
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request", "method", r.Method, "path", r.URL.Path)
		m.ServeHTTP(w, r)
	})
}
func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	var ready atomic.Bool
	ready.Store(true)
	srv := &http.Server{Addr: env("ADDR", ":8080"), Handler: handler(&ready), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	done := make(chan error, 1)
	go func() { done <- srv.ListenAndServe() }()
	select {
	case err := <-done:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
		return
	case <-ctx.Done():
	}
	ready.Store(false)
	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdown); err != nil {
		slog.Error("shutdown", "error", err)
		srv.Close()
	}
}
