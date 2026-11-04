# logical_clock

Four small, single-process simulations of logical clocks used in distributed systems: Lamport, vector, hybrid logical (HLC) and a simplified interval tree clock (ITC).

## What the code does

Each file is a standalone program with its own `func main` (module `logical_clock`, no dependencies).

- **`lamport_clock.go`**: two clocks A and B. `Tick` increments; `Send` ticks and returns the value; `Receive` takes `max(local, received) + 1`. Output: `A event: 1`, `A send: 2`, `B receive: 3`, `B event: 4`, `B send: 5`, `A receive: 6`.
- **`vector__clock.go`**: two-process vector clocks. `Receive` takes the element-wise max and then ticks its own slot. Output: `A: [1 0]`, `A send: [2 0]`, `B: [2 1]`, `B: [2 2]`.
- **`hybrid_clock.go`**: HLC timestamps `{Wall, Logical}`. `Now` uses the physical time if it moved forward, otherwise bumps the logical counter; `Update` merges a remote timestamp using the max wall time and adjusts the logical part. Two goroutines create local events, then A "sends" to B. Values are nanosecond wall times, so output differs per run; when the clock has not advanced, you see the logical part grow (for example `{... 1}`, `{... 2}`, `{... 3}`).
- **`interval_tree_clock.go`**: a tree of counters plus a `path` that identifies a process's interval. Demonstrates `Event`, `Fork`, `Join`, `Clone` and `Compare` (Equal/Before/After/Concurrent), printing the tree after each step.

## Concepts you can learn

- **Happens-before with Lamport clocks**: if a -> b then L(a) < L(b), but not the reverse; a single counter cannot detect concurrency.
- **Vector clocks**: one counter per process, so comparing vectors can tell causality from concurrency.
- **Hybrid logical clocks**: stay close to physical time (useful for timestamps humans read) while still respecting causality when clocks tie or skew.
- **Interval tree clocks (idea)**: identities can be forked and joined dynamically instead of fixing the number of processes up front.

## How to run

The files cannot be built together (`main redeclared`), so `go run .` and `go build ./...` fail by design. Run one file at a time:

```bash
cd logical_clock
go run lamport_clock.go
go run vector__clock.go
go run hybrid_clock.go
go run interval_tree_clock.go
```

## Notes

- **ITC is not a correct ITC.** `Fork` returns a clock that shares the **same root pointer**, so P1 and P2 mutate one tree; after "independent events" both print identical trees and `Compare` reports `Equal` (not `Concurrent`). Only `Clone` (deep copy) yields a real `Before`/`After` result. Real ITCs also split IDs and events separately; this is a simplified sketch.
- `HLC` has no locking; it is safe here only because each node is used by one goroutine at a time.
- The vector clock assumes a fixed number of processes (`n = 2`).
