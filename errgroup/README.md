# errgroup

Fetching several URLs concurrently with `golang.org/x/sync/errgroup`, limited to 2 in flight, and collecting the first error.

## What the code does

`main.go` has a list of four URLs, the last one deliberately invalid (`https://invalid-url`). For each URL it calls `g.Go(...)` on an `errgroup.Group` with `g.SetLimit(2)`, so at most two fetches run at once. `fetchPage` does an `http.Get`, reads the body, and stores it in a shared map guarded by a `sync.Mutex`, printing `Successfully fetched ...` lines along the way.

`g.Wait()` returns the first non-nil error. Because `https://invalid-url` cannot be fetched, the expected result is `failed to fetch https://invalid-url: ...` followed by `Error occured:  failed to fetch https://invalid-url: ...`. The "all fetched" branch, which prints the first 100 bytes of each body, only runs if every fetch succeeds.

## Concepts you can learn

- **`errgroup.Group`**: like `sync.WaitGroup`, but each goroutine returns an `error` and `Wait` reports the first one.
- **`SetLimit(n)`**: bounds concurrency; `g.Go` blocks when `n` goroutines are already running, without writing your own semaphore.
- **Protecting shared state**: goroutines write into one map, so a mutex is required (maps are not safe for concurrent writes).
- **Loop variable capture**: the closure uses `url` directly, which is safe since Go 1.22 gives each iteration its own variable (the module declares `go 1.24.0`).

## How to run

Own `go.mod` (`module errgroup`). Needs internet access.

```bash
cd errgroup
go run .
```

## Notes

- A plain `errgroup.Group` (not `errgroup.WithContext`) does not cancel the other goroutines when one fails; all four fetches still run to completion.
- `content[:100]` panics if any body is shorter than 100 bytes; this path is only reached when all fetches succeed, which the invalid URL prevents.
- Output order of the fetch lines depends on network timing. Whether `https://www.codeheim.io` responds was not verified.
