# async

Background task processing with [asynq](https://github.com/hibiken/asynq), a Redis-backed task queue: a producer enqueues an "email" task and a consumer (worker server) picks it up and handles it.

## What the code does

The module (`module asy`, own `go.mod`) has three packages:

- `tasks/const.go`: task type names. Only `TypeEmail = "email:deliver"` is used; `TypeReport` and `TypeImageProcessing` are declared but have no handler or producer.
- `tasks/email.go`: `EmailPayload` (address, subject, body). `NewEmailTask` JSON-encodes the payload into an `asynq.Task`. `EmailTaskHandler` decodes it and only logs `Sending Email to User: email=..., subject=...`; no email is actually sent. A JSON decode error is wrapped with `asynq.SkipRetry` so asynq does not retry it.
- `producer/main.go`: connects to Redis at `127.0.0.1:6379`, enqueues one email task for `test@test.com`, and logs `enqueued task: id= <id> queue= <queue>` (the queue is `default`, since no queue option is passed).
- `consumer/main.go`: starts an `asynq.Server` with `Concurrency: 10` and weighted queues `critical: 6`, `default: 3`, `low: 1`, registers `EmailTaskHandler` for `email:deliver` on a `ServeMux`, and runs until stopped.

## Concepts you can learn

- **Producer / consumer decoupling through a queue**: the producer only writes to Redis and exits; work happens later in a separate process, which survives restarts and spreads load.
- **Task type routing**: `ServeMux.HandleFunc(type, handler)` maps a task type string to a handler, like HTTP routing.
- **Serializing payloads**: tasks carry `[]byte`, so the payload struct is JSON-marshalled on enqueue and unmarshalled in the handler.
- **Retry control**: returning an error makes asynq retry; wrapping with `asynq.SkipRetry` marks errors (like a malformed payload) that retrying cannot fix.
- **Weighted priority queues**: the consumer polls `critical`/`default`/`low` in a 6:3:1 ratio.

## How to run

Prerequisite: a Redis server listening on `127.0.0.1:6379` (address is hard-coded), for example `docker run -p 6379:6379 redis`.

```bash
cd async
go run ./consumer      # terminal 1: keeps running, logs each handled task
go run ./producer      # terminal 2: enqueues one task and exits
```

## Notes

- The module root has no `main` package, so `go run .` inside `async/` does not work; run the `consumer` and `producer` sub-packages as shown.
- The Redis address is a constant in both `main.go` files; change it there if your Redis is elsewhere.
