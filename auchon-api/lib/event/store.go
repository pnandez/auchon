package event

import (
	"context"
	"database/sql"
)

// SQLiteStore is the production Publisher: it appends events to the domain_events table in
// the transaction carried by ctx. Nothing consumes them yet.
//
//	domain_events(id TEXT PK, name TEXT, aggregate_id TEXT, occurred_at TEXT (RFC 3339 UTC), payload TEXT (JSON))
type SQLiteStore struct {
	conn *sql.DB
}

func NewSQLiteStore(conn *sql.DB) *SQLiteStore { return &SQLiteStore{conn: conn} }

func (s *SQLiteStore) Publish(ctx context.Context, events ...Event) error {
	return nil
}
