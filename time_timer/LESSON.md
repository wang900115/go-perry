# 2026-11-10 · time.Timer / Ticker / AfterFunc and the Go 1.23 timer semantics

## Why learn this
- **Timeouts in event loops**: "close the connection if nothing arrives for 30s", "flush the batch after 100ms of quiet". This is a `select` with a timer that is reset on every event.
- **Periodic work**: heartbeats, metrics flush, cache cleanup with `time.Ticker`.
- **Delayed callbacks**: retry with backoff, debounce, lease expiry with `time.AfterFunc`.
- **Old code is full of workarounds** (`if !t.Stop() { <-t.C }`) that were needed before Go 1.23. You need to know why they existed and when they are no longer required.

This continues the concurrency series (`channel`, `context`, `sync_*`): a timer is just a channel that the runtime sends to.

## Core concepts
1. `time.NewTimer(d)` returns a `*Timer` whose channel `C` receives the current time once, after `d`. `time.After(d)` is `NewTimer(d).C` (https://pkg.go.dev/time#NewTimer, https://pkg.go.dev/time#After).
2. `t.Stop()` prevents the timer from firing; it returns `true` if the call stopped the timer, `false` if it had already expired or been stopped (https://pkg.go.dev/time#Timer.Stop).
3. `t.Reset(d)` changes the timer to expire after `d` (https://pkg.go.dev/time#Timer.Reset).
4. `time.NewTicker(d)` sends repeatedly; it **drops** ticks for slow receivers rather than queueing them (https://pkg.go.dev/time#NewTicker).
5. `time.AfterFunc(d, f)` calls `f` in its own goroutine; its `C` field is nil (https://pkg.go.dev/time#AfterFunc).

### What changed in Go 1.23 (https://go.dev/doc/go1.23#timer-changes)
- **Unreferenced timers/tickers are garbage-collected** even if `Stop` was never called. `time.After` in a loop no longer leaks until the timer fires.
- **Timer channels are synchronous** (unbuffered): `len(t.C)` and `cap(t.C)` are 0. After `Stop` or `Reset` returns, **no stale value from before the call will be received**.
- These apply only when the main module's `go.mod` says `go 1.23` or later. `GODEBUG=asynctimerchan=1` restores the old behaviour.

### How it works (runtime)
- Each P (logical processor in the scheduler) keeps its own heap of timers ordered by expiry time (`runtime/time.go`). The scheduler checks the heap when it looks for work, and the netpoller (epoll on Linux, IOCP on Windows) is given a timeout equal to the nearest timer, so a sleeping `M` wakes up in time. There is no dedicated "timer goroutine".
- Before 1.23 the channel had a 1-element buffer, and the runtime sent into it the moment the timer expired. If nobody received, the value sat in the buffer, which is why `Reset` could be followed by an immediate stale receive and why the drain idiom was needed.
- Since 1.23 the channel still has an internal 1-element slot, but the runtime reports `len`/`cap` as 0 and, on `Stop`/`Reset`, discards any value that was not yet received. In the example, `Stop` on an expired-but-unreceived timer returns `true` on go1.26. That return value is what I observed; I am **not certain** the docs guarantee it, so do not build logic on it. What *is* documented is that no stale value will be received after `Stop` returns.

## Full example
See `main.go` in this directory. Verified with `go run .` (3 runs, identical output) and `go vet` on go1.26 linux/amd64; `GOOS=windows go vet .` also passes. Expected output: see [README.md](README.md). Running with `GODEBUG=asynctimerchan=1` shows the old behaviour (verified).

Parts 2-3 depend on timing (tens of milliseconds of margin). On a very loaded machine they could in theory print a different result.

## Common pitfalls
1. **`time.After` in a hot loop**: since Go 1.23 it is collected, but it still allocates a new timer on every iteration, and it resets the timeout on every loop iteration (including iterations triggered by other cases). Use one `NewTimer` + `Reset` when you mean an idle timeout.
2. **Copying old drain idioms into Go 1.23+ code**: `if !t.Stop() { <-t.C }` can now **block forever**, because after `Stop` the channel is guaranteed to be empty. Just call `Stop`/`Reset`. (The idiom is still needed if your module's `go.mod` is below 1.23.)
3. **Assuming `AfterFunc`'s `Stop` waits for `f`**: if `f` already started, `Stop` returns `false` and does **not** wait for `f` to finish. `f` runs in another goroutine, so shared data needs synchronization (the example uses `atomic.Int32`).
4. **Forgetting `Ticker.Stop`**: since 1.23 an unreferenced ticker is collected, but a ticker still referenced by a running goroutine keeps ticking. Always `defer tk.Stop()`.

## Exercise
Write a `Debouncer` that delays calling a function until no new `Trigger()` has happened for `d`:

```go
type Debouncer struct { /* ... */ }
func NewDebouncer(d time.Duration, f func()) *Debouncer
func (db *Debouncer) Trigger()
func (db *Debouncer) Stop()
```

Requirements:
- 10 calls to `Trigger()` 10ms apart with `d = 50ms` call `f` exactly once, about 50ms after the last call.
- `Trigger()` is safe to call from several goroutines.
- After `Stop()`, `f` is never called again.
- Implement it twice: once with `time.AfterFunc`, once with a goroutine + `NewTimer` + `select`. Which one is easier to make correct? Run it with `go run -race`.

## Suggested directory
`time_timer/`
