package event_test

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/auchon/auchon-api/lib/db"
	"github.com/auchon/auchon-api/lib/event"
	"github.com/auchon/auchon-api/lib/test"
)

var eventsSchema = fstest.MapFS{"00001_events.sql": {Data: []byte(`
-- +goose Up
CREATE TABLE domain_events (
  id TEXT PRIMARY KEY, name TEXT NOT NULL, aggregate_id TEXT NOT NULL,
  occurred_at TEXT NOT NULL, payload TEXT NOT NULL
);
-- +goose Down
DROP TABLE domain_events;
`)}}

func openStore(t *testing.T) (*sql.DB, *event.Store) {
	t.Helper()
	conn, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	test.AssertNoError(t, err)
	t.Cleanup(func() { conn.Close() })
	test.AssertNoError(t, db.Migrate(context.Background(), conn, eventsSchema))
	return conn, event.NewStore(conn)
}

func TestStore(t *testing.T) {
	ctx := context.Background()

	t.Run("Given two events, When they are published", func(t *testing.T) {
		conn, store := openStore(t)
		first := boxLabelled{Base: event.NewBase("box-1", at), Label: "Cables"}
		second := boxLabelled{Base: event.NewBase("box-2", at), Label: "Tools"}
		test.AssertNoError(t, store.Publish(ctx, first, second))

		t.Run("Then each is stored with its name, aggregate, moment and payload", func(t *testing.T) {
			var name, aggregate, occurred, payload string
			err := conn.QueryRow(`SELECT name, aggregate_id, occurred_at, payload FROM domain_events WHERE id = ?`,
				first.EventID()).Scan(&name, &aggregate, &occurred, &payload)
			test.AssertNoError(t, err)
			if name != "test.BoxLabelled" || aggregate != "box-1" || occurred != "2026-11-20T18:00:00Z" {
				t.Errorf("row = %q %q %q", name, aggregate, occurred)
			}
			if payload != `{"Label":"Cables"}` {
				t.Errorf("payload = %s", payload)
			}
		})

		t.Run("And both are stored", func(t *testing.T) {
			var n int
			conn.QueryRow(`SELECT COUNT(*) FROM domain_events`).Scan(&n)
			if n != 2 {
				t.Errorf("events = %d, want 2", n)
			}
		})
	})

	t.Run("Given a transaction that fails after publishing, Then the events are rolled back", func(t *testing.T) {
		conn, store := openStore(t)
		db.WithTx(ctx, conn, func(ctx context.Context) error {
			store.Publish(ctx, boxLabelled{Base: event.NewBase("box-1", at), Label: "Cables"})
			return errors.New("store failed")
		})
		var n int
		conn.QueryRow(`SELECT COUNT(*) FROM domain_events`).Scan(&n)
		if n != 0 {
			t.Errorf("events = %d, want 0", n)
		}
	})
}
