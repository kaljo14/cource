package main

import (
	"context"
	"database/sql"
)

type Postgres struct{ DB *sql.DB }

func (p Postgres) Ping(ctx context.Context) error { return p.DB.PingContext(ctx) }
func (p Postgres) List(ctx context.Context) ([]Note, error) {
	rows, err := p.DB.QueryContext(ctx, "SELECT id, text, created_at FROM notes ORDER BY id DESC LIMIT 100")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	notes := []Note{}
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Text, &n.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, rows.Err()
}
func (p Postgres) Add(ctx context.Context, text string) (Note, error) {
	var n Note
	err := p.DB.QueryRowContext(ctx, "INSERT INTO notes(text) VALUES ($1) RETURNING id,text,created_at", text).Scan(&n.ID, &n.Text, &n.CreatedAt)
	return n, err
}
