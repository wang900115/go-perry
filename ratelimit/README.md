# ratelimit

Per-client HTTP rate limiting with `golang.org/x/time/rate` (token bucket), plus a client that sends requests fast enough to trigger it.

## What the code does

The module (`module ratelimiter`) has two programs:

- `server/main.go`: listens on `0.0.0.0:8080` and serves `{"message":"Hello, World!"}` on `/`. `rateLimiterMiddleware` extracts the client IP from `r.RemoteAddr`, keeps one `rate.Limiter` per IP in a `sync.Map` (`LoadOrStore`), and allows 2 requests/second with a burst of 10. Rejected requests get `429 Too Many Requests` and `{"error":"Too many request"}`.
- `client/main.go`: sends 50 `GET http://localhost:8080/` requests, 100 ms apart (about 10 req/s), printing `Request: N: Status: <code>` for each.

Expected client output: about the first 10 requests return `200` (the burst), after which most return `429` and roughly one in five returns `200` as tokens refill at 2 per second. The exact pattern depends on timing.

## Concepts you can learn

- **Token-bucket rate limiting**: `rate.NewLimiter(r, b)` refills `r` tokens per second up to `b`; `Allow()` takes a token or reports `false` without blocking.
- **Per-key limiters**: storing one limiter per client IP in a `sync.Map` lets each client have its own budget.
- **HTTP middleware**: wrapping an `http.Handler` to reject requests before they reach the real handler.
- **Using `429` correctly** for throttled requests.

## How to run

Own `go.mod`; the module root has no `main` package, so run the sub-packages:

```bash
cd ratelimit
go run ./server     # terminal 1
go run ./client     # terminal 2
```

## Notes

- The per-IP map is never cleaned up, so it grows with every distinct client IP.
- `LoadOrStore` builds a new `rate.Limiter` on every request, even when one already exists (it is then discarded).
- Behind a proxy or load balancer, `RemoteAddr` is the proxy's address, so all clients would share one limiter.
- The client never closes `resp.Body`, and its `Response` type is unused.
- `go.mod` marks `golang.org/x/time` as `// indirect` even though the server imports it directly.
