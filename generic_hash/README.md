# generic_hash

Custom hashing and equality for a struct type through the generic `maphash.Hasher[T]` interface from `hash/maphash`. Two values count as equal when only their `Major.Minor` matches.

## What the code does

`main.go`:

- `Version{Major, Minor, Patch}` is the value type.
- `releaseLine` implements `maphash.Hasher[Version]`:
  - `Hash(h *maphash.Hash, v)` writes only `"<Major>.<Minor>"` into the hash.
  - `Equal(a, b)` compares only `Major` and `Minor`.
- `main` compares `1.27.0` with `1.27.4` using `Equal`. It then hashes both with two `maphash.Hash` values that share one seed and compares the `Sum64()` results.

Actual output:

```
true
true
```

## Concepts you can learn

- **Hash/Equal consistency**: if `Equal(a, b)` is true, `Hash` must produce the same hash for both. Here both ignore `Patch`, so the contract holds. Breaking it makes hash-based lookups miss keys.
- **Domain-specific equivalence**: a custom hasher can group values by a key you define ("release line"), without changing the type's own `==`.
- **Seeds**: `maphash` output depends on the seed. Two hashes are only comparable when they use the same `Seed`, so the code calls `SetSeed(seed)` on both. A random seed per process also makes hashes unpredictable across runs.
- **Writing to `maphash.Hash` as an `io.Writer`**: `fmt.Fprintf(h, ...)` feeds the bytes in.

## How to run

Own `go.mod` (`go 1.27.1`, standard library only):

```bash
cd generic_hash
go run .
```

Verified with go1.27.1. `maphash.Hasher` is a newer standard-library API, so older Go toolchains may not have it.

## Notes

- The example only calls the hasher directly. It does not plug it into a container (such as a generic hash map or set), so it shows the interface contract, not a full use case.
- `fmt.Fprintf` allocates and formats text. A real hasher would usually write the integers directly (for example with `maphash.WriteComparable` or binary writes), but the code here does not.
