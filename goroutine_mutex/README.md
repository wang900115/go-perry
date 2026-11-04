# goroutine_mutex

A ticket-selling simulation where 2000 goroutines share one counter protected by a `sync.Mutex`. Compare with `goroutine_confinement/`, which avoids the lock by confining the counter to one goroutine.

## What the code does

`main.go` starts with 500 tickets and launches 2000 `buyticket` goroutines. Each one locks a package-level `sync.Mutex`, checks the count, and either decrements it and prints `User N Buy ticket, remaining M` or prints `User N No ticket can buy`, then unlocks. `main` waits on a `sync.WaitGroup`.

Expected result: exactly 500 purchases and 1500 "No ticket can buy" lines, with the remaining count going from 499 down to 0. Which user IDs get tickets is random.

## Concepts you can learn

- **Critical sections**: the check-then-decrement must happen atomically; otherwise two goroutines can both see `tickets > 0` and oversell.
- **`sync.Mutex`** serializes access to shared memory between goroutines.
- **`sync.WaitGroup`** with `defer wg.Done()` to wait for all goroutines.
- **Mutex vs confinement**: same problem as `goroutine_confinement/`, different tool; the mutex version keeps the logic in each goroutine instead of a manager.

## How to run

Own `go.mod` (`module mutex`):

```bash
cd goroutine_mutex
go run .
go run -race .     # remove the Lock/Unlock to see the race detector fire; requires cgo on Windows
```

## Notes

- `fmt.Printf` runs while the lock is held, which lengthens the critical section; fine for a demo, but in real code keep I/O outside the lock.
