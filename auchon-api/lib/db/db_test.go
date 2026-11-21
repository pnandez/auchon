package db_test

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/auchon/auchon-api/lib/db"
)

var migrations = fstest.MapFS{
	"00001_boxes.sql": {Data: []byte(`
-- +goose Up
CREATE TABLE boxes (id TEXT PRIMARY KEY, label TEXT NOT NULL);
-- +goose Down
DROP TABLE boxes;
`)},
}

func openMigrated(t *testing.T) *sql.DB {
	t.Helper()
	conn, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if conn == nil {
		t.Fatal("Open returned no connection")
	}
	t.Cleanup(func() { conn.Close() })
	if err := db.Migrate(context.Background(), conn, migrations); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	return conn
}

func countBoxes(t *testing.T, conn *sql.DB) int {
	t.Helper()
	var n int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM boxes`).Scan(&n); err != nil {
		t.Fatalf("count: %v", err)
	}
	return n
}

func TestOpenAndMigrate(t *testing.T) {
	t.Run("Given a new database file", func(t *testing.T) {
		conn := openMigrated(t)

		t.Run("Then foreign keys are enforced", func(t *testing.T) {
			var on int
			conn.QueryRow(`PRAGMA foreign_keys`).Scan(&on)
			if on != 1 {
				t.Errorf("foreign_keys = %d, want 1", on)
			}
		})

		t.Run("Then the migrations are applied", func(t *testing.T) {
			if n := countBoxes(t, conn); n != 0 {
				t.Errorf("boxes = %d, want 0", n)
			}
		})

		t.Run("When migrating again, Then nothing fails", func(t *testing.T) {
			if err := db.Migrate(context.Background(), conn, migrations); err != nil {
				t.Errorf("second Migrate: %v", err)
			}
		})
	})
}

func TestWithTx(t *testing.T) {
	ctx := context.Background()

	t.Run("Given fn succeeds, Then its writes are committed", func(t *testing.T) {
		conn := openMigrated(t)
		err := db.WithTx(ctx, conn, func(ctx context.Context) error {
			_, err := db.From(ctx, conn).ExecContext(ctx, `INSERT INTO boxes VALUES ('b1', 'Cables')`)
			return err
		})
		if err != nil {
			t.Fatalf("WithTx: %v", err)
		}
		if n := countBoxes(t, conn); n != 1 {
			t.Errorf("boxes = %d, want 1", n)
		}
	})

	t.Run("Given fn fails, Then its writes are rolled back and the error returned", func(t *testing.T) {
		conn := openMigrated(t)
		boom := errors.New("boom")
		err := db.WithTx(ctx, conn, func(ctx context.Context) error {
			db.From(ctx, conn).ExecContext(ctx, `INSERT INTO boxes VALUES ('b1', 'Cables')`)
			return boom
		})
		if !errors.Is(err, boom) {
			t.Errorf("err = %v, want boom", err)
		}
		if n := countBoxes(t, conn); n != 0 {
			t.Errorf("boxes = %d, want 0", n)
		}
	})

	t.Run("Given a nested WithTx whose outer fn fails, Then the inner writes are rolled back too", func(t *testing.T) {
		conn := openMigrated(t)
		db.WithTx(ctx, conn, func(ctx context.Context) error {
			db.WithTx(ctx, conn, func(ctx context.Context) error {
				_, err := db.From(ctx, conn).ExecContext(ctx, `INSERT INTO boxes VALUES ('b1', 'Cables')`)
				return err
			})
			return errors.New("outer fails")
		})
		if n := countBoxes(t, conn); n != 0 {
			t.Errorf("boxes = %d, want 0", n)
		}
	})

	t.Run("Given no transaction in ctx, Then From returns the database", func(t *testing.T) {
		conn := openMigrated(t)
		if got := db.From(ctx, conn); got != db.Executor(conn) {
			t.Errorf("From() = %T, want *sql.DB", got)
		}
	})
}
