# unique

The `unique` package (Go 1.23+): interning values with `unique.Make` to save memory on duplicates and to get cheap equality checks.

## What the code does

- `main.go` measures heap growth (`runtime.MemStats.Alloc` after a forced GC) for two cases:
  - **A**: keeps 10,000 separate copies (`strings.Clone`) of the same 37-byte string in a `[]string`.
  - **B**: clones the string the same way but stores `unique.Make(str)` in a `[]unique.Handle[string]`. Every call returns the same canonical handle, and the clones become garbage.

  It prints `A: <n> bytes` and `B: <n> bytes`. B is expected to be much smaller, and it was in a local run; the exact numbers vary by run and platform.
- `main_test.go` benchmarks comparing two equal `structA{Content string}` values directly (`s1 == s2`) against comparing their `unique.Handle`s (`h1 == h2`).

## Concepts you can learn

- **Interning**: `unique.Make` deduplicates equal values into one canonical copy. Many duplicate strings (tags, keys, hostnames) then cost the memory of one.
- **Handle comparison is pointer comparison**: two handles are equal if and only if the values are equal, and comparing them does not look at the contents. This is O(1) instead of O(len).
- **Weak canonical map**: values that are no longer referenced through any handle can be garbage collected, so interning does not leak memory forever.
- **Measuring allocations**: `runtime.GC()` + `ReadMemStats` + `runtime.KeepAlive` give a rough before/after heap delta.

## How to run

```bash
cd unique
go run .
go test -bench . -benchmem ./...
```

## Notes

- `MemStats.Alloc` deltas are approximate. Slice growth from `append` is included in both A and B.
- In the benchmarks, the compiler may simplify comparisons of package-level values. Treat sub-nanosecond results with caution.
- Benchmark names contain a typo (`Comparisong`). The module is named `testunique`.
