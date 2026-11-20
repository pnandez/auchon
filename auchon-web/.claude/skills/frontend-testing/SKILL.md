---
name: frontend-testing
description: Generate Vitest + React Testing Library tests for inmobil-app/app components, hooks, and utilities. Triggers on testing, spec files, coverage, Vitest, RTL, unit tests, integration tests, or write/review test requests.
---

# Frontend Testing Skill

## Inmobil scope

Applies to **`inmobil-app/app/`** (Vite + React + MUI v7 + TanStack Router). Does **not** apply to:

- `inmobil-app/landing/` — Astro; uses its own testing approach (component tests with `@astro/test` if needed, otherwise rely on visual checks for marketing content).
- `inmobil-api/` — Go service, tested with `go test`.

All paths in this skill resolve relative to `inmobil-app/app/`. Run pnpm commands from inside that directory (or `pnpm -C inmobil-app/app …` from the repo root).

> Use Vitest mock/timer APIs (`vi.*`).

## When to Apply This Skill

Apply this skill when the user:

- Asks to **write tests** for a component, hook, or utility under `app/src/`
- Asks to **review existing tests** for completeness
- Mentions **Vitest**, **React Testing Library**, **RTL**, or **spec files**
- Requests **test coverage** improvement
- Mentions **testing**, **unit tests**, or **integration tests** for frontend code
- Wants to understand **testing patterns** in `inmobil-app/app/`

**Do NOT apply** when:

- User is asking about Go API tests in `inmobil-api/`
- User is asking about Astro component tests in `landing/`
- User is asking about E2E tests (Playwright)
- User is only asking conceptual questions without code context

## Quick Reference

### Key Commands

Run from `inmobil-app/app/` (or prefix with `pnpm -C inmobil-app/app`):

```bash
# Run all tests
pnpm test

# Watch mode
pnpm test --watch

# Run specific file
pnpm test path/to/file.spec.tsx

# Generate coverage report
pnpm test --coverage
```

### File Naming

- Test files: `ComponentName.spec.tsx` inside a same-level `__tests__/` directory
- Placement rule: Component, hook, and utility tests must live in a sibling `__tests__/` folder at the same level as the source under test. For example, `foo/index.tsx` maps to `foo/__tests__/index.spec.tsx`, and `foo/bar.ts` maps to `foo/__tests__/bar.spec.ts`.
- Integration tests that cross feature boundaries: `app/src/__tests__/` directory

## Test Structure Template

```typescript
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import Component from './index'

// ✅ Import real project components (DO NOT mock these)
// import { SuspenseLoader } from '~components/SuspenseLoader'
// import { ChildComponent } from './child-component'

// ✅ Mock external dependencies only
vi.mock('~features/properties/api/propertiesApi')

// TanStack Router: mock the hooks you actually use
vi.mock('@tanstack/react-router', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@tanstack/react-router')>()
  return {
    ...actual,
    useNavigate: () => vi.fn(),
    useParams: () => ({ id: '1' }),
  }
})

// Shared state for mocks (if needed)
let mockSharedState = false

describe('ComponentName', () => {
  beforeEach(() => {
    vi.clearAllMocks()  // ✅ Reset mocks BEFORE each test
    mockSharedState = false  // ✅ Reset shared state
  })

  // Rendering tests (REQUIRED)
  describe('Rendering', () => {
    it('should render without crashing', () => {
      // Arrange
      const props = { title: 'Test' }

      // Act
      render(<Component {...props} />)

      // Assert
      expect(screen.getByText('Test')).toBeInTheDocument()
    })
  })

  // Props tests (REQUIRED)
  describe('Props', () => {
    it('should apply custom className', () => {
      render(<Component className="custom" />)
      expect(screen.getByRole('button')).toHaveClass('custom')
    })
  })

  // User Interactions
  describe('User Interactions', () => {
    it('should handle click events', () => {
      const handleClick = vi.fn()
      render(<Component onClick={handleClick} />)

      fireEvent.click(screen.getByRole('button'))

      expect(handleClick).toHaveBeenCalledTimes(1)
    })
  })

  // Edge Cases (REQUIRED)
  describe('Edge Cases', () => {
    it('should handle null data', () => {
      render(<Component data={null} />)
      expect(screen.getByText(/no data/i)).toBeInTheDocument()
    })

    it('should handle empty array', () => {
      render(<Component items={[]} />)
      expect(screen.getByText(/empty/i)).toBeInTheDocument()
    })
  })
})
```

## Testing Workflow (CRITICAL)

### ⚠️ Incremental Approach Required

**NEVER generate all test files at once.** For complex components or multi-file directories:

1. **Analyze & Plan**: List all files, order by complexity (simple → complex)
1. **Process ONE at a time**: Write test → Run test → Fix if needed → Next
1. **Verify before proceeding**: Do NOT continue to next file until current passes

```
For each file:
  ┌────────────────────────────────────────┐
  │ 1. Write test                          │
  │ 2. Run: pnpm test <file>.spec.tsx      │
  │ 3. PASS? → Mark complete, next file    │
  │    FAIL? → Fix first, then continue    │
  └────────────────────────────────────────┘
```

### Complexity-Based Order

Process in this order for multi-file testing:

1. 🟢 Utility functions (simplest)
1. 🟢 Custom hooks
1. 🟡 Simple components (presentational)
1. 🟡 Medium components (state, effects)
1. 🔴 Complex components (API via `useSuspenseQuery`, routing)
1. 🔴 Integration tests (feature entry / index files - last)

### When to Refactor First

- **Complexity > 50**: Break into smaller pieces before testing
- **500+ lines**: Consider splitting before testing
- **Many dependencies**: Extract logic into hooks first

## Testing Strategy

### Path-Level Testing (Directory Testing)

When assigned to test a directory/path, test **ALL content** within that path:

- Test all components, hooks, utilities in the directory (not just `index` file)
- Use incremental approach: one file at a time, verify each before proceeding
- Goal: 100% coverage of ALL files in the directory

### Integration Testing First

**Prefer integration testing** when writing tests for a feature directory under `app/src/features/<feature>/`:

- ✅ **Import real project components** directly (including primitives in `~components/` and feature siblings)
- ✅ **Only mock**: API services under `~features/<feature>/api/*`, TanStack Router hooks when navigation/params drive the test, complex context providers
- ❌ **DO NOT mock** primitives under `~components/` (e.g. `SuspenseLoader`, `CustomAppBar`) or MUI components
- ❌ **DO NOT mock** sibling/child components in the same directory

> See [Test Structure Template](#test-structure-template) for correct import/mock patterns.

### Suspense + React Query

Components in `app/` follow the Suspense-first doctrine from `frontend-dev-guidelines`. Tests for them should:

- Render inside a `QueryClientProvider` test helper (place in `app/src/test/queryClient.tsx`).
- Wrap the component under test in `<Suspense fallback={…}>` if the component itself doesn't.
- Assert against the resolved UI; use `await screen.findBy…` rather than asserting on loading fallbacks unless the loading boundary itself is what you're testing.

## Core Principles

### 1. AAA Pattern (Arrange-Act-Assert)

Every test should clearly separate:

- **Arrange**: Setup test data and render component
- **Act**: Perform user actions
- **Assert**: Verify expected outcomes

### 2. Black-Box Testing

- Test observable behavior, not implementation details
- Use semantic queries (`getByRole` with accessible `name`, `getByLabelText`, `getByPlaceholderText`, `getByText`, and scoped `within(...)`)
- Treat `getByTestId` as a last resort. If a control cannot be found by role/name, label, landmark, or dialog scope, fix the component accessibility first instead of adding or relying on `data-testid`.
- Remove production `data-testid` attributes when semantic selectors can cover the behavior. Keep them only for non-visual mocked boundaries, canvas/chart output, or third-party widgets with no accessible DOM in the test environment.
- Do not assert decorative icons by test id. Assert the named control that contains them, or mark decorative icons `aria-hidden`.
- Avoid testing internal state directly
- **Prefer pattern matching over hardcoded strings** in assertions:

```typescript
// ❌ Avoid: hardcoded text assertions
expect(screen.getByText('Loading...')).toBeInTheDocument()

// ✅ Better: role-based queries
expect(screen.getByRole('status')).toBeInTheDocument()

// ✅ Better: pattern matching
expect(screen.getByText(/loading/i)).toBeInTheDocument()
```

### 3. Single Behavior Per Test

Each test verifies ONE user-observable behavior:

```typescript
// ✅ Good: One behavior
it('should disable button when loading', () => {
  render(<Button loading />)
  expect(screen.getByRole('button')).toBeDisabled()
})

// ❌ Bad: Multiple behaviors
it('should handle loading state', () => {
  render(<Button loading />)
  expect(screen.getByRole('button')).toBeDisabled()
  expect(screen.getByText('Loading...')).toBeInTheDocument()
  expect(screen.getByRole('button')).toHaveClass('loading')
})
```

### 4. Semantic Naming

Use `should <behavior> when <condition>`:

```typescript
it('should show error message when validation fails')
it('should call onSubmit when form is valid')
it('should disable input when isReadOnly is true')
```

## Required Test Scenarios

### Always Required (All Components)

1. **Rendering**: Component renders without crashing
1. **Props**: Required props, optional props, default values
1. **Edge Cases**: null, undefined, empty values, boundary conditions

### Conditional (When Present)

| Feature | Test Focus |
|---------|-----------|
| `useState` | Initial state, transitions, cleanup |
| `useEffect` | Execution, dependencies, cleanup |
| Event handlers | All onClick, onChange, onSubmit, keyboard |
| `useSuspenseQuery` | Resolved data render, error boundary behavior |
| TanStack Router | Navigation, params, search params |
| `useCallback`/`useMemo` | Referential equality |
| Context | Provider values, consumer behavior |
| Forms | Validation, submission, error display |

## Coverage Goals (Per File)

For each test file generated, aim for:

- ✅ **100%** function coverage
- ✅ **100%** statement coverage
- ✅ **>95%** branch coverage
- ✅ **>95%** line coverage

## Project Configuration

These files live under `inmobil-app/app/` once the sub-app is scaffolded:

- `vite.config.ts` — Vite/Vitest configuration
- `src/test/setup.ts` — Test environment setup (jest-dom matchers, global mocks)
- `src/test/queryClient.tsx` — React Query test wrapper

Modules are not mocked automatically. Global mocks (e.g. browser shims) belong in `src/test/setup.ts`; mock feature-specific modules locally in test files.
