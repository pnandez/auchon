# 0003: godog for BDD acceptance tests

Status: accepted (2026-11-21)

## Context

Features are described as Gherkin `.feature` files (AI-skills `bdd-acceptance` module). The
first choice, gherkingen, only generates a `_test.go` skeleton once: after that, editing a
scenario does not change or break its test, so the two drift apart silently.

Libraries that run the feature file itself were compared on GitHub (2026-11-21):

| | godog | gobdd | gherkingen | gocuke |
|---|---|---|---|---|
| Stars | 2,693 | 147 | 99 | 35 |
| Contributors | 84 | 16 | 5 | 7 |
| Latest release | v0.16.0 (2026-07) | v1.1.4 (2025-07) | v4.1.1 (2026-07) | v1.1.1 (2024-03) |
| Commits, last 12 months | 93 (20 by people) | 4 | 3 | 0 |

## Decision

Use [godog](https://github.com/cucumber/godog), Cucumber's own Go implementation, as described
in the AI-skills `bdd-acceptance` module:

- `TestFeatures` in each context's `acceptance/` runs every `.feature` in that folder through
  `go test`, one subtest per scenario
- Strict mode: a step without a definition fails the run
- Steps are methods on a per-scenario state struct, registered with anchored regular
  expressions

## Consequences

- The feature file is the test input: changing a step's wording makes it undefined and fails
  the suite, so feature files and tests cannot drift
- Editors with Cucumber support can jump from a step to its Go function
- Steps are bound by regular expressions, which is more typing than gocuke's method-name
  matching
