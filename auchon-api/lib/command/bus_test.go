package command_test

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/auchon/auchon-api/lib/command"
	"github.com/auchon/auchon-api/lib/db"
	"github.com/auchon/auchon-api/lib/test"
)

type labelBox struct{ ID, Label string }

func (labelBox) CommandName() string { return "test.LabelBox" }

type labelBoxHandler struct {
	conn     *sql.DB
	received []labelBox
	fail     error
}

func (h *labelBoxHandler) Execute(ctx context.Context, cmd labelBox) error {
	h.received = append(h.received, cmd)
	if _, err := db.From(ctx, h.conn).ExecContext(ctx, `INSERT INTO boxes VALUES (?, ?)`, cmd.ID, cmd.Label); err != nil {
		return err
	}
	return h.fail
}

func setup(t *testing.T) (*sql.DB, *command.Bus, *labelBoxHandler) {
	t.Helper()
	conn, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	test.AssertNoError(t, err)
	t.Cleanup(func() { conn.Close() })
	test.AssertNoError(t, db.Migrate(context.Background(), conn, fstest.MapFS{
		"00001_boxes.sql": {Data: []byte("-- +goose Up\nCREATE TABLE boxes (id TEXT PRIMARY KEY, label TEXT);\n-- +goose Down\nDROP TABLE boxes;\n")},
	}))
	bus := command.NewBus(conn)
	h := &labelBoxHandler{conn: conn}
	command.Register[labelBox](bus, h)
	return conn, bus, h
}

func boxes(t *testing.T, conn *sql.DB) int {
	var n int
	test.AssertNoError(t, conn.QueryRow(`SELECT COUNT(*) FROM boxes`).Scan(&n))
	return n
}

func TestBus(t *testing.T) {
	ctx := context.Background()

	t.Run("Given a registered handler", func(t *testing.T) {
		t.Run("When its command is dispatched", func(t *testing.T) {
			conn, bus, h := setup(t)
			err := bus.Dispatch(ctx, labelBox{ID: "b1", Label: "Cables"})

			t.Run("Then the handler receives the typed command", func(t *testing.T) {
				test.AssertNoError(t, err)
				if len(h.received) != 1 || h.received[0].Label != "Cables" {
					t.Errorf("received = %v", h.received)
				}
			})

			t.Run("And its writes are committed", func(t *testing.T) {
				if n := boxes(t, conn); n != 1 {
					t.Errorf("boxes = %d, want 1", n)
				}
			})
		})

		t.Run("When the handler fails", func(t *testing.T) {
			conn, bus, h := setup(t)
			h.fail = errors.New("box is full")
			err := bus.Dispatch(ctx, labelBox{ID: "b1", Label: "Cables"})

			t.Run("Then the error is returned and the writes are rolled back", func(t *testing.T) {
				if !errors.Is(err, h.fail) {
					t.Errorf("err = %v, want %v", err, h.fail)
				}
				if n := boxes(t, conn); n != 0 {
					t.Errorf("boxes = %d, want 0", n)
				}
			})
		})

		t.Run("When the same command is registered again, Then it panics", func(t *testing.T) {
			_, bus, h := setup(t)
			defer func() {
				if recover() == nil {
					t.Error("expected a panic")
				}
			}()
			command.Register[labelBox](bus, h)
		})
	})

	t.Run("Given a command with no handler, When dispatched, Then it returns NoHandlerError", func(t *testing.T) {
		_, bus, _ := setup(t)
		err := bus.Dispatch(ctx, unknownCommand{})
		test.AssertErrorIsOfType[command.NoHandlerError](t, err)
	})
}

type unknownCommand struct{}

func (unknownCommand) CommandName() string { return "test.Unknown" }
