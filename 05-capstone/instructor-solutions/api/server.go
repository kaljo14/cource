package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Serve(handler http.Handler) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	server := &http.Server{Addr: ":8080", Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	done := make(chan error, 1)
	go func() { done <- server.ListenAndServe() }()
	select {
	case err := <-done:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("serve", "error", err)
			os.Exit(1)
		}
		return
	case <-ctx.Done():
	}
	drain, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(drain); err != nil {
		slog.Error("shutdown", "error", err)
		server.Close()
	}
}
