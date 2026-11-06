# sync_cond

`sync.Cond`: letting goroutines sleep until a shared condition becomes true, then waking one (`Signal`) or all (`Broadcast`) of them.

## What the code does

- **Part 1 – Broadcast**: three runners wait on `cond.Wait()` until `ready` becomes true; a single `Broadcast()` wakes them all.
- **Part 2 – Bounded queue**: `BoundedQueue` (capacity 2) uses one `sync.Mutex` and two `sync.Cond`s sharing it:
  - `notFull` – producers wait here when the queue is full.
  - `notEmpty` – consumers wait here when the queue is empty.
  - `Close()` sets `closed` and `Broadcast()`s so every consumer wakes, sees `closed`, and exits.

Expected output (observed on Linux with go1.26.3, also under `-race`):

```
== Part 1: Broadcast 一次叫醒全部 ==
all started: [runner-1 runner-2 runner-3]
== Part 2: 用兩個 Cond 做有界佇列 ==
consumed 10 items, sum = 55
```

The runner order is not guaranteed, so it is sorted before printing.

## Concepts you can learn

- **`Wait` semantics**: must be called with `c.L` held; it atomically unlocks, parks the goroutine, and re-locks before returning. Internally the runtime's `notifyList` hands out a ticket *before* unlocking, so a `Signal` between unlock and park is not lost.
- **Always wait in a `for` loop**: when `Wait` returns, the condition may no longer hold (another goroutine may have taken the item, or `Broadcast` woke many).
- **Do not copy a `Cond`**: it carries `noCopy`/`copyChecker`; copying panics at runtime and `go vet` (copylocks) reports it.
- **Channel vs Cond**: `Broadcast` ≈ `close(ch)` (but repeatable), `Signal` ≈ sending one value. `Wait` has no timeout or `context` support.

Reference: https://pkg.go.dev/sync#Cond

## How to run

```bash
cd sync_cond
go run .
```
