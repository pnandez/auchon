---
name: frontend-testing
description: Write Vitest + React Testing Library tests for auchon-web — commands, Dexie hooks, the sync engine, components and routes. Triggers on testing, test files, coverage, Vitest, RTL, unit or integration tests, or writing/reviewing tests in auchon-web.
---

# Frontend Testing (auchon-web)

Applies to `auchon-web/`. Run commands from that folder. Adapted from the AI-skills
`frontend-testing` skill. Does not apply to `auchon-api/` (Go, `go test`).

## TDD first

Tests are written **before** the implementation, run red, and **shown to the user for
approval**. Only then is the code written. Commit the approved tests (`test: …`), then the
implementation (`feat: …`).

## Commands

```bash
npm test                         # all tests once
npm run test:watch               # watch mode
npm test -- src/commands/addItem # one file or folder
npm test -- --coverage           # coverage report
```

## Files

- Tests sit **next to** the code: `addItem.ts` → `addItem.test.ts`,
  `PlaceCard.tsx` → `PlaceCard.test.tsx`
- Shared helpers in `src/test/`: `setup.ts` (jest-dom, `fake-indexeddb/auto`), and helpers
  such as `renderRoute()` and `resetDb()`

## What to test, and how

Order work from simple to complex, one file at a time (write → run → next):

1. **Commands** (`src/commands/`): given a local state, the command changes the tables
   **and** appends exactly one outbox envelope with the right `type` and payload, in one
   transaction. Invalid input throws the typed error and writes nothing
2. **Hooks** (`features/*/hooks/`): seed Dexie, render the hook with `renderHook`, assert the
   result; change the data and assert the hook updates (live query)
3. **Sync engine** (`src/sync/`): mock `fetch` only. Push sends the outbox in order and removes
   applied/duplicate entries; rejected ones become sync problems; pull applies changes and moves
   the cursor; offline (fetch throws) keeps the outbox intact
4. **Components**: render with real hooks against a seeded Dexie; interact with `userEvent`
5. **Routes / pages**: render through the real router with a memory history (below)

### Mocking policy

- ✅ Use the **real** Dexie database (fake-indexeddb). Reset it before each test
- ✅ Use real shadcn components, real child components, the real router
- ✅ Mock only the network (`fetch`) and time (`vi.useFakeTimers()`) when they matter
- ❌ Never mock Dexie, commands, hooks or sibling components to make a test pass

### Rendering a route

```tsx
import { render, screen } from '@testing-library/react'
import { RouterProvider, createMemoryHistory } from '@tanstack/react-router'
import { buildRouter } from '@/app/router'

function renderRoute(path: string) {
  const router = buildRouter({ history: createMemoryHistory({ initialEntries: [path] }) })
  return render(<RouterProvider router={router} />)
}

it('should list the basement on the places page', async () => {
  await db.locations.add(basement())
  renderRoute('/')
  expect(await screen.findByRole('link', { name: /basement/i })).toBeInTheDocument()
})
```

### Test data

- Builders with valid defaults and overrides, e.g. `aLocation({ name: 'Basement' })`,
  `anItem({ quantity: 3 })`, in `src/test/builders.ts`
- Never build test data through the sync pull path unless the pull is what is being tested

## Principles

- **Arrange / Act / Assert** in every test
- **Black box**: assert what the user sees or what lands in Dexie/outbox, not internal state
- **Semantic queries**: `getByRole` with `name`, `getByLabelText`, `getByText`, `within(...)`.
  `getByTestId` is a last resort; fix accessibility instead
- Async UI: `await screen.findBy…`, not assertions on the first render (live queries resolve
  after it)
- **One behaviour per test**, named `should <behaviour> when <condition>`
- Prefer regex text matching (`/basement/i`) over exact strings

## Required scenarios

- Commands: happy path, each validation failure, outbox envelope shape
- Components: renders, loading (`undefined`) state, empty state, each interaction
- Edge cases: empty lists, zero quantity, removed (tombstoned) records hidden, very long names

## Coverage

Aim for 100% of functions and statements and >95% of branches in `src/commands/` and
`src/sync/`; components and hooks are covered by behaviour, not by a number.
