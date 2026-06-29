---
name: gokit-dev
description: "Develop/contribute to the gokit library itself (this repo, github.com/tmoeish/gokit). Use when adding, fixing, or reviewing helpers in gokit's own xxx packages (slicex, mapx, strx, setx, etc.). Covers the package catalog, where new code goes, the conventions every change must follow, and the make-check harness. NOT for merely using gokit in another project — that's the gokit-usage skill."
---

# gokit (contributor guide)

`github.com/tmoeish/gokit` is a general-purpose, generics-first Go utility
library — Go's answer to Guava/Lodash. A flat set of small, independent `xxx`
packages at the module root, each complementing one part of the standard
library. **It is a library: no app, no `main`.** Every exported symbol is public
API — add functions, never break existing signatures.

The authoritative rules are in the repo's `AGENTS.md`. This skill is the quick
operational map.

## First, reuse before writing

Before writing a utility, check whether gokit already has it. Package → domain:

| Need… | Package | Examples |
|------|---------|----------|
| slice transforms / set ops | `slicex` | `Map Filter Reduce GroupBy Chunk Unique Intersection MinBy` |
| map helpers | `mapx` | `Keys Values Merge Pick Omit Invert MapValues GroupBy` |
| a Set collection | `setx` | `New Add Contains Union Intersection Difference IsSubsetOf` |
| string helpers | `strx` | `IsBlank ToSnakeCase ToCamelCase Slugify Truncate Between Pad*` |
| math | `mathx` | `Clamp Abs Sum Average GCD LCM Pow RoundTo SafeDivide` |
| any→T conversion | `conv` | `ToString ToInt64 ToFloat64 ToBool` (+ `Must*`) |
| time | `timex` | `BeginOfDay EndOfMonth DaysBetween AgeAt AddWorkDays ParseDate` |
| randomness | `randx` | `String Int UUID WeightedChoice Shuffle Sample` |
| hashing/encoding | `cryptox` | `SHA256 HMACSHA256 Base64Encode HexEncode` |
| files/IO | `iox` | `ReadFile WriteLines FileExists CopyFile TempDir EnsureDir` |
| JSON/YAML | `jsonx` | `Marshal Unmarshal MarshalPretty Clone JSONToYAML Get` |
| network | `netx` | `LocalIP IsValidIP IsIPInCIDR GetFreePort ParseHostPort` |
| pointers | `ptrx` | `Of Deref IsNil CoalescePtr` |
| panic-on-error | `must` | `Must OK Assert` |
| retries/backoff | `retry` | `Do DoWithResult With{MaxAttempts,ExponentialBackoff,Jitter}` |
| concurrency | `syncx` | `SafeMap Once[T] WaitGroupCtx Notifier` |
| context values | `contextx` | `WithValue Value WithReqID ReqID WithLogger Logger` |
| structured logging | `logx` | `NewLogger Info/Error(ctx,…) LogError CallerLoc` |
| HTTP envelope (no framework) | `httpx` | `Code NewSuccessResponse NewErrorResponse NewPageResponse` |
| echo + httpx glue | `echox` | `Success Error Page RequestID() ErrorHandler` |

Run `go doc github.com/tmoeish/gokit/<pkg>` for the full signature list.

## Where new code goes

1. **Default: extend the existing package** that owns the domain. A new string
   helper goes in `strx`, a new slice helper in `slicex` — not a new package.
2. **New package only for a new domain.** Then create `<pkg>/<pkg>.go` +
   `<pkg>/<pkg>_test.go`, give the package a doc comment that credits any library
   it borrows from, and update **all four** discovery surfaces together:
   `README.md` table, `doc.go` list, `.claude/skills/gokit/SKILL.md` (this file),
   and `skills/gokit-usage/SKILL.md` (the consumer catalog).
3. **Framework/heavy deps go in their own package**, isolated like `echox`
   isolates `labstack/echo` from the framework-agnostic `httpx`. Never pull such
   a dep into a leaf utility package.

## Conventions (match existing code exactly)

- **Generics-first.** Type parameters (`[T any]`, `[K comparable]`,
  `[T cmp.Ordered]`, local `Number`) over `any`/reflection.
- **Pure.** Return new slices/maps; mutate inputs only when the name says so
  (`Fill`, `Add`, `Remove`). Document ordering (or "unspecified order").
- **Errors are values.** Fallible funcs return `error`; add a `Must*` variant
  only when panicking is a sane caller choice. Lookups return `(T, bool)`.
- **Doc comments** on every exported symbol, starting with its name
  (`// Map applies …`) — enforced by `revive`.
- **Don't duplicate** `slices`/`maps`/`cmp`/`strings`/`strconv`; extend/compose.

## Tests

Black-box: `package <pkg>_test`, standard library `testing` only (no assertion
frameworks), `t.Fatalf("Name: got %v, want %v", got, want)`. Cover happy path +
edge cases (empty, nil, single element, boundaries). See `mapx/map_test.go` or
`setx/set_test.go` as templates.

## Harness — run before declaring done

```bash
make check     # gofmt-check + go vet + golangci-lint + go test -race -cover
make fmt       # gofmt -w (fix formatting)
make test      # tests only
```

A change isn't done until `make check` is clean. CI (`.github/workflows/ci.yml`)
runs the same on Go 1.22/1.24/1.26 plus golangci-lint. Config `.golangci.yml` is
the **v2 schema** (`version: "2"`; `gofmt`/`goimports` live under `formatters:`)
— don't rewrite it as v1.

## Commit, privacy & release (public OSS repo)

- **Commit identity = the public handle**, set repo-locally:
  `tmoeish <tmoeish@users.noreply.github.com>`. Never author commits under a
  personal real name/email; never commit `.claude/settings.local.json`.
- **No private/historical info** in code or comments (internal repo names, ticket
  IDs, hostnames, employer). Credit only public OSS in "inspired by" notes.
- **Commit messages**: Conventional Commits (`feat:`/`fix:`/`docs:`), imperative,
  ending with `Co-Authored-By: Claude <noreply@anthropic.com>`.
- **Versioning is deliberate, not per-commit.** Don't auto-bump. Cut a tag only
  when asked. SemVer pre-1.0: patch = fixes/docs/tooling, minor = new
  package/function; avoid breaking changes pre-`v1.0.0`.
- **On release**: tag annotated `vX.Y.Z` + push (the tag *is* the Go release);
  `gh release create` if `gh` is authed; **bump `version` in both
  `.claude-plugin/plugin.json` and `marketplace.json`** to match; verify with
  `go list -m github.com/tmoeish/gokit@vX.Y.Z`.
- **Don't push or tag unprompted** — these are outward-facing.

## Definition of done

1. Conventions + doc comments followed.
2. Tests added; `make test` passes with `-race`.
3. `make check` clean.
4. New/renamed package → README table, `doc.go`, **and both skill catalogs**
   updated.
5. No breaking changes to exported signatures (add, don't mutate).
6. Commit follows the privacy/message rules; no version bump or push unless the
   user asked to release.
