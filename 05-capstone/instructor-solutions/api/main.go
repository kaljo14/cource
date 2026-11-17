package main

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"time"
)

var version = "dev"

func main() {
	// lib/pq reads PGHOST, PGPORT, PGUSER, PGPASSWORD and PGDATABASE from environment.
	db, err := sql.Open("postgres", "sslmode=disable connect_timeout=3")
	if err != nil {
		slog.Error("database setup failed")
		os.Exit(1)
	}
	defer db.Close()
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)
	// Idempotent, simple classroom migration. Real projects should version migrations.
	initialized := false
	for attempt := 0; attempt < 60; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS notes (id BIGSERIAL PRIMARY KEY, text VARCHAR(200) NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now())`)
		cancel()
		if err == nil {
			initialized = true
			break
		}
		slog.Info("waiting for database", "attempt", attempt+1)
		time.Sleep(2 * time.Second)
	}
	if !initialized {
		slog.Error("database initialization timed out")
		os.Exit(1)
	}
	Serve(Handler(Postgres{DB: db}, version))
}
