# 2026-11-13 · log/slog structured logging

## Why learn this
- **Logs are read by machines first**: Loki, Elasticsearch, CloudWatch and Datadog can filter `status=500 user=perry` only if the log line has fields, not a sentence built with `fmt.Sprintf`.
- **It is in the standard library** since Go 1.21, so libraries can accept a `*slog.Logger` (or `slog.Handler`) without forcing zap/zerolog on their users. The `logger` topic in this repo uses zap; `slog` is the stdlib answer to the same problem, and zap/zerolog can act as slog backends.
- **Request-scoped context**: attach `request_id`, `trace_id` (see the `otel` topic) and `user` once, and every later log line carries them.
- It builds on earlier lessons: a handler writes to an `io.Writer` (`io_reader_writer`), values come from `context` (`context`), and `LevelVar` is an atomic (`sync_atomic`).

## Core concepts
1. **Logger -> Handler**: `*slog.Logger` is the front end you call (`Info`, `Warn`, `Error`, `Debug`, `...Context`). It builds a `slog.Record` (time, level, message, attributes) and passes it to a `slog.Handler`, which decides whether to emit it and how to format it (https://pkg.go.dev/log/slog).
2. **Built-in handlers**: `slog.NewTextHandler` (`key=value`) and `slog.NewJSONHandler` (one JSON object per line). `HandlerOptions` sets `Level`, `AddSource` and `ReplaceAttr` (rename/drop/rewrite any attribute) (https://pkg.go.dev/log/slog#HandlerOptions).
3. **Attributes**: pass alternating key-value pairs (`"user", "perry"`) or typed `slog.Attr`s (`slog.String`, `slog.Int`, `slog.Group`). Typed attrs avoid boxing into `any` and cannot be mismatched.
4. **`With` and `WithGroup`**: `logger.With(attrs...)` returns a new logger that adds those attrs to every record; `WithGroup("req")` nests all later attrs under `req` (`req.method=GET` in text, `"req":{...}` in JSON).
5. **Levels**: `Debug(-4) < Info(0) < Warn(4) < Error(8)`. `HandlerOptions.Level` accepts a `slog.Leveler`; pass a `*slog.LevelVar` to change the level at runtime (https://pkg.go.dev/log/slog#LevelVar).
6. **`LogValuer`**: a type with `LogValue() slog.Value` is resolved only when the record is actually handled. Use it to redact secrets or to defer expensive work (https://pkg.go.dev/log/slog#LogValuer).
7. **Go 1.26**: `slog.NewMultiHandler(h1, h2, ...)` fans one record out to several handlers, each with its own level (https://pkg.go.dev/log/slog#NewMultiHandler). Go 1.24 added `slog.DiscardHandler`.

### How it works (standard library)
Based on `src/log/slog/logger.go` and `handler.go` in go1.26; internal details may change between versions.
- **Enabled is checked first**: `Logger.log` calls `handler.Enabled(ctx, level)` before building the `Record`. If it returns false, nothing else happens: no `Record`, no `LogValue()` call. That is why Part 3 shows `0` calls at INFO. But the **arguments themselves are still evaluated** by Go before `Debug(...)` is called, so `logger.Debug("x", "v", expensiveFunc())` still runs `expensiveFunc`. Only a `LogValuer` (or an explicit `if logger.Enabled(...)`) is lazy.
- **Record** stores the first few attrs inline (an array of 5 in the current implementation) and the rest in a slice, to avoid allocations for typical calls. `Record.Clone()` returns a copy with no shared state; call it before modifying a record you did not create, as the custom handler in Part 5 does.
- **`With` pre-formats**: the built-in handlers format the attrs passed to `WithAttrs` once and keep the bytes, so a logger created with `With(...)` is cheap to reuse per call. That is also why a custom handler **must** implement `WithAttrs`/`WithGroup` by wrapping the inner handler's result; returning the inner handler directly would drop the custom behavior for every derived logger.
- **`slog.SetDefault(l)`** also redirects the old `log` package (`log.Printf`) to `l`'s handler, so legacy code gets structured output.

## Full example
See `main.go`. Verified with `go run .` (3 runs, identical output), `go vet .` and `GOOS=windows go vet .` on go1.26 linux/amd64. Expected output is in [README.md](README.md).

Parts:
1. Key-value pairs vs typed attrs; a dangling key becomes `!BADKEY`; `Debug` is filtered at INFO.
2. `With("service", "api").WithGroup("req")` and `slog.Group` with the JSON handler.
3. `LevelVar` switched from INFO to DEBUG at runtime; a counting `LogValuer` proves lazy resolution.
4. A `Token` type whose `LogValue` returns `REDACTED`.
5. A custom `requestIDHandler` that reads `request_id` from the context in `InfoContext`.
6. `slog.NewMultiHandler`: text at INFO and JSON at WARN.

## Common pitfalls
1. **Mismatched key-value pairs**: `logger.Info("msg", "user")` or `logger.Info("msg", userID)` produces `!BADKEY=...`. `go vet` (the `slog` analyzer) catches literal calls, but not args built at runtime (Part 1). Prefer typed `slog.String(...)` in hot or important paths.
2. **Assuming disabled logs are free**: arguments are evaluated before the level check. `logger.Debug("dump", "body", string(hugeBytes))` converts the bytes even at INFO. Use a `LogValuer` or guard with `logger.Enabled(ctx, slog.LevelDebug)`.
3. **Custom handler forgets `WithAttrs`/`WithGroup` or ignores the context**: if `WithAttrs` returns `h.next.WithAttrs(...)` unwrapped, every `logger.With(...)` silently loses your handler. And `logger.Info(...)` passes `context.Background()`; only `InfoContext(ctx, ...)` gives your handler the request context.
4. **Logging secrets**: logging a whole struct (config, request, user) prints every field, tokens included. Give sensitive types a `LogValue` method (Part 4) or drop keys in `ReplaceAttr`.

## Exercise (no answer given)
Build a tiny HTTP-free "request pipeline" logger:
- Write a `levelOverrideHandler` that wraps another handler and lets **one component** (`component=db`) log at DEBUG while everything else stays at INFO. Hint: think about where the component name is known (`WithAttrs`) versus where the level is decided (`Enabled`).
- Use `slog.NewMultiHandler` so ERROR records also go to a JSON handler writing to a `bytes.Buffer`, and print how many bytes ended up in that buffer at the end.
- Add a `User` struct with `ID`, `Email` and `Password`. Implement `LogValue` so logs show `id` and a masked email (`p***@example.com`) but never the password.

Then answer: why does `logger.With("component", "db")` call your handler's `WithAttrs` only once, while `Enabled` runs on every log call? What does that mean for where you should do expensive work?

## Suggested directory
`slog_structured`
