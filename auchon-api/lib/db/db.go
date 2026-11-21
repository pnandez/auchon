// Package db is the kernel's SQLite access: opening, migrating, and the transaction carried
// in the request context.
package db

import (
	"context"
	"database/sql"
	"io/fs"
)

// Executor is what adapters run SQL on: the transaction in the context, or the database.
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Open opens the SQLite file at path with foreign keys on, WAL and a busy timeout.
func Open(path string) (*sql.DB, error) {
	return nil, nil
}

// Migrate applies the goose migrations in fsys (SQL files at its root).
func Migrate(ctx context.Context, conn *sql.DB, fsys fs.FS) error {
	return nil
}

// WithTx runs fn in a transaction carried by ctx: committed when fn returns nil, rolled back
// otherwise. If ctx already carries one, fn joins it.
func WithTx(ctx context.Context, conn *sql.DB, fn func(ctx context.Context) error) error {
	return nil
}

// From returns the transaction in ctx, or conn when there is none.
func From(ctx context.Context, conn *sql.DB) Executor {
	return conn
}
