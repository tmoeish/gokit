# gokit

[![Go Reference](https://pkg.go.dev/badge/github.com/tmoeish/gokit.svg)](https://pkg.go.dev/github.com/tmoeish/gokit)
[![CI](https://github.com/tmoeish/gokit/actions/workflows/ci.yml/badge.svg)](https://github.com/tmoeish/gokit/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tmoeish/gokit)](https://goreportcard.com/report/github.com/tmoeish/gokit)
[![Go Version](https://img.shields.io/github/go-mod/go-version/tmoeish/gokit)](go.mod)
[![Release](https://img.shields.io/github/v/tag/tmoeish/gokit?label=release&sort=semver)](https://github.com/tmoeish/gokit/releases)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

> A batteries-included standard-library companion for Go — the Swiss Army knife
> that Java developers know as [Guava](https://github.com/google/guava) and
> JavaScript developers know as Lodash.

`gokit` is a collection of small, focused, **generics-first** packages that fill
the gaps in Go's standard library: slices, maps, strings, math, time, crypto,
I/O, networking, concurrency, retries, structured logging, and HTTP plumbing.
Every package is independent — import only what you need, and pull in only the
dependencies that package requires.

```go
import "github.com/tmoeish/gokit/slicex"

evens := slicex.Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
// []int{2, 4}
```

## Install

```bash
go get github.com/tmoeish/gokit@latest
```

Requires **Go 1.22+** (uses generics and `math/rand/v2`).

## Design principles

- **One concern per package.** Each `xxx` package owns one domain and uses the
  conventional `x` suffix (`slicex`, `mapx`, `strx`, …) to avoid clashing with
  the std-lib package it complements.
- **Generics over `interface{}`.** Functions are type-safe; no reflection on hot
  paths.
- **Pure functions, no surprises.** Helpers return new values instead of mutating
  inputs unless the name says otherwise (e.g. `Fill`).
- **Errors are values, panics are opt-in.** Fallible helpers return `error`; the
  `must` package and `Must*` variants exist for the cases where panicking is what
  you actually want.
- **Minimal dependencies.** The only third-party deps are `google/uuid`,
  `yaml.v3`, and `labstack/echo/v4` — and the last is only pulled in if you
  import `echox`.

## Packages

| Package    | What it gives you                                                                 | Inspired by |
|------------|-----------------------------------------------------------------------------------|-------------|
| `slicex`   | Map/Filter/Reduce, set ops, chunk/window/zip, sampling, min/max-by               | samber/lo, lancet, Guava |
| `mapx`     | Keys/Values/Entries, Merge, Pick/Omit, Invert, MapValues, GroupBy                | samber/lo, lancet |
| `setx`     | Generic hash `Set[T]`: union/intersection/difference, subset/superset            | Guava `Sets`, lancet |
| `strx`     | Blank checks, case conversion (camel/snake/kebab/pascal), slugify, pad, between  | Guava `Strings`, lancet |
| `mathx`    | Clamp, Abs, Sum/Average, GCD/LCM, Pow, prime/fib, rounding, safe divide          | lancet, Guava `math` |
| `conv`     | Safe `any → string/int/float/bool` conversions with explicit errors              | spf13/cast |
| `timex`    | Begin/End of day/week/month/year, work-days, age, durations, parsing             | jinzhu/now, lancet |
| `randx`    | Random strings/charsets, ints, UUIDs, weighted choice, sampling                  | lancet, Guava |
| `cryptox`  | MD5/SHA/HMAC hashing, base64/hex (encode + verify)                               | std crypto wrappers |
| `iox`      | File read/write/append, lines, exists, copy, temp files, ensure-dir              | lancet `fileutil` |
| `jsonx`    | JSON + YAML marshal/unmarshal, pretty-print, deep clone, path `Get`              | std encoding wrappers |
| `netx`     | Local IPs, IP/CIDR validation, free ports, host:port parsing                     | lancet `netutil` |
| `ptrx`     | `Of`/`Deref` pointer helpers, nil-coalescing                                     | AWS `aws.String`, samber/lo |
| `must`     | `Must`/`OK`/`Assert` panic-on-error helpers                                       | samber/lo `Must` |
| `retry`    | Context-aware retries with constant/linear/exponential backoff + jitter          | avast/retry-go, cenkalti/backoff |
| `syncx`    | `SafeMap`, generic `Once[T]`, `WaitGroupCtx`, `Notifier`                          | std `sync`, samber/lo |
| `contextx` | Type-safe context values; package-defined request metadata helpers               | std `context` |
| `logx`     | Structured logging on `log/slog` with context propagation + caller/stack helpers | std `log/slog` |
| `httpx`    | Framework-agnostic response envelope, application error codes, pagination        | — |
| `echox`    | Wires `httpx` + `contextx` into `labstack/echo` (responses, error handler, mw)   | — |

Full, runnable API docs live on
[pkg.go.dev/github.com/tmoeish/gokit](https://pkg.go.dev/github.com/tmoeish/gokit).

## Examples

### slicex — transform and aggregate

```go
people := []Person{{"Ann", 30}, {"Bob", 25}, {"Cy", 30}}

byAge := slicex.GroupBy(people, func(p Person) int { return p.Age })
names := slicex.Map(people, func(p Person) string { return p.Name })
oldest, _ := slicex.MaxBy(people, func(p Person) int { return p.Age })
```

### setx — set algebra

```go
a := setx.New(1, 2, 3)
b := setx.New(2, 3, 4)
a.Intersection(b).ToSlice() // [2 3]
a.Union(b).Len()            // 4
a.IsSubsetOf(b)             // false
```

### retry — resilient calls

```go
err := retry.Do(ctx, func() error {
    return callFlakyAPI()
},
    retry.WithMaxAttempts(5),
    retry.WithExponentialBackoff(100*time.Millisecond, 2.0),
    retry.WithJitter(0.2),
)
```

### echox — HTTP responses for echo

```go
e := echo.New()
e.HTTPErrorHandler = echox.ErrorHandler
e.Use(echox.RequestID())

e.GET("/users/:id", func(c echo.Context) error {
    u, err := store.Find(c.Param("id"))
    if err != nil {
        return echox.Error(c, httpx.CodeRecordNotFound, err)
    }
    return echox.Success(c, u)
})
```

## Project layout

The repository follows
[Effective Go](https://go.dev/doc/effective_go) and the standard Go module
conventions: a flat set of single-purpose packages at the module root, each in
its own directory with co-located `*_test.go` files. There is no `pkg/` or
`src/` nesting — import paths read as `github.com/tmoeish/gokit/<pkg>`.

```
gokit/
├── slicex/      mapx/      setx/     strx/      mathx/
├── conv/        timex/     randx/    cryptox/   iox/
├── jsonx/       netx/      ptrx/     must/      retry/
├── syncx/       contextx/  logx/     httpx/     echox/
├── go.mod
├── Makefile         # dev harness: make check / test / lint / fmt
├── AGENTS.md        # rules for humans and coding agents
└── README.md
```

## Using gokit with AI coding agents

gokit ships a **reuse-first skill** (`gokit-usage`) that teaches coding agents
(Claude Code, etc.) to reach for gokit instead of hand-rolling slice/map/string/
time/math/crypto/retry helpers in *your* project. Install it either way:

**As a Claude Code plugin** (recommended — auto-updates with the repo):

```text
/plugin marketplace add tmoeish/gokit
/plugin install gokit
```

**Or copy the skill** into your project or personal skills directory:

```bash
# project-local
mkdir -p .claude/skills/gokit-usage
curl -sL https://raw.githubusercontent.com/tmoeish/gokit/main/skills/gokit-usage/SKILL.md \
  -o .claude/skills/gokit-usage/SKILL.md
```

The skill catalogs every package and maps "I'm about to write X" → "use gokit's Y."

> Note: the `.claude/skills/gokit/` directory in *this* repo is a separate
> **contributor** skill (`gokit-dev`) for developing gokit itself — not needed to
> consume the library.

## Contributing & developing

Run the full local check (format, vet, lint, test with race + coverage) before
sending changes:

```bash
make check
```

See **[AGENTS.md](AGENTS.md)** for the conventions every change must follow and
**[CONTRIBUTING.md](CONTRIBUTING.md)** for the workflow. The same rules are
encoded as a Claude Code skill under `.claude/skills/gokit/` so coding agents
pick them up automatically.

## License

MIT
