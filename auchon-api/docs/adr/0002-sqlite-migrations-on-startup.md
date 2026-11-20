# 0002: SQLite, migrations applied on startup

Status: accepted (2026-11-20)

## Context

The API runs as one binary on a home server for a single user. The AI-skills migrations module
assumes a separate migrate step before deploys, built for multi-instance services with rolling
deploys.

## Decision

- SQLite via `modernc.org/sqlite` (pure Go, so cross-compiling for ARM needs no CGO).
- SQL migrations in `migrations/`, embedded with `embed.FS` and applied by `pressly/goose` when
  the binary starts, before it serves requests.
- Every up migration has a down; applied migrations are never edited.

## Consequences

- Deploying is replacing the binary (or image) and restarting.
- Tests open a temp SQLite file and migrate it the same way, so integration and acceptance tests
  need no Docker and run under `go test -short`.
- If the API ever runs as several instances, revisit this and adopt the migrations module.
