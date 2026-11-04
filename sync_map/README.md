# sync_map

An HTTP Fibonacci service (Gin) that caches results in a `sync.Map` and expires each entry after 5 minutes.

## What the code does

`main.go` starts a Gin server on `:8080` with one route, `GET /fibonacci?n=<int>`:

- If `n` is not an integer, it returns `400` with `{"error":"Invalid input"}`. A missing `n` defaults to `0`.
- On a cache hit (`cache.Load(n)`) it prints `Cache Hit!` to the server console and returns `{"result": ...}`.
- On a miss it prints `Cache miss!`, computes the value iteratively with `math/big` (so large `n` does not overflow), stores it, and starts a goroutine that deletes the entry after 5 minutes.

`calculateFibonacci` returns 0 for `n <= 0` and 1 for `n == 1`; for example `n=10` returns `55`.

## Concepts you can learn

- **`sync.Map`**: a concurrency-safe map with no explicit locking. It suits caches where keys are written once and read many times, which is the case here.
- **Concurrent HTTP handlers**: Gin runs each request in its own goroutine, so the shared cache must be safe for concurrent use; a plain `map` would race.
- **Arbitrary-precision integers** with `math/big` for values that exceed `int64`.
- **Time-based expiry** with a per-entry goroutine + `time.Sleep` (simple, but see Notes).

## How to run

Own `go.mod` (`module sync_map`).

```bash
cd sync_map
go run .
# in another terminal
curl "http://localhost:8080/fibonacci?n=10"     # Cache miss!
curl "http://localhost:8080/fibonacci?n=10"     # Cache Hit!
```

## Notes

- One goroutine per cached key sleeping 5 minutes is fine for a demo but does not scale to many keys; a real cache would use a single janitor or a TTL library.
- Two concurrent misses for the same `n` both compute the value and both start an expiry goroutine.
- There is no upper bound on `n`, so a very large value makes the handler do a lot of work.
- The `r.Run(port)` error is ignored.
