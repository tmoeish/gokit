# AGENTS.md — working rules for gokit

Rules for humans and coding agents working in this repository. Read this before
writing code. (Claude Code also loads `CLAUDE.md`, which points here.)

## What this project is

`gokit` (`github.com/tmoeish/gokit`) is a **general-purpose Go utility library** —
Go's answer to Java's Guava / JavaScript's Lodash. It is a flat collection of
small, independent, generics-first packages, each complementing one part of the
standard library. There is **no application, no `main`, no server** here — only
reusable library code.

Because it's a library, the bar is different from app code:

- Every exported symbol is public API. Renames and signature changes are
  **breaking** — avoid them; add new functions instead of changing old ones.
- Correctness, clear doc comments, and tests matter more than cleverness.

## Repository map

```
<pkg>/<pkg>.go        implementation        (package <pkg>)
<pkg>/<pkg>_test.go    black-box tests        (package <pkg>_test)
doc.go                 module landing doc     (package gokit, no symbols)
Makefile               dev harness — `make check`
.golangci.yml          lint config
.github/workflows/ci.yml  CI (build/test on Go 1.22–1.24 + lint)
README.md              user-facing overview + package table
.claude/skills/gokit/  this guidance, as a Claude Code skill
```

Current packages (20): `slicex mapx setx strx mathx conv timex randx cryptox iox
jsonx netx ptrx must retry syncx contextx logx httpx echox`. See README.md for
what each does.

## The build/test harness — always run it

```bash
make check      # fmt-check + vet + lint + test(-race -cover) — the gate
make test       # go test -race -cover ./...
make fmt        # gofmt -w (run this if fmt-check fails)
```

A change is **not done** until `make check` is clean. If `golangci-lint` isn't
installed, `make lint` skips with a notice — that's fine locally, but CI enforces
it, so keep code lint-clean.

## Conventions (match the existing code)

**Package naming.** One domain per package, named `<domain>x` to avoid colliding
with the std-lib package it extends (`slicex` ↔ `slices`, `strx` ↔ `strings`).
Lowercase, no underscores. The package doc comment is one or two sentences and
**credits the libraries it borrows from** (see any existing file, e.g.
`slicex/slice.go`).

**Where new code goes.** Add helpers to the **existing** package that owns the
domain. Only create a new package for a genuinely new domain, and when you do:
add `<pkg>/<pkg>.go` + `<pkg>/<pkg>_test.go`, a package doc comment, and a row in
the README table + `doc.go` list.

**Function style.**
- Generics-first: use type parameters (`[T any]`, `[K comparable]`,
  `[T cmp.Ordered]`, the local `Number` constraint) instead of `any`/reflection.
- Pure: return new slices/maps; never mutate inputs unless the name says so
  (`Fill`, `Add`, `Remove`). Document ordering guarantees or their absence.
- Fallible work returns `error`. Provide a `Must*` variant (or use the `must`
  package) only when panicking is a reasonable caller choice. "Not found" lookups
  return `(T, bool)` — see `slicex.First`, `mapx.GetOrDefault`.
- Don't reimplement what `slices`, `maps`, `cmp`, `strings`, `strconv` already do
  well. Extend or compose; don't duplicate.

**Doc comments.** Every exported func/type/const starts with its name
(`// Map applies fn …`). This is enforced by `revive`.

**Dependencies.** Keep core packages dependency-light. The whole module currently
depends only on `google/uuid`, `yaml.v3`, and `labstack/echo/v4`. A
framework/third-party dependency must be **isolated in its own package** so
importers of other packages don't pay for it — `echox` isolates `echo` away from
the framework-agnostic `httpx`. Run `make tidy` after touching imports.

## Tests

- Black-box: `package <pkg>_test`, importing the package under its real path.
- Plain standard-library `testing` — **no assertion frameworks**. Use
  `t.Fatalf("Name: got %v, want %v", got, want)`.
- Cover the happy path **and** edge cases: empty input, nil, single element,
  boundary values. Table-driven tests are welcome where they read well.

## Definition of done for a change

1. Code + doc comments follow the conventions above.
2. New/changed behavior has tests; `make test` passes with `-race`.
3. `make check` is clean (fmt, vet, lint, test).
4. README package table and `doc.go` updated if you added/renamed a package.
5. No breaking changes to existing exported signatures (add, don't mutate).

## Common pitfalls (don't)

- Don't add a `cmd/`, `internal/`, `pkg/`, or `src/` layer — packages stay flat
  at the module root (this is the idiomatic layout for a library; see README).
- Don't introduce mutation-by-default helpers or hidden globals.
- Don't pull a heavy dependency into a leaf utility package.
- Don't leave trailing blank lines / unformatted files — `gofmt` and CI will
  reject them.
- Don't change an exported function's signature to "fix" a caller; add a new
  function.
