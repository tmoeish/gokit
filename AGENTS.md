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
<pkg>/<pkg>.go             implementation       (package <pkg>)
<pkg>/<pkg>_test.go         black-box tests      (package <pkg>_test)
doc.go                      module landing doc   (package gokit, no symbols)
Makefile                    dev harness — `make check`
.golangci.yml               lint config (golangci-lint v2 schema)
.github/workflows/ci.yml    CI (build/test on Go 1.22/1.24/1.26 + lint)
LICENSE                     MIT
README.md                   user-facing overview + package table + badges
.claude/skills/gokit/       contributor skill (gokit-dev) — THIS guidance
skills/gokit-usage/         consumer skill — for agents USING gokit elsewhere
.claude-plugin/             plugin + marketplace manifests (ship the consumer skill)
```

Current packages (20): `slicex mapx setx strx mathx conv timex randx cryptox iox
jsonx netx ptrx must retry syncx contextx logx httpx echox`. See README.md for
what each does.

There are **two skills**, keep them distinct: `.claude/skills/gokit/` is the
contributor guide (this file in skill form); `skills/gokit-usage/` teaches agents
to *use* gokit in other projects and is shipped via the `.claude-plugin/`
marketplace. When you add/rename a package, update **both** catalogs (see below).

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
domain. Only create a new package for a genuinely new domain, and when you do,
update **all four** discovery surfaces in the same change (CI/review will assume
they're in sync):

1. `<pkg>/<pkg>.go` + `<pkg>/<pkg>_test.go` with a package doc comment that
   credits any library it borrows from.
2. The package table in `README.md`.
3. The package list in `doc.go`.
4. **Both** skill catalogs: `.claude/skills/gokit/SKILL.md` (contributor) and
   `skills/gokit-usage/SKILL.md` (consumer "I'm about to write X → use Y" table).

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

## Git, commits & privacy

This is a **public open-source repo**. Keep private/historical info out of it —
forever, because git history is permanent.

- **Commit identity is the public handle, not a personal one.** This repo's
  commits are authored as `tmoeish <tmoeish@users.noreply.github.com>` (set
  repo-locally via `git config user.name/user.email`). Do **not** commit under a
  personal real name or private email. Don't change global git config; set it on
  the repo.
- **No private/historical traces in code or comments**: no internal project
  names, ticket IDs, internal hostnames/URLs, employer names, or
  "adapted from <internal repo>" notes. Borrowed-from credit should name only
  public OSS (samber/lo, lancet, Guava…). A quick sweep before committing:
  `grep -rniE "<personal-name>|@<corp>|internal-host|jira|gitlab|gitee" .`
- **`.gitignore` excludes local tooling state**: `.claude/settings*.json` are
  personal and must stay untracked. The shipped skills (`.claude/skills/`,
  `skills/`) and plugin manifests (`.claude-plugin/`) **are** tracked.
- **Commit messages**: Conventional-Commits style
  (`feat: …`, `fix: …`, `docs: …`, `feat(scope): …`), imperative mood, and end
  with the agent co-author trailer:
  `Co-Authored-By: Claude <noreply@anthropic.com>`.
- **Push only when asked.** Pushing and tagging are outward-facing — don't do
  them unprompted.

## Releasing & versioning

**Do not auto-bump the version on every commit.** Most changes land on `main`
and are captured by the *next* deliberate tag. Cut a release only when there's a
meaningful, releasable delta and the user asks for it.

- **SemVer, pre-1.0 conventions** (we are in `v0.x`): bump **patch**
  (`v0.0.x`) for fixes, docs, tooling, and skill/plugin updates; bump **minor**
  (`v0.x.0`) for new packages or new public functions. Breaking changes are
  avoided entirely until a deliberate `v1.0.0` (and a `/v2+` module-path suffix
  thereafter — don't add one now).
- **Keep plugin versions in sync.** When you tag a release, bump the `version`
  field in `.claude-plugin/plugin.json` **and** `.claude-plugin/marketplace.json`
  to match, in the same change.
- **Release steps** (once the user authorizes a release):
  1. `make check` clean; working tree clean on `main`.
  2. `git tag -a vX.Y.Z -m "…"` (annotated, with a short changelog) and
     `git push origin vX.Y.Z`. For Go modules the **pushed tag is the release** —
     `go get github.com/tmoeish/gokit@vX.Y.Z` works off it.
  3. If `gh` is authenticated, create the GitHub Release page:
     `gh release create vX.Y.Z --title vX.Y.Z --notes "…"` (add `--latest` for
     the newest). If `gh` isn't authed, the tag alone still publishes the version.
  4. Verify: `go list -m github.com/tmoeish/gokit@vX.Y.Z` resolves via the proxy.

## Go open-source best practices (already wired — keep them true)

- **godoc**: every exported symbol documented; `doc.go` is the module landing
  page on pkg.go.dev. Don't break the doc-comment convention (`revive` guards it).
- **Badges in README** (pkg.go.dev, CI, Go Report Card, Go version, release,
  license) must keep pointing at real endpoints.
- **golangci-lint is v2** (`.golangci.yml` starts with `version: "2"`). In v2,
  `gofmt`/`goimports` are *formatters*, not linters — keep them under the
  `formatters:` block. Don't downgrade the config to the v1 schema.
- **CI** runs the build/test matrix + lint on every push/PR; keep the Go matrix
  current (oldest supported + latest two).
- **LICENSE** (MIT) and module path (`github.com/tmoeish/gokit`) are fixed —
  don't rename the module.

## Definition of done for a change

1. Code + doc comments follow the conventions above.
2. New/changed behavior has tests; `make test` passes with `-race`.
3. `make check` is clean (fmt, vet, lint, test).
4. For a new/renamed package: README table, `doc.go`, **and both skill catalogs**
   updated (see "Where new code goes").
5. No breaking changes to existing exported signatures (add, don't mutate).
6. Commit follows the Git/commits/privacy rules; no version bump or push unless
   the user asked to release.

## Common pitfalls (don't)

- Don't add a `cmd/`, `internal/`, `pkg/`, or `src/` layer — packages stay flat
  at the module root (this is the idiomatic layout for a library; see README).
- Don't introduce mutation-by-default helpers or hidden globals.
- Don't pull a heavy dependency into a leaf utility package.
- Don't leave trailing blank lines / unformatted files — `gofmt` and CI will
  reject them. (`make fmt` fixes them.)
- Don't change an exported function's signature to "fix" a caller; add a new
  function.
- Don't write the `.golangci.yml` in the old v1 schema — it's **v2**
  (`version: "2"`, formatters under `formatters:`). v1 configs fail to load.
- Don't commit `.claude/settings.local.json` or author commits under a personal
  name/email — see "Git, commits & privacy".
- Don't tag/push a release unprompted, and don't forget to bump the
  `.claude-plugin/*` `version` fields when you do tag one.
