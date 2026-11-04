# workerpool

A worker pool built with the [tunny](https://github.com/Jeffail/tunny) library, sized to the number of CPUs, used to "send" 100 emails. See `worker_pool/` for a hand-written version.

## What the code does

`main.go`:

- Prints `Number of CPUs: N` (`runtime.NumCPU()`) and creates `tunny.NewFunc(numCPUs, fn)`, a pool of `N` workers.
- `fn` receives an `interface{}` payload, type-asserts it to `map[string]interface{}`, extracts `email`, `subject`, and `body`, and calls `SendEmail`, which prints the email and sleeps 2 s. A missing field returns an `error` value.
- Starts 100 goroutines, each calling `pool.Process(data)` (blocks until a worker is free and returns the result). A `nil` result prints `Mail sent!`.
- Polls `pool.QueueLength()` once per second, printing it, until it reaches 0, then sleeps 3 s and exits.

With `N` workers and 2 s per email, sending all 100 takes about `100 / N * 2` seconds.

## Concepts you can learn

- **Library-managed worker pools**: tunny keeps `N` long-lived workers and queues `Process` calls, so callers get bounded concurrency with a synchronous API.
- **Sizing a pool to `runtime.NumCPU()`**: a common default for CPU-bound work. This workload is a sleep, so it is effectively I/O-bound and the CPU count is just a convenient number.
- **`interface{}` payloads and type assertions** with the `, ok` form to fail safely on bad input.
- **Observing in-flight work** with `QueueLength()` (jobs waiting for a worker plus jobs running).

## How to run

Own `go.mod` (`module workerpool`):

```bash
cd workerpool
go run .
```

## Notes

- The exit condition is racy. The first `QueueLength()` check can run before any goroutine has called `Process`, see 0, and stop the loop; then `main` exits after 3 s with most emails unsent. (`QueueLength()` counts every in-flight `Process` call, waiting or running, so once the loop has seen a non-zero value, 0 does mean all jobs finished.) A `sync.WaitGroup` would be the reliable way to wait.
- Error results from `Process` are ignored (only `nil` is reported).
- `go.mod` marks `github.com/Jeffail/tunny` as `// indirect` even though it is imported directly; `go mod tidy` would fix that.
