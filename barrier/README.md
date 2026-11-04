# barrier

Memory-ordering ("happens-before") demos: three ways to make a write in one goroutine safely visible to another, using atomics, a mutex, and a channel.

## What the code does

`main.go` sets a shared `data` variable to 42 in a goroutine and prints it from `main`, using a different synchronization tool in each function:

- `main()`: the goroutine writes `data = 42` then `ready.Store(true)` (`atomic.Bool`); `main` spins on `ready.Load()` and prints `42`. The comments describe the store as a release barrier and the load as an acquire barrier.
- `main2()`: protects `data` with a `sync.Mutex` in both goroutines (unlock = release, lock = acquire).
- `main3()`: passes the value over an unbuffered channel (send = release, receive = acquire). Note it declares a local `data := 42`, so it does not touch the package-level `data`.
- A commented-out version at the bottom uses a plain `bool` flag with no synchronization: the broken case the other versions fix.

Only `main()` runs. `main2` and `main3` are never called; to try them you have to edit the code (for example, rename one to `main`).

## Concepts you can learn

- **Go memory model / happens-before**: without synchronization, another goroutine is not guaranteed to see a write (or may see it out of order). Atomics, mutexes, and channels all create the ordering edges that make it visible.
- **Release/acquire semantics**: `atomic.Bool.Store` publishes earlier writes; a `Load` that observes the value sees them too. This is why `data` can be a plain `int`.
- **Mutex and channel as barriers**: they give the same visibility guarantee, plus mutual exclusion or value hand-off.
- **Why the plain-flag version is wrong**: it is a data race; the compiler may hoist the `ready` read out of the loop (spin forever) or `data` may not be visible.

## How to run

Own `go.mod` (`module barrier`):

```bash
cd barrier
go run .
go run -race .     # race detector (requires cgo / a C compiler on Windows)
```

## Notes

- `main2()` as written does not guarantee it prints 42: if `main` takes the lock before the goroutine does, it prints `0`. The mutex gives visibility, not ordering of who goes first.
- The spin loop in `main()` burns CPU while waiting; it is fine for a demo, not for real code.
