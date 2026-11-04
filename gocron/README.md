# gocron

In-process job scheduling with [`go-co-op/gocron/v2`](https://github.com/go-co-op/gocron). The program sets up an interval job, a cron-expression job, and a daily job at fixed times.

## What the code does

`main.go` creates one scheduler and registers three jobs:

1. **Duration job**: every 30 seconds, logs `Every 30 seconds`. It's named `job: Every 30 seconds` and has three event listeners:
   - `BeforeJobRuns` logs `Job starting: <uuid>, <name>`.
   - `AfterJobRuns` logs `Job completed: ...`.
   - `AfterJobRunsWithError` logs `Job has an error: ...`.

   The job's UUID is logged once at startup.
2. **Cron job**: `*/10 * * * *` (5-field cron without seconds, because `withSeconds=false`). It runs every 10 minutes and logs `Cronjob: Every 10 mins`.
3. **Daily job**: interval 1 day, at 18:07:00 and 05:30:00 local time. It logs `Dailyjob  Runs everyday`.

`s.Start()` starts the scheduler. A goroutine waits for SIGINT or SIGTERM, calls `s.Shutdown()`, and exits.

## Concepts you can learn

- **Three ways to define a schedule**: fixed interval (`DurationJob`), cron expression (`CronJob`), and calendar times (`DailyJob` + `NewAtTimes`).
- **Tasks with arguments**: `gocron.NewTask(fn, args...)` binds parameters to the job function when you register it.
- **Job lifecycle hooks**: event listeners give you a place for metrics or logging around every run, without changing the job itself.
- **Graceful shutdown**: `Shutdown()` stops the scheduler before the process exits.

## How to run

```bash
cd gocron
go run .
```

You don't need any external services. Stop it with Ctrl+C.

## Notes

- `main` ends with an empty `for {}` loop, which keeps one CPU core at 100%. Blocking on the signal channel (or `select {}`) would avoid that.
- None of the `s.NewJob` errors are checked. For the first job, `err` is assigned but never read, so a failed registration would call `j.ID()` on a nil job.
- The task functions can't return an error, so the `AfterJobRunsWithError` listener never fires.
- `go.mod` marks every dependency as `// indirect`, even those imported directly.
