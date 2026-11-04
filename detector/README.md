# detector

Causality detection with version vectors and dotted version vectors (DVV): deciding whether one replica's update happened before another or whether they are concurrent.

## What the code does

Two standalone programs (module `detector`, no dependencies), each with its own `func main`.

- **`version-vector.go`**: `VersionVector map[NodeID]uint64` with `Increment`, `Update` (element-wise max), `Clone`, `HappensBefore`, `Equal`, `Concurrent` and a sorted `String`. Scenario: A makes two events, B copies A's state and adds one event, A makes a third event independently. Output:
  ```
  A: {A:1}
  A: {A:2}
  B: {A:2, B:1}
  A: {A:3}
  A happens-before B: false
  B happens-before A: false
  A concurrent B: false
  ```
- **`dotted-version.go`**: a `DVV` = causal `Context` (version vector) + one `Dot` (actor, counter) naming a single event. Provides `Next`, `Advance`, `NewDVV`, `Merge`, `Contains` and `HappensBefore`. Scenario creates A1, A2, B1 (from A2's history) and A3, prints each DVV's Context and Dot (A1 `{} (A,1)`; A2 `{A:1} (A,2)`; B1 `{A:2} (B,1)`; A3 `{A:1} (A,2)`), then:
  ```
  A2 -> B1: true
  A2 -> A3: false
  B1 -> A3: false
  ```

## Concepts you can learn

- **Version vectors**: per-replica counters; `a < b` when every entry of `a` is `<=` and at least one is `<`. If neither is less, the updates are concurrent (a conflict to resolve).
- **Merging histories**: element-wise max combines what two replicas have seen.
- **Dots**: a dot identifies exactly one event, which lets DVVs distinguish "the event this value came from" from "everything already seen", avoiding false conflicts that plain version vectors can create.

## How to run

The two files redeclare `VersionVector`, so `go run .` / `go build ./...` fail. Run each file on its own:

```bash
cd detector
go run version-vector.go
go run dotted-version.go
```

## Notes

- **`Concurrent` always returns `false`.** `Equal` is implemented as "neither happens-before the other", which is also true for concurrent vectors, so `Concurrent`'s `!Equal` check cancels it out. That is why the output says `A concurrent B: false` even though `{A:3}` and `{A:2, B:1}` are concurrent.
- In `dotted-version.go`, A3 is built from `a2.Context` (which does not include A2's own dot), so A3 gets the dot `(A,2)`, the same as A2. The printed "A3" is therefore not a distinct third event, and `A2 -> A3: false` follows from that.
