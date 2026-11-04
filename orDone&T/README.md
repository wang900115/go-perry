# orDone&T

The **or-done** and **tee** channel patterns, applied to a toy transaction stream that is split into a fraud-detection pipeline and an analytics pipeline.

## What the code does

`main.go`:

- `transactionGenerator(done)` emits up to 200 random amounts (`rand.Intn(1000)`), sleeping 100 ms after each send, and stops early if `done` is closed (printing `Transaction Generator: Shutting down ...`).
- `tee(input)` copies every value to two output channels. It sends to both using a `select` and sets each channel variable to `nil` after it has been sent to, so the value goes to each output exactly once in whichever order is ready first.
- `orDone(done, input)` wraps a channel so a `for range` over it stops as soon as `done` is closed, even if `input` is never closed.
- `fraudDetectionPipeline` prints `Suspicious transaction detected` for amounts > 700, otherwise `Normal transaction`.
- `analyticsPipeline` prints each transaction and, when it stops, the integer average.
- `main` wires generator -> tee -> both pipelines, lets them run 10 s, closes `done`, waits 1 s, and prints `Main: All pipeline stopped`.

Since the generator produces about 10 values per second, roughly 100 of the 200 transactions are processed before shutdown. Values are random, so output differs every run.

## Concepts you can learn

- **Or-done**: reading from a channel you do not control can block forever; wrapping it in a `select` with `done` makes the consumer cancellable.
- **Tee**: fan one stream out to several consumers; each consumer receives every value.
- **`nil` channels in `select`**: a send to a `nil` channel never proceeds, so setting `o1 = nil` disables that case after it fires.
- **Closing a channel as a broadcast**: `close(done)` wakes every goroutine waiting on it.
- **`defer close(out)`** in each stage so downstream `range` loops end.

## How to run

This directory has no `go.mod`; it belongs to the root module. The `&` in the directory name must be quoted in the shell. From the repo root:

```bash
go run "orDone&T/main.go"
```

Runs for about 11 s.

## Notes

- `tee` does not watch `done`. If it has received a value when the pipelines stop reading, it stays blocked on that send forever; otherwise it exits when the generator closes `source`. Either way the program only ends cleanly because `main` exits.
- The shutdown messages (`... shut down`, the average) may print in any order, and `main` waits only a fixed 1 s rather than synchronizing with the pipelines.
- Code comments are in Chinese.
