# 2026-11-07 · sync/atomic: typed atomics, CAS loops, and atomic.Pointer

## Why learn this
- **Counters**: QPS, error counts, connection counts — many goroutines increment, reads are occasional. A mutex is overkill; `atomic.Int64` is just right.
- **Config hot-swap**: read-heavy, write-rare data (feature flags, routing tables, certificates). `atomic.Pointer[T]` with copy-on-write lets readers proceed without any lock.
- **State flags**: server ready / draining / closed.
- **Reading the standard library**: the `sync.Once` fast path is an `atomic.Bool`; `sync.Mutex`, `WaitGroup`, and `context` are built on atomics.

This follows `sync_cond` → `sync_once`, going one level deeper in the sync series: what locks are made of.

## Core concepts
1. **Atomicity covers a single operation only**: `Add`, `Load`, `Store`, `Swap`, `CompareAndSwap`. A `Load` followed by a `Store` is two operations, and another goroutine can interleave between them → use a CAS loop or a mutex.
2. **Typed API (Go 1.19+)**: `atomic.Int32/Int64/Uint32/Uint64/Uintptr/Bool/Pointer[T]`.
   - No accidental non-atomic access to the same variable.
   - The zero value is ready to use.
   - Automatically 64-bit aligned, even on 32-bit platforms.
   - Embeds `noCopy`, so `go vet` flags copies.
   - The old function style `atomic.AddInt64(&x, 1)` still works, but on 386/ARM32 `x` must be 64-bit aligned or it panics (https://pkg.go.dev/sync/atomic#pkg-note-BUG).
3. **CAS loop**: `for { old := x.Load(); new := f(old); if x.CompareAndSwap(old, new) { break } }`. A failed CAS means someone else changed the value; reload, recompute, retry. This is the basic building block of lock-free structures.
4. **`atomic.Pointer[T]` + copy-on-write**: never mutate a published object. Build a new one and `Store` the new pointer. A reader's `Load` always sees one complete version; old versions are reclaimed by the GC once unreferenced.
5. **Go 1.23 adds `And` / `Or`** on `Int32/Int64/Uint32/Uint64/Uintptr`, **returning the old value** (https://go.dev/doc/go1.23#minor_library_changes).
6. **Memory model**: Go atomics behave as sequentially consistent. If atomic op B observes the value written by atomic op A, then A is synchronized before B. So "write the data, then atomically Store the pointer" guarantees a reader that Loads the pointer sees the complete data (https://go.dev/ref/mem#atomic).

### How the compiler/runtime implements it
- On common platforms the compiler treats `sync/atomic` calls as **intrinsics** and emits CPU instructions directly instead of a function call. To my knowledge, on amd64:
  - `Add` → `LOCK XADDQ`
  - `CompareAndSwap` → `LOCK CMPXCHGQ`
  - `Load` → plain `MOV` (x86 is TSO, a plain load is sufficient)
  - `Store` → `XCHGQ` (implicitly locked, acts as a full barrier)
- Verify yourself: `go build -gcflags=-S . 2>&1 | grep -E "XADD|CMPXCHG|XCHG"`
- **Uncertain**: `And`/`Or` return the old value, but amd64 `LOCK OR` does not. I believe the compiler emits a CAS loop when the result is used and `LOCK OR` when it is not, but I have not verified this across versions — trust what `-gcflags=-S` shows.
- The `-race` detector understands atomics and does not report them, but it does catch a variable accessed both atomically and non-atomically.

## Full example
See `main.go` in this directory (run with `go run .`, `go run -race .`, and `go vet` on go1.26.3 linux/amd64; `GOOS=windows go build` also succeeds). Requires Go 1.23 (`Or`/`And`; integer `for range` needs 1.22).

```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ---------- 1. Counter: atomic.Int64 ----------

type Stats struct {
	requests atomic.Int64 // zero value is ready to use, no init needed
	errors   atomic.Int64
}

// ---------- 2. CAS loop: lock-free "update only if larger" ----------

type MaxTracker struct {
	max atomic.Int64
}

func (m *MaxTracker) Observe(v int64) {
	for {
		old := m.max.Load()
		if v <= old {
			return // no update needed
		}
		if m.max.CompareAndSwap(old, v) {
			return // new value installed
		}
		// CAS failed: another goroutine changed it after our Load; reload and retry
	}
}

// ---------- 3. atomic.Pointer: config hot-swap (copy-on-write) ----------

type Config struct {
	Version  int
	Features []string
}

var current atomic.Pointer[Config]

func loadConfig() *Config { return current.Load() }

func reload(version int, features ...string) {
	// Build a brand-new Config and swap it in with a single Store; the GC reclaims the old one.
	// Never mutate the object returned by current.Load().
	current.Store(&Config{Version: version, Features: features})
}

// ---------- 4. Flag bits: Go 1.23 Or / And ----------

const (
	FlagReady uint32 = 1 << iota
	FlagDraining
	FlagClosed
)

func main() {
	fmt.Println("== 1. atomic.Int64 counter ==")
	var s Stats
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range 1000 {
				s.requests.Add(1)
				if (i*1000+j)%100 == 0 {
					s.errors.Add(1)
				}
			}
		}()
	}
	wg.Wait()
	fmt.Printf("requests=%d errors=%d\n", s.requests.Load(), s.errors.Load())

	fmt.Println("== 2. CAS loop max tracker ==")
	var mt MaxTracker
	for g := range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range int64(10000) {
				mt.Observe((v*7919 + int64(g)*31) % 99991)
			}
		}()
	}
	wg.Wait()
	fmt.Println("max =", mt.max.Load())

	fmt.Println("== 3. atomic.Pointer config hot-swap ==")
	reload(1, "login")
	cfg := loadConfig() // readers get a snapshot of one complete version
	reload(2, "login", "dark-mode")
	fmt.Printf("old snapshot: v%d %v\n", cfg.Version, cfg.Features)
	fmt.Printf("current:      v%d %v\n", loadConfig().Version, loadConfig().Features)
	swapped := current.CompareAndSwap(cfg, &Config{Version: 99})
	fmt.Println("CAS with stale pointer succeeded?", swapped)

	fmt.Println("== 4. Flag bits Or / And (Go 1.23+) ==")
	var state atomic.Uint32
	old := state.Or(FlagReady | FlagDraining) // returns the OLD value
	fmt.Printf("Or  before=%03b after=%03b\n", old, state.Load())
	old = state.And(^FlagDraining) // clear Draining
	fmt.Printf("And before=%03b after=%03b\n", old, state.Load())

	fmt.Println("== 5. atomic.Value panics on inconsistent type ==")
	var v atomic.Value
	v.Store("hello")
	func() {
		defer func() { fmt.Println("recover:", recover()) }()
		v.Store(42)
	}()
	fmt.Println("value is still:", v.Load())
}
```

Expected output:
```
== 1. atomic.Int64 counter ==
requests=50000 errors=500
== 2. CAS loop max tracker ==
max = 99990
== 3. atomic.Pointer config hot-swap ==
old snapshot: v1 [login]
current:      v2 [login dark-mode]
CAS with stale pointer succeeded? false
== 4. Flag bits Or / And (Go 1.23+) ==
Or  before=000 after=011
And before=011 after=001
== 5. atomic.Value panics on inconsistent type ==
recover: sync/atomic: store of inconsistently typed value into Value
value is still: hello
```
The max in part 2 depends only on the fixed input set, not on goroutine scheduling, so it is always 99990.

## Common pitfalls
1. **"Load then Store" is not atomic**:
   ```go
   if x.Load() < limit { x.Store(x.Load() + 1) } // two goroutines can both pass the check
   ```
   Use a CAS loop, or just a mutex. When several fields must stay consistent together (e.g. `count` and `sum`), atomics cannot help — use a mutex, or bundle them in a struct and swap it with `atomic.Pointer`.
2. **Mutating the object behind an `atomic.Pointer`**: `cfg := current.Load(); cfg.Features = append(cfg.Features, "x")` races with every concurrent reader. Treat published objects as immutable; replace, never modify.
3. **Mixing atomic and plain access, or copying atomics**:
   - Using `atomic.AddInt64` in one place and `n++` or a plain read elsewhere is a data race.
   - Passing a struct containing `atomic.Int64` by value copies its state; `go vet` (copylocks) warns.
   - `atomic.Value` fixes its concrete type on the first `Store`; storing a different type or `nil` panics.

## Exercise (no answer given)
Implement a **lock-free token bucket rate limiter** with atomics:
- `type Limiter struct { tokens atomic.Int64; capacity int64 }`
- `TryAcquire() bool`: if `tokens > 0`, atomically decrement and return `true`; otherwise return `false`.
  - **Must use a CAS loop** — do not `Add(-1)` and then check.
  - Think about it: what goes wrong if you `Add(-1)`, see a negative value, and add it back?
- `Refill(n int64)`: add `n` tokens without exceeding `capacity` (also a CAS loop).
- In `main`:
  1. With capacity = 100, start 1000 goroutines that each call `TryAcquire` once; count successes with an `atomic.Int64` and verify exactly 100.
  2. Call `Refill(30)`, run again, and verify 30 successes.
  3. Confirm no races with `go run -race .`.
- Bonus: write a benchmark with `b.RunParallel` comparing the CAS version to a `sync.Mutex` version.

## Suggested directory name
`sync_atomic`

## References
- https://pkg.go.dev/sync/atomic
- https://go.dev/ref/mem#atomic
- https://go.dev/doc/go1.19#atomic_types
- https://go.dev/doc/go1.23#minor_library_changes
