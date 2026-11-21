# auchon-web

Auchon's offline-first PWA. Product rules live in the root `../CLAUDE.md`; the sync design in
`../auchon-api/docs/adr/0001-*.md`.

## Stack

- React 19 + TypeScript, Vite, `vite-plugin-pwa` (service worker, manifest)
- TanStack Router, file-based (`src/routes/`, generated `src/routeTree.gen.ts`)
- shadcn/ui + Tailwind v4 (`src/components/ui/`, `cn()` in `src/lib/utils.ts`)
- Dexie + `dexie-react-hooks` for IndexedDB
- Vitest + React Testing Library + `fake-indexeddb`

## Structure

- `src/db/`: Dexie schema: `locations`, `items` (local read model), `outbox` (pending
  commands), `meta` (sync cursor)
- `src/commands/`: one module per command. Each applies the change to the local tables and
  appends its envelope to `outbox` in **one Dexie transaction**
- `src/sync/`: push/pull engine and sync status
- `src/features/<feature>/`: screens and components (locations, items, search)
- `src/routes/`: TanStack Router route files (thin: params + the feature's page)
- `src/app/`: router setup, shell, sync badge
- `src/components/ui/`: shadcn primitives (`npx shadcn@latest add <name>`)
- `src/api/`: types generated from `../auchon-api/docs/openapi/openapi.yaml`

## Rules

- **Offline first**: components read only from Dexie (`useLiveQuery`) and never call the API.
  Only `src/sync/` talks to the server
- Every mutation is a command through `src/commands/`, never a direct table write from a
  component
- Ids are UUIDv7 generated here (`uuid` package, `v7()`)
- Quantity changes are deltas (`+1`, `-2`), never "set to N"
- API types are generated, never hand-written
- Mobile first: large tap targets, works one-handed
- **TDD, tests shown first**: Vitest + React Testing Library tests are written, run red and
  shown for approval before any implementation. IndexedDB is `fake-indexeddb` in tests

## Skills

- `.claude/skills/frontend-dev-guidelines/`: component, page, routing, styling and data rules
- `.claude/skills/frontend-testing/`: how to write the Vitest + RTL tests
- Both are adapted from AI-skills for this stack; follow these copies, not the originals

## Commands

- Dev: `npm run dev` (proxies `/api` to `localhost:8080`)
- Tests: `npm test` (watch: `npm run test:watch`)
- Build: `npm run build` (outputs `dist/` with the service worker)
- Lint: `npm run lint`
