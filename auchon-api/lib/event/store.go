package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/auchon/auchon-api/lib/db"
)

// Store is the production Publisher: it appends events to the domain_events table of the
// database it is given, inside the transaction carried by ctx. Nothing consumes them yet.
//
//	domain_events(id TEXT PK, name TEXT, aggregate_id TEXT, occurred_at TEXT (RFC 3339 UTC), payload TEXT (JSON))
type Store struct {
	conn *sql.DB
}

func NewStore(conn *sql.DB) *Store { return &Store{conn: conn} }

func (s *Store) Publish(ctx context.Context, events ...Event) error {
	exec := db.From(ctx, s.conn)
	for _, e := range events {
		payload, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("encode %s: %w", e.EventName(), err)
		}
		_, err = exec.ExecContext(ctx,
			`INSERT INTO domain_events (id, name, aggregate_id, occurred_at, payload) VALUES (?, ?, ?, ?, ?)`,
			e.EventID(), e.EventName(), e.AggregateID(), e.OccurredAt().UTC().Format(time.RFC3339Nano), string(payload))
		if err != nil {
			return fmt.Errorf("store %s: %w", e.EventName(), err)
		}
	}
	return nil
}
