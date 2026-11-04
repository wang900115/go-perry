# cbor

Encoding Go values as CBOR (RFC 8949, a binary JSON-like format) with `github.com/fxamacker/cbor/v2`, focusing on **canonical encoding**: the same logical value always produces the same bytes.

## What the code does

`main.go`:

- Defines `User{ID, Name, Admin}` with `cbor:"..."` tags. It is only used by a **commented-out** `main` that marshals and unmarshals a `User` round-trip.
- The active `main` builds two `map[string]any` values with the same entries (`name`, `age`, `active`) written in different literal orders. It encodes both with `cbor.CanonicalEncOptions().EncMode()` and prints each encoding and its SHA-256 digest.

Actual output (both lines are identical):

```
Canonical CBOR A: a3636167651819646e616d6565506572727966616374697665f5
SHA-256        A:        9c7906a22d876562c4f86cba2d6409a8ec46e2e8caf292918a04c1034b1b3cc5
Canonical CBOR B: a3636167651819646e616d6565506572727966616374697665f5
SHA-256        B:        9c7906a22d876562c4f86cba2d6409a8ec46e2e8caf292918a04c1034b1b3cc5
```

The keys come out as `age`, `name`, `active`: shorter keys first, then bytewise order.

## Concepts you can learn

- **Deterministic serialization**: Go map iteration order is random, so a plain encoder can produce different bytes for equal maps. Canonical mode sorts map keys and uses the shortest integer and length forms, so equal values give equal bytes.
- **Why that matters**: hashing, signing, deduplication, and content addressing all need byte-identical output. Here, the equal SHA-256 digests show that A and B are the same value.
- **Struct tags for CBOR field names** (in the commented-out example), similar to `encoding/json`.

## How to run

Own `go.mod` (module `cbor`):

```bash
cd cbor
go run .
```

To try the struct round-trip, uncomment the first `main` and comment out the canonical one. Only one `main` can be active.

## Notes

- `go.mod` lists `fxamacker/cbor/v2` as `// indirect` although it is imported directly. It still builds.
- The type `User` is unused while the first `main` is commented out.
