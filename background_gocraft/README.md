# background_gocraft

Redis-backed background jobs with [`gocraft/work`](https://github.com/gocraft/work). An enqueuer pushes jobs, and a worker pool processes them through middleware.

## What the code does

There are two sub-programs, and each has its own `go.mod` (module names `1` and `2`). Both use the namespace `demo_app` and Redis at **`localhost:6379`**.

- `enqueuer/enqueuer.go` enqueues two jobs and exits:
  - `email` with `{userID: 10, subject: "testing"}`
  - `report` with `{userID: 5}`
- `process_jobs/worker.go` starts a worker pool with concurrency 10 and a `Context` type. It runs until SIGINT or SIGTERM, then calls `pool.Stop()`.
  - Middleware runs in this order:
    1. `Log` prints `Starting a new Job: <name> with ID: <id>`.
    2. `FindCurrentUser` builds a fake `User{ID, Email: "test<ID>@test.com"}` from the `userID` arg and stores it on the context.
  - `email` (priority 10, MaxFails 1) prints `Sending email to test10@test.com with subject testing` and sleeps 2s.
  - `report` (priority 10, MaxFails 1) prints `Preparing report ...`. It then loops 360 times, sleeping 10s each time and calling `job.Checkin("i = N")`, which takes about 1 hour. Finally it enqueues an `email` job with subject `"Report is ready"`.

## Concepts you can learn

- **Typed job context + middleware**: gocraft/work creates a fresh `Context` per job. Middleware (`func(*Context, *Job, NextMiddlewareFunc)`) can load shared data, such as the current user, before the handler runs.
- **Job options**: `Priority` weights how often the pool picks each job type, and `MaxFails` caps retries before a job goes to the dead queue.
- **Check-ins for long-running jobs**: `job.Checkin` records progress that you can see in gocraft's tooling, which is useful for jobs that run for minutes or hours.
- **Argument helpers**: `ArgInt64` and `ArgString` collect conversion errors, and `ArgError()` reports them in one check.
- **Redis connection pooling**: `redigo`'s `redis.Pool` with `MaxActive` and `Wait: true` blocks callers instead of opening unlimited connections.

## How to run

Prerequisite: Redis on `localhost:6379`. For example:

```bash
docker run --rm -p 6379:6379 redis
```

Start the worker first, then the enqueuer:

```bash
cd background_gocraft/process_jobs
go run .
```

```bash
cd background_gocraft/enqueuer
go run .
```

## Notes

- The `report` job deliberately runs for about 1 hour.
- `sendEmail` dereferences `c.currentUser`, which is only set when the job has a `userID` arg. An `email` job without `userID` would panic with a nil pointer.
- The enqueuer ignores the error from the second `Enqueue` call. So does the worker's chained `Enqueue` in `report`.
- Not run against a live Redis.
