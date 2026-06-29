// Package gokit is the module root for github.com/tmoeish/gokit, a
// batteries-included companion to the Go standard library — the Swiss Army
// knife that Java calls Guava and JavaScript calls Lodash.
//
// gokit is a collection of small, focused, generics-first packages. Import only
// the ones you need:
//
//   - slicex  — generic slice utilities (map/filter/reduce, set ops, chunking)
//   - mapx    — generic map utilities (keys/values, merge, pick/omit, invert)
//   - setx    — a generic hash Set with union/intersection/difference
//   - strx    — string helpers (case conversion, slugify, blank checks, pad)
//   - mathx   — generic math (clamp, gcd/lcm, rounding, primes, safe divide)
//   - conv    — safe any → string/number/bool conversions with explicit errors
//   - timex   — time helpers (begin/end of period, work-days, age, parsing)
//   - randx   — random strings, ints, UUIDs, weighted choice, sampling
//   - cryptox — hashing (MD5/SHA/HMAC) and base64/hex encoding
//   - iox     — file and I/O helpers (read/write lines, exists, copy, temp)
//   - jsonx   — JSON and YAML marshal/unmarshal, pretty-print, deep clone
//   - netx    — network helpers (local IPs, IP/CIDR validation, free ports)
//   - ptrx    — pointer helpers (Of/Deref, nil-coalescing)
//   - must    — Must/OK/Assert panic-on-error helpers
//   - retry   — context-aware retries with configurable backoff
//   - syncx   — concurrent structures (SafeMap, generic Once, WaitGroupCtx)
//   - contextx — type-safe context values and request metadata helpers
//   - logx    — structured logging on log/slog with context propagation
//   - httpx   — framework-agnostic HTTP response envelope and error codes
//   - echox   — bridges httpx + contextx into the labstack/echo framework
//
// This root package intentionally exposes no symbols; see the subpackages.
package gokit
