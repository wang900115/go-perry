# swisstable

A simplified, generic Swiss Table hash map written for learning. It shows the core ideas behind the map design Go adopted in 1.24, but it is not the runtime's implementation.

## What the code does

`main.go` (package `swiss_table`, a library with no `main` and no tests) implements `Table[K comparable, V any]`:

- **Layout**: a `ctrl []uint8` metadata byte per slot plus a parallel `entries` slice. The capacity is rounded up to a power of two (minimum 8). `ctrlEmpty = 0x80` and `ctrlDeleted = 0xFE` mark free slots.
- **Hashing**: `hashKey` uses a SplitMix64-style `mix64` for integer keys and FNV-1a followed by `mix64` for strings. Other types fall back to `formatValue`.
- **`splitHash`**: H1 = `hash & mask` picks the start slot, and H2 = the top 7 bits form a fingerprint (forced to 1 if it is 0).
- **`Get` / `Set` / `Delete`**: probe in groups of 8 consecutive slots. A slot whose control byte differs from the fingerprint is skipped without comparing keys, and probing stops at the first empty slot.
- **Tombstones**: `Delete` marks a slot `ctrlDeleted` instead of emptying it, so later probe chains are not cut short. `Set` reuses the first tombstone it passes.
- **Resizing**: `Set` doubles capacity when the load factor would exceed 0.75, and `resize` re-inserts only live entries, which also drops tombstones.
- `Len`, `Capacity` and `LoadFactor` are accessors.

## Concepts you can learn

- **Open addressing**: all entries live in flat arrays, with no per-bucket linked lists. This is cache-friendly and needs fewer allocations.
- **Control bytes and fingerprints**: a 1-byte H2 check rejects most non-matching slots without loading or comparing the full key.
- **Group probing**: slots are scanned in fixed-size groups. Real Swiss Tables compare a whole group with SIMD or word tricks; here it is a plain loop.
- **Tombstones**: you cannot simply empty a deleted slot in open addressing without breaking lookups, so deletes leave a marker.
- **Load factor and rehashing**: growing at 0.75 keeps probe sequences short.

## How to run

```bash
cd swisstable
go build ./...
```

There is nothing to `go run`. Import the package or add your own `_test.go` to try it. Comments mention benchmarks, but none exist in the directory.

## Notes

- Key types outside the `hashKey` switch (`uint64`, `uint32`, `uint`, `int`, `int64`, `int32`, `string`) fall back to `formatValue`, which returns `""` for them. All such keys (for example structs, `int8`, `float64`) hash to the same value. Lookups stay correct but degrade to linear scans.
- The load factor counts only live entries. Heavy delete-then-insert workloads can fill the table with tombstones until the next resize.
- Probing walks consecutive groups (`start + probe`), which is effectively linear probing rather than the triangular probing used in production Swiss Tables.
