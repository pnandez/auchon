# 0001: Offline sync with a command outbox and a change-sequence pull

Status: accepted (2026-11-20)

## Context

Auchon is used mostly in a basement with no signal. Every change must work offline on the phone
and reach the server later. There is a single user, sometimes on more than one device.

## Decision

- **Client-generated ids** (UUIDv7) for locations and items, so creation needs no server.
- **Command outbox on the client**: each user action is a command envelope
  `{commandId, type, payload, createdAt}`. It is applied optimistically to the local IndexedDB
  read model and appended to an outbox in the same Dexie transaction.
- **Push**: `POST /api/sync/push` takes the outbox in order. The sync controller (UI layer,
  anti-corruption) maps each envelope to an application command and dispatches it on the command
  bus. Each result is `applied`, `duplicate` or `rejected` (with the error category).
- **Idempotency**: `processed_commands(command_id, result)` is written in the same transaction
  as the command, so retrying a batch never applies anything twice.
- **Pull**: `GET /api/sync/pull?since=<seq>` returns locations and items changed after `seq`,
  tombstones included, plus the new cursor. Every write stamps the row with the next value of a
  server-wide change sequence. It is a read model: no entities are loaded.
- **Conflicts**: commands are applied in arrival order. Name, notes and moves are last write
  wins; quantities are deltas, so concurrent changes add up.
- **Deletes are soft** (tombstones) so other devices learn about them on pull.

## Consequences

- The client never waits for the server; the UI reads only from IndexedDB.
- Rejected commands (e.g. removing a box that still holds items) are shown to the user. After
  the pull, the local state matches the server again.
- Every command must be safe to receive twice and to be validated late, long after the user
  made the change.
- No event consumers are needed: the pull is served from read-model columns, not from events.
