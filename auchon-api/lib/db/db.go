// Package db is the kernel's SQLite access: opening, migrating, and the transaction carried
// in the request context.
package db

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

// Executor is what adapters run SQL on: the transaction in the context, or the database.
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Open opens the SQLite file at path with foreign keys on, WAL and a busy timeout. SQLite has
// one writer, so the pool holds one connection: transactions queue instead of failing busy.
func Open(path string) (*sql.DB, error) {
	dsn := "file:" + path + "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_txlock=immediate"
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	conn.SetMaxOpenConns(1)
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	return conn, nil
}

// Migrate applies the goose migrations in fsys (SQL files at its root).
func Migrate(ctx context.Context, conn *sql.DB, fsys fs.FS) error {
	provider, err := goose.NewProvider(goose.DialectSQLite3, conn, fsys)
	if err != nil {
		return fmt.Errorf("migrations: %w", err)
	}
	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

type txKey struct{}

// WithTx runs fn in a transaction carried by ctx: committed when fn returns nil, rolled back
// otherwise. If ctx already carries one, fn joins it.
func WithTx(ctx context.Context, conn *sql.DB, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return fn(ctx)
	}
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	if err := fn(context.WithValue(ctx, txKey{}, tx)); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

// From returns the transaction in ctx, or conn when there is none.
func From(ctx context.Context, conn *sql.DB) Executor {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return conn
}
