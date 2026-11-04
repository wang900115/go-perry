# pipeline

A small library of generic, context-aware channel pipeline helpers (`ReadOrDone`, `SendOrDone`, `FromSlice`, `Merge`, `Tee`, `Stage`) with unit tests. There is no runnable program; the tests are the way to exercise it.

## What the code does

`main.go` (package `main`, but **without a `func main`**) defines:

- `ReadOrDone[T](ctx, in)`: forwards values from `in` until `in` closes or `ctx` is done.
- `SendOrDone[T](ctx, out, v)`: sends `v` unless `ctx` is done first.
- `FromSlice[T](ctx, slice)`: returns a buffered channel pre-filled with the slice and already closed (`ctx` is unused).
- `Merge[T](ctx, cs...)`: fan-in; one goroutine per input, a `sync.WaitGroup` closes the output when all inputs are drained.
- `Tee[T](ctx, in)`: duplicates each value into two output channels buffered to 100.
- `Stage[In, Out](ctx, maxWorker, in, fn)`: runs `fn` on each input in its own goroutine, limiting concurrency with a `golang.org/x/sync/semaphore.Weighted`. Returns an output channel and an error channel; both are closed after it acquires the full semaphore weight, i.e. after all workers finish. Errors that are `context.Canceled` are dropped; semaphore errors are logged with `github.com/labstack/gommon/log`.

`main_test.go` tests each helper: ordering and cancellation for `ReadOrDone`, empty and non-empty `FromSlice`, `Merge` collecting 6 values from 3 channels, `Tee` delivering identical values to both outputs, `Stage` results/errors/max-2-concurrency/cancellation, and `SendOrDone` with a canceled context. Several cancellation sub-tests only `t.Logf` the count instead of asserting.

## Concepts you can learn

- **Generics for channel plumbing**: one implementation of `Merge`/`Tee`/`Stage` works for any element type.
- **Context-aware sends and receives**: `ReadOrDone`, `SendOrDone`, and `Merge` pair each blocking channel operation with `<-ctx.Done()` so a canceled pipeline does not leak goroutines; `Tee` and `Stage`'s worker sends do not (see Notes).
- **Fan-in with `WaitGroup`**: close the merged channel only after every producer is done.
- **Bounded parallelism with a weighted semaphore**, and waiting for all workers by acquiring the full weight before closing outputs.
- **Separate result and error channels** for a processing stage.

## How to run

Own `go.mod` (`module pipeline`). `go run .` and `go build .` fail with `function main is undeclared in the main package`. Use the tests:

```bash
cd pipeline
go test ./...
go test -v -run TestStage ./...
go test -race ./...     # requires cgo on Windows
```

## Notes

- `Stage`'s channels are unbuffered: the caller must drain **both** the output and error channels concurrently (as the tests do), or workers block.
- `Tee`'s comment says the sends "won't block", but the sends are not context-aware and do block once a 100-item buffer is full.
- `testMu` in `main_test.go` is declared but unused.
