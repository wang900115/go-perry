# sync_atomic

`sync/atomic` typed atomics (`atomic.Int64`, `atomic.Uint32`, `atomic.Pointer[T]`, `atomic.Value`): lock-free counters, CAS loops, copy-on-write config hot-swap, and Go 1.23 `Or`/`And` bit flags.

## Run

```
go run .
go run -race .
```

## What the code does

- **Part 1 – counter**: 50 goroutines × 1000 `Add(1)` on `atomic.Int64` → always exactly 50000.
- **Part 2 – CAS loop**: `MaxTracker.Observe` updates a shared max without a mutex using `Load` + `CompareAndSwap` retry.
- **Part 3 – `atomic.Pointer[Config]`**: readers `Load` an immutable snapshot; the writer builds a new `Config` and `Store`s it (copy-on-write). A CAS with a stale pointer fails.
- **Part 4 – `Or` / `And` (Go 1.23+)**: set/clear flag bits atomically; both return the *old* value.
- **Part 5 – `atomic.Value` pitfall**: storing a different concrete type panics.

Expected output (observed on Linux with go1.26.3, also under `-race`; `GOOS=windows go build` also succeeds):

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

## Concepts you can learn

- Atomics make a **single** read/write/RMW indivisible; they do not make *multiple* operations a transaction (`Load` then `Store` is still a race window → use CAS or a mutex).
- Go's memory model: atomic ops behave as sequentially consistent and create happens-before edges.
- Typed atomics (Go 1.19+) are auto 64-bit aligned even on 32-bit platforms and embed `noCopy`.
- `atomic.Pointer[T]` is type-safe; prefer it over `atomic.Value`.

References: <https://pkg.go.dev/sync/atomic>, <https://go.dev/ref/mem#atomic>, <https://go.dev/doc/go1.19#atomic_types>, <https://go.dev/doc/go1.23#minor_library_changes>
