# context

Five short programs covering the basic `context` package APIs: the root context, cancellation, timeouts, and request-scoped values.

## What the code does

- `01.go`: prints the root context from `context.Background()` (output: `This is the root context: context.Background`).
- `02.go`: `context.WithCancel`. A goroutine blocks on `ctx.Done()`; `main` calls `cancel()` after 1 s and the goroutine prints `Goroutine cancelled: context canceled`. `main` then sleeps 500 ms so the print has time to happen.
- `03.go`: a 2 s `context.WithTimeout` raced against `time.After(3 * time.Second)` in a `select`. The context fires first and prints `Timeout: context deadline exceeded`.
- `04.go`: `context.WithValue` with a custom key type `contextKey`. `process` reads the value back and prints `User ID from context: 42`.
- `05.go`: a 3 s timeout around a 5-step loop where each step sleeps 1 s. It prints `Working on step 1` to `3`, then typically `Stopped early: context deadline exceeded` instead of finishing all 5 steps.

## Concepts you can learn

- **`context.Background()`** is the root of every context tree; derived contexts add cancellation, deadlines, or values.
- **Cancellation propagates through `ctx.Done()`**: goroutines `select` on it to stop cooperatively, and `ctx.Err()` says why (`context.Canceled` vs `context.DeadlineExceeded`).
- **Always call `cancel`** (`defer cancel()` in 03 and 05) to release the timer and resources a derived context holds.
- **Typed context keys**: an unexported named type (`contextKey`) instead of a plain `string` avoids key collisions between packages.
- **Checking for cancellation inside a loop** (05) with `select { case <-ctx.Done(): ... default: ... }`.

## How to run

`context/` has no `go.mod`; it belongs to the root module, and each file declares its own `func main`. From the repo root:

```bash
go run context/01.go
go run context/02.go
go run context/03.go
go run context/04.go
go run context/05.go
```

## Notes

- In `05.go` the 3 s deadline and the end of the third 1 s sleep land at almost the same moment, and cancellation is only checked between sleeps, so whether step 4 starts depends on timing.
