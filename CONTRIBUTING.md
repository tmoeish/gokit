# Contributing to gokit

Thanks for helping build Go's Swiss Army knife. This guide covers the workflow;
the **conventions** every change must follow live in [AGENTS.md](AGENTS.md).

## Workflow

1. **Pick the right package.** New helpers belong in the existing `xxx` package
   that owns their domain (string → `strx`, slice → `slicex`, …). Only create a
   new package for a genuinely new domain — see AGENTS.md for the bar.
2. **Write the helper + doc comment.** Every exported symbol needs a doc comment
   starting with its name. Match the style of the surrounding code.
3. **Write a test.** Tests live in the same directory as `package <pkg>_test`
   (black-box). Cover the happy path and at least one edge case (empty input,
   nil, boundary).
4. **Run `make check`** and make sure it is clean.
5. **Open a PR** describing what you added and which peer library / Guava feature
   it mirrors, if any.

## Local commands

```bash
make fmt     # gofmt -w the tree
make vet     # go vet ./...
make lint    # golangci-lint (if installed)
make test    # go test -race -cover ./...
make check   # fmt-check + vet + lint + test — run this before every PR
```

You don't need `golangci-lint` installed to contribute — `make lint` skips
gracefully if it's missing — but CI runs it, so installing it locally is
recommended:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## What makes a good gokit addition

- **General-purpose.** If only one project would ever use it, it doesn't belong
  here.
- **Generics-first.** Prefer type parameters over `any`/reflection.
- **Pure and predictable.** Return new values; don't mutate inputs unless the
  name says so. Document ordering guarantees (or the lack of them).
- **Std-lib aware.** Don't reimplement something `slices`, `maps`, `cmp`, or
  `strings` already does well. Wrap or extend, don't duplicate.
- **Dependency-light.** Adding a third-party dependency to a core package needs a
  strong justification. Framework-specific code goes in its own package (see how
  `echox` isolates the `echo` dependency from `httpx`).

## Reporting bugs

Open an issue with a minimal reproduction (ideally a failing test). PRs with a
fix and a regression test are even better.
