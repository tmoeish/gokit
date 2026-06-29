# CLAUDE.md

This project's working rules for coding agents live in **[AGENTS.md](AGENTS.md)** —
read it before making changes. Quick reminders:

- `gokit` is a flat, generics-first Go **utility library** (Go's Guava/Lodash).
  No app, no `main`. Every exported symbol is public API; add functions, don't
  break signatures.
- New helpers go in the existing `<domain>x` package that owns them; only create
  a new package for a new domain (and update README, `doc.go`, **and both skill
  catalogs**: `.claude/skills/gokit/` and `skills/gokit-usage/`).
- Run **`make check`** (gofmt + vet + lint + `go test -race -cover`) before
  considering any change done.
- Tests are black-box (`package <pkg>_test`), standard `testing` only, covering
  edge cases.
- Public OSS repo: commit as `tmoeish <tmoeish@users.noreply.github.com>` (no
  personal name/email, no `.claude/settings.local.json`); Conventional-Commits
  messages. **Don't bump the version or push/tag unless asked to release** — and
  when releasing, sync `.claude-plugin/*` `version` fields.

Full conventions, the harness, and pitfalls: see **[AGENTS.md](AGENTS.md)**.
