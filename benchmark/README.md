# benchmark

Record an execution trace with `runtime/trace`. Despite the directory name, there are no `testing` benchmarks here.

## What the code does

`main.go` creates `trace.out` in the current directory and starts tracing. It prints `Hello world` 1000 times, then stops the trace (via `defer`). The closing comment contrasts the two tools:

- **Profile**: resource usage (CPU, memory) to find hotspots.
- **Trace**: a timeline of execution events (goroutines, syscalls, GC, scheduling) to diagnose latency and concurrency problems.

The directory also contains committed artifacts: `trace.out`, `cpu.prof` and `mem.prof` (gzip-compressed pprof data). The code only produces `trace.out`. The code does not show how the `.prof` files were generated.

## Concepts you can learn

- **`runtime/trace`**: `trace.Start(w)` / `trace.Stop()` capture a detailed event timeline with low enough overhead for short production captures.
- **`go tool trace`**: opens the timeline in a browser, where you can see goroutine scheduling, blocking on syscalls (here, the 1000 `Println` writes to stdout) and GC.
- **Profile vs trace**: a profile says where time is spent in aggregate, and a trace says when and why goroutines were running or waiting.

## How to run

```bash
cd benchmark
go run .
go tool trace trace.out

# inspect the committed profiles
go tool pprof cpu.prof
go tool pprof mem.prof
```

## Notes

- Errors from `os.Create` and `trace.Start` are ignored, and the file is never explicitly closed.
- Running the program overwrites the committed `trace.out`.
