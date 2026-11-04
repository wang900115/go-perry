# pool

Reusing short-lived buffers with `sync.Pool` and counting how many allocations actually happen.

## What the code does

`main.go` creates a `sync.Pool` whose `New` function allocates a 1 KiB `[]byte`, prints `.`, and increments `allocationCnt`. It then starts 1000 goroutines, one every 10 ms. Each goroutine `Get`s a buffer, prints `-`, holds it for 100 ms, and `Put`s it back. After all finish it prints `Number of allocations: N`.

Because each buffer is held for about 100 ms and a new goroutine starts every 10 ms, only around ten buffers are in use at any moment, so `N` is expected to be far smaller than 1000. The exact number varies by run and is not fixed by the code.

## Concepts you can learn

- **`sync.Pool`** caches temporary objects so hot paths can reuse them instead of allocating, reducing GC pressure.
- **`New` is only called on a pool miss**, which is what the allocation counter makes visible.
- **Pool contents are not guaranteed**: the runtime may drop pooled objects during GC (the comment at the top of the file says so), so a pool is a cache, not a fixed-size object store.

## How to run

This directory has no `go.mod`; it belongs to the root module. From the repo root:

```bash
go run pool/main.go
```

Runs for roughly 10 s.

## Notes

- `allocationCnt++` inside `New` is not synchronized, and `New` can be called from several goroutines at once, so this is a data race (`go run -race pool/main.go` should be able to flag it). The printed count is approximate.
- Storing a `[]byte` (not a pointer) in the pool causes an allocation for the interface conversion on each `Put`; pools usually hold pointers such as `*[]byte` or `*bytes.Buffer`.
- A commented-out block shows the single `Get`/`Put` round trip.
