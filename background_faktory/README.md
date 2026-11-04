# background_faktory

Background job processing with [Faktory](https://github.com/contribsys/faktory). A producer pushes a job, and a worker processes it and chains a follow-up job.

## What the code does

There are two sub-programs, and each has its own `go.mod`:

- `producer/main.go` connects with `faktory.Open()` and pushes one job: type `report`, arg `"test@test.com"`, queue `critical`. If the push fails, it prints `Error pushing job`.
- `worker/main.go` creates a `faktory_worker_go` manager with concurrency 5 and a shutdown timeout of 25s. It processes queues in strict priority order: `critical`, then `default`, then `low_priority`. It registers two job types:
  - `report(addr)`: logs the job ID, prints `Preparing report for the user: <addr>`, and sleeps 10s. It then pushes an `email` job with args `(addr, "Report is ready!")` to the default queue, using `help.With`.
  - `email(addr, subject)`: logs the job ID, prints `Sending mail to <addr> with subject <subject>`, and sleeps 5s.

The full flow is: run the producer once, the worker picks up `report` from `critical`, then it enqueues and runs `email` from `default`.

## Concepts you can learn

- **Producer/worker decoupling through a job server**: the producer returns immediately, and the slow work runs elsewhere with retries and persistence handled by Faktory.
- **Strict priority queues**: `ProcessStrictPriorityQueues` always drains `critical` before it looks at lower queues. With non-strict (weighted) queues, lower queues can still be served.
- **Job chaining**: a job can enqueue follow-up jobs. `help.With` borrows a client from the worker's connection pool.
- **Job helpers**: `worker.HelperFor(ctx)` exposes the job ID and related data inside the handler.
- **Graceful shutdown**: `ShutdownTimeout` gives in-flight jobs time to finish when the worker is stopped.

## How to run

Prerequisite: a running Faktory server. `faktory.Open()` reads `FAKTORY_URL` (or `FAKTORY_PROVIDER`) and otherwise uses `localhost:7419` (that default comes from the faktory client library, not this repo). For example:

```bash
docker run --rm -p 7419:7419 -p 7420:7420 contribsys/faktory
```

Then, in two terminals:

```bash
cd background_faktory/worker
go run .
```

```bash
cd background_faktory/producer
go run .
```

The Faktory web UI is served on port 7420 by the server, not by this code.

## Notes

- The two modules pin different `faktory` versions (producer v1.9.2, worker v1.8.0), and both mark their direct dependencies as `// indirect`.
- The `email` handler type-asserts `args[0]` and `args[1]` without checks, so a job with missing args panics inside the handler.
- The producer never closes its client.
- Not run against a live Faktory server.
