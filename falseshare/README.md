# falseshare

False sharing: two goroutines updating different variables that happen to sit on the same CPU cache line, and how padding fixes it.

## What the code does

- `main.go` defines two structs:
  - `Counter{A, B int64}`: `A` and `B` are adjacent, so they normally share one 64-byte cache line.
  - `CustomCounter{A int64; _ cpu.CacheLinePad; B int64}`: `golang.org/x/sys/cpu.CacheLinePad` inserts padding sized for the target architecture's cache line, so `A` and `B` end up on different lines.
- `counter_test.go` (external test package) has two benchmarks. Each starts two goroutines: one runs `atomic.AddInt64(&counter.A, 1)` `b.N` times and the other does the same for `B`.
  - `BenchmarkFalseSharing` uses `Counter`.
  - `BenchmarkPaddedCounter` uses `CustomCounter`.

The padded version is expected to be faster per op, and it was in a local run. The ratio depends on the CPU.

## Concepts you can learn

- **Cache coherence works per line, not per variable**: when one core writes to a line, other cores must invalidate their copy. Two cores writing to unrelated fields on the same line keep stealing it from each other.
- **Padding**: separating hot, independently written fields by at least one cache line removes that contention, at the cost of memory.
- **`cpu.CacheLinePad`**: a portable padding type, so you do not hard-code 64 bytes.
- **Benchmarking concurrency**: the benchmark needs at least 2 Ps (`GOMAXPROCS >= 2`) for the goroutines to actually run in parallel.

## How to run

```bash
cd falseshare
go test -bench . -benchmem ./...
go test -bench . -cpu 1,2,4 ./...   # compare parallelism levels
```

`go build` / `go run` do not work: `main.go` is `package main` without `func main`, so the build fails with `function main is undeclared`. Only `go test` is meaningful.

## Notes

- The module is named `falsesharee` (with a double `e`). The test imports it under that path.
- Each benchmark iteration count `b.N` applies per goroutine, so ns/op is the wall time for `b.N` increments on each of the two counters.
