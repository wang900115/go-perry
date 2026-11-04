# lsm

A minimal Log-Structured Merge (LSM) tree key-value store: writes go to a write-ahead log and an in-memory table, which is flushed to immutable sorted files (SSTables) on disk.

## What the code does

- `mem/mem-table.go` - `MemTable`: a `map[string][]byte` guarded by a `sync.RWMutex`. `Put`/`Get` copy the value bytes; `Flush` swaps in an empty map and returns the old one.
- `wal/wal.go` - `WAL`: an append-only `wal.log`. Each record is `keyLen(uint32) | valLen(uint32) | key | value` (little-endian), and every `Append` calls `fsync`. `ReplayInto` reads the log back into a memtable; `Truncate` deletes and recreates the file.
- `ssl/ssl.go` - SSTable writer/reader. `WriteSSTable` sorts the keys and writes data records, then an index (`count`, then `keyLen | key | offset` per key), then an 8-byte trailer holding the index offset. `LoadSSTable` reads the trailer and index into a map; `Get` seeks to the record offset and reads the value.
- `lsm/lsm.go` - `LSM` ties it together:
  - `OpenLSM` replays the WAL and loads existing `sst-*.db` files, newest first.
  - `Put` appends to the WAL, then writes to the memtable.
  - `Get` checks the memtable, then each SSTable from newest to oldest.
  - `FlushMemToSST` writes the memtable to `sst-<unixnano>.db` and truncates the WAL.
  - `CompactTwoNewest` merges the two newest SSTables into a new one and deletes the inputs.
- `main.go` - wipes and recreates `./data_lsm_example`, puts `apple`/`banana`/`cherry`, flushes, overwrites `banana` with `green`, adds `date`, flushes again, prints five lookups, compacts, and prints them again.

Actual output:

```
apple => red
banana => green
cherry => darkred
date => brown
fig => <nil>
After compaction:
apple => red
banana => yellow
cherry => darkred
date => brown
fig => <nil>
```

Note `banana => yellow` after compaction. This is a bug, explained under Notes.

## Concepts you can learn

- **Write-ahead logging**: the record goes to an fsynced log before the memtable, so the memtable can be rebuilt after a crash (`ReplayInto`). This gives durability at the cost of one sync per write.
- **Memtable + immutable SSTables**: writes are sequential appends and in-memory updates. The cost moves to reads, which may have to check several files.
- **Newest-wins read path**: `Get` searches newest to oldest, so a newer write shadows an older value without rewriting old files.
- **On-disk layout with an index and trailer**: a fixed-size trailer points to the index, and the index maps each key to a byte offset. A lookup needs only a seek, not a scan.
- **Compaction**: merging SSTables reduces how many files a read must check and drops overwritten values. The merge order decides which version survives.

## How to run

Own `go.mod` (module `lsm`, no external dependencies):

```bash
cd lsm
go run .
```

## Notes

- **Compaction bug**: `CompactTwoNewest` copies the newest table (`sa`) into the merged map first, then the older one (`sb`), so older values overwrite newer ones. That is why `banana` reverts to `yellow`.
- `main.go` runs `os.RemoveAll("./data_lsm_example")` relative to the current directory. Running it inside `lsm/` deletes and rewrites the `data_lsm_example/` files that are committed to the repo.
- `Get` re-opens each SSTable and re-reads its whole index on every lookup. There is no index caching, bloom filter, or block cache.
- There is no delete operation (no tombstones), and no automatic flush or compaction threshold. `main.go` triggers both manually.
- `MemTable` and `LSM.mu` are separate locks, so `Put` (WAL append and memtable write) and `FlushMemToSST` are not atomic relative to each other.
