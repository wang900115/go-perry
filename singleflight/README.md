# singleflight

Deduplicating concurrent identical requests with `golang.org/x/sync/singleflight`: callers that ask for the same key while a call is in flight share its result.

## What the code does

`main.go` starts 10 goroutines that all call `g.Do("get-google", getData)`. `getData` increments a mutex-protected counter and does `http.Get("https://google.com")`. Goroutines 0 to 3 start immediately; then `main` sleeps 1 s and starts goroutines 4 to 9. Each goroutine prints `status: <code>, shared: <bool>`, and at the end `main` prints `original function was called N times`.

Goroutines that call `Do` while a call for the key is still running wait for it and get its result instead of triggering a new request; when that happens every caller of that call, including the first, sees `shared: true`. The expected `N` is small (likely 2: one call for each burst), but it depends on how long the HTTP request takes relative to the 1 s gap and on goroutine scheduling, so it is not guaranteed.

## Concepts you can learn

- **Request coalescing / cache-stampede protection**: under a burst of identical requests, only one hits the backend.
- **`singleflight.Group.Do(key, fn)`** returns `(value, err, shared)`; `shared` reports whether the result was given to more than one caller.
- **Keys are only deduplicated while in flight**: once a call finishes, the next `Do` runs `fn` again. That is why the second burst after 1 s typically triggers a new call. It is not a cache.

## How to run

Own `go.mod` (`module singleflight`). Needs internet access to `https://google.com`.

```bash
cd singleflight
go run .
```

## Notes

- Response bodies are never closed, and all sharing callers receive the same `*http.Response`; reading its body from several goroutines would not be safe.
- Errors returned by `getDataSingleFlight` are ignored because it runs as `go getDataSingleFlight(...)`.
- The final `Printf` has no trailing newline.
