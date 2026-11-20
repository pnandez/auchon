# Auchon

Home inventory: what is stored where. Starts with the basement and grows to more places
(garage, attic…). Places contain nested locations (areas, shelves, boxes) that hold items with
quantities.

## Offline first

The basement has no signal. The phone is the working copy and the server is the sync point and
backup. Every change must work with no network and sync later; never build a feature that needs
the server to answer before the user sees the result.

## Layout

- `auchon-api/`: Go API + SQLite, hexagonal/DDD/CQRS. See `auchon-api/CLAUDE.md`
- `auchon-web/`: React + TypeScript PWA (Vite, Dexie). See `auchon-web/CLAUDE.md`
- Contract: `auchon-api/docs/openapi/openapi.yaml`. The web app's API types are generated from
  it, never hand-written
- Decisions: `auchon-api/docs/adr/`

## Rules

- Single user, no auth; runs on a home server on the home network
- **TDD, tests shown first**: write the tests for a change (feature scenarios, unit tests), see
  them fail, show them for approval, and only then implement
- Conventional commits (`feat:`, `fix:`, `test:`, `refactor:`, `docs:`, `chore:`)
