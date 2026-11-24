# map_semantics

How Go maps behave at the language level.

- A `nil` map can be read, ranged over and `len`-ed, but writing to it panics.
- A map value is a pointer to a runtime structure: assigning or passing a map shares the entries. Use `maps.Clone` for a (shallow) copy.
- Use the comma-ok form `v, ok := m[k]` to tell a missing key from a zero value.
- Iteration order is unspecified and randomized. Use `slices.Sorted(maps.Keys(m))` (Go 1.23) for a deterministic order.
- Map elements are not addressable: `m[k].Field = x` does not compile. Copy out and store back, or store pointers.
- Deleting entries during `range` is safe; entries added during `range` may or may not be visited.
- `NaN` keys can be inserted but never looked up or deleted; `clear(m)` (Go 1.21) removes them.
- `map[K]struct{}` is the set idiom; interface keys with non-comparable dynamic types panic at run time.
- Concurrent map writes are a fatal error that `recover` cannot catch.

```bash
go run .
```

See [LESSON.md](LESSON.md) for the full lesson.
