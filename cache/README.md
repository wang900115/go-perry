# cache

A generic TTL cache wrapper around `github.com/jellydator/ttlcache/v3`, and an HTTP client decorator that caches whole responses keyed by the serialized request.

## What the code does

Module `cache` (`go 1.26.3`). Both files are `package main`, but there is **no `func main`**, so this is library-style code with nothing to run.

- **`simple_cache.go`**:
  - `SimpleCache[K, V]` interface: `Add`, `AddWithTTL`, `Get`, `GetWithLoader`, `Remove`, `Keys`, `Values`, `Len`, `OnExpiration`.
  - `NewSimpleCache` configures ttlcache with touch-on-hit disabled (reads do not extend TTL), optional capacity and default TTL, starts its expiry goroutine, and registers `runtime.AddCleanup` to `Stop()` it when the wrapper is garbage-collected.
  - `GetWithLoader` loads missing keys through a callback that returns value, TTL and error, and surfaces the loader's error.
  - `evictExpired` runs `DeleteExpired` at most once per hour, triggered by writes/loads.
  - `Keys`/`Values` skip expired items; `OnExpiration` registers a callback only for TTL expirations and returns an unregister function.
- **`cached_http_client.go`**:
  - `HTTPClient` wraps anything with `Do(*http.Request)`. Cache size limit 1000.
  - `Do` serializes the request (method, headers, URL, base64 body) to JSON and uses that string as the key. On a miss it rebuilds the request from the key, calls the wrapped client, and stores the full response in wire format (`resp.Write`). On a hit it parses the stored text back with `http.ReadResponse`.
  - Logs with Navidrome's `log.Trace` (key, whether it was cached, duration).

## Concepts you can learn

- **Generic interfaces over a third-party cache**: callers depend on `SimpleCache[K, V]`, not on ttlcache directly.
- **Loader pattern**: the cache, not the caller, decides when to fetch, which centralizes miss handling.
- **TTL semantics**: disabling touch-on-hit gives a fixed lifetime from insertion.
- **`runtime.AddCleanup`** (Go 1.24+): release background resources tied to an object without finalizers.
- **Decorator over `http.Client`**: same `Do` signature, transparent caching.
- **Request/response serialization**: turning HTTP messages into strings so they can be cache keys and values.

## How to run

```bash
cd cache
go vet ./...     # type-checks
```

`go build` / `go run .` fail with `function main is undeclared in the main package`. To try it, add a `main` that calls `NewHTTPClient(http.DefaultClient, time.Minute).Do(req)` twice.

## Notes

- Requires Go 1.26+: it uses `new(expr)` (for example `new(time.Now().Add(...))`).
- Bug: in `Do`, if `deserializeReq` fails, the inner `req` is `nil` and `log.Trace(req.Context(), ...)` panics.
- Reading `req.Body` to build the key consumes the caller's body.
- The request sent on a miss is rebuilt with `http.NewRequest`, so the caller's context (cancellation/deadline) is dropped.
- Every header is part of the key, so requests differing only in, say, a trace header never share a cache entry.
- The code imports `github.com/navidrome/navidrome` just for logging, which pulls in that project as a dependency. (Inference: the structure resembles Navidrome's own cache utilities; not verified.)
