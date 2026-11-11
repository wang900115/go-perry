# 2026-11-11 · defer / panic / recover

## Why learn this
- **Servers must not die because one request panics.** `net/http` recovers panics per connection; gRPC interceptors, job runners and worker pools need the same "recover and turn it into an error" wrapper.
- **Panics in goroutines kill the whole process.** Every `go` statement in a worker pool or pipeline is a place where one bad input can crash everything.
- **Cleanup that always runs**: `defer f.Close()`, `defer mu.Unlock()`, `defer cancel()`. You need to know exactly when the deferred call runs and which values it sees.
- It connects the previous lessons: `error_wrapping` (a recovered panic becomes a wrapped error) and the concurrency series (goroutine boundaries).

## Core concepts
1. `defer f(args)` evaluates `f` and `args` **immediately**, but runs the call when the surrounding function returns (normally or by panic). Multiple defers run in **LIFO** order (https://go.dev/ref/spec#Defer_statements).
2. A deferred closure can read and modify **named result parameters** after `return` has assigned them.
3. `panic(v)` stops normal execution of the current function, runs the goroutine's deferred calls one by one, going up the stack. If nothing recovers, the program prints the panic value and stack traces and exits with status 2 (https://go.dev/ref/spec#Handling_panics).
4. `recover()` stops the panicking sequence and returns the panic value, **but only when called directly by a deferred function**. Anywhere else it returns `nil` (spec, same link).
5. Runtime errors (index out of range, nil map write, nil dereference) panic with a value that implements `runtime.Error` (https://pkg.go.dev/runtime#Error).
6. Since Go 1.21, `panic(nil)` panics with `*runtime.PanicNilError`, so `recover()` never returns `nil` for a real panic (https://go.dev/doc/go1.21#language, https://pkg.go.dev/runtime#PanicNilError).
7. `runtime.Goexit()` runs deferred calls but is not a panic: `recover()` returns `nil` (https://pkg.go.dev/runtime#Goexit).

Background reading: https://go.dev/blog/defer-panic-and-recover

### How it works (runtime)
Based on `src/runtime/panic.go` and the compiler; internal details change between versions.
- **Open-coded defers (Go 1.14+)**: if a function has a small number of defers (at most 8) and none are in a loop, the compiler inlines the deferred calls at each return point and tracks which ones are active with a bitmask. This makes `defer` almost free (https://go.dev/doc/go1.14#runtime). Defers in a loop (like Part 1's `for` loop) fall back to `_defer` records linked onto the goroutine.
- **`gopanic`**: creates a panic record on the goroutine and runs the pending defers, including open-coded ones found by scanning stack frames using compiler metadata.
- **`gorecover`**: checks that it is called from the deferred function that the current panic is running, at the expected stack frame. A helper called by that deferred function is one frame deeper, so the check fails and it returns `nil` (Part 3).
- After a successful recover, the runtime unwinds the stack to the frame of the function that registered the defer and makes it **return normally to its caller**, with whatever values its named results hold. That is why `safeCall` needs `(err error)` as a named result.
- Panics are per goroutine: the defer chain belongs to the goroutine, so a parent goroutine cannot recover a child's panic.

## Full example
See `main.go`. Verified with `go run .` (3 runs, identical output) and `go vet` on go1.26 linux/amd64; `GOOS=windows go vet .` also passes. Expected output is in [README.md](README.md).

Parts:
1. Defer order, argument evaluation vs closure capture, modifying a named result.
2. `safeCall`: panic -> wrapped error (`%w` twice, so both `ErrPanic` and the original `runtime.Error` are in the chain).
3. `recover` inside a helper returns `nil`; `Goexit` is not a panic.
4. `panic(nil)` -> `*runtime.PanicNilError`.
5. `goSafe`: recover inside each goroutine and report through a channel.

## Common pitfalls
1. **Recover in a helper**: `defer func() { logPanic() }()` where `logPanic` calls `recover()` does nothing; the panic continues. Call `recover()` directly in the deferred function (or `defer logPanic()` so that `logPanic` itself *is* the deferred function).
2. **Recover in the parent for a child goroutine**: a `defer recover...` in `main` does not catch a panic in `go work()`; the process crashes. Each goroutine needs its own recover (Part 5).
3. **Unnamed result + recover**: with `func f() error { defer func(){ recover() }(); ... }` the function returns the zero value (`nil`), silently hiding the failure. Use a named result and set it.
4. **Defer in a loop**: `for _, f := range files { defer f.Close() }` keeps every file open until the function returns. Move the body into a function or close explicitly.
5. **Re-panicking**: if you recover only to log and then `panic(r)` again, the crash message reads `panic: first [recovered, repanicked]` instead of printing the value twice (verified on go1.26; the wording is tested in `runtime/crash_test.go`). I believe this format was introduced in Go 1.25 (https://go.dev/doc/go1.25), but I could not open the release notes from this environment to confirm, so treat the version as **unverified**.

## Exercise (no answer given)
Write a `Recover` HTTP middleware **without** using a real server in `main` (use `net/http/httptest`):
- `func Recover(next http.Handler) http.Handler` returns 500 and the body `internal error` when `next` panics, and logs the panic value plus `runtime/debug.Stack()`.
- It must **not** swallow `http.ErrAbortHandler` (re-panic with it; see https://pkg.go.dev/net/http#ErrAbortHandler).
- Test three handlers with `httptest.NewRecorder`: one that works, one that indexes out of range, one that calls `panic(nil)`. Print the status codes.
- Bonus: what happens if the handler already wrote the header (`w.WriteHeader(200)`) before panicking? Explain what the client sees.

## Suggested directory
`defer_panic_recover`
