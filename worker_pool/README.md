# worker_pool

A hand-written worker pool: a fixed number of goroutines pull tasks of different types from a shared channel. See `workerpool/` for the same idea built on the `tunny` library.

## What the code does

- `worker_pool.go`: defines a `Task` interface with one method, `Process()`, and two implementations: `EmailTask` (prints `Sending Email to :  <email>`, sleeps 2 s) and `ImageProcessingTask` (prints `Processing the image :  <url>`, sleeps 5 s). `WorkerPool.Run()` creates a channel buffered to `len(Tasks)`, starts `concurrency` worker goroutines that `range` over it, enqueues every task, closes the channel, and waits on a `sync.WaitGroup` (one `Done` per processed task).
- `main.go`: builds 20 tasks (10 email, 10 image, alternating) and runs them with `concurrency: 5`.

Up to 5 tasks run at a time. Output order is not fixed. With 10 x 2 s + 10 x 5 s = 70 s of work spread across 5 workers, 14 s is the lower bound; with this alternating task order the run takes about 16 s.

## Concepts you can learn

- **Worker pool**: bound concurrency to `N` goroutines instead of one goroutine per task, which caps resource usage.
- **Interfaces for heterogeneous work**: the pool only knows `Task.Process()`, so new task types need no pool changes.
- **Closing the job channel** tells workers to exit their `range` loop once the queue is drained.
- **`sync.WaitGroup` per task** (not per worker) to know when all work is finished.

## How to run

Own `go.mod` (`module worker_pool`); `main.go` and `worker_pool.go` are one package, so run the directory, not a single file:

```bash
cd worker_pool
go run .
```

## Notes

- `wg.Add(len(wp.Tasks))` is called after the workers start but before any task is sent, so it is still safe: no `Done` can happen before the `Add`.
- `WorkerPool` contains a `sync.WaitGroup`, so it must not be copied after use; `main` calls `Run` on a local variable through a pointer receiver, which is fine.
