# weak_pointer

The `weak` package (Go 1.24+): holding a reference to an object without keeping it alive, combined with `runtime.AddCleanup`.

## What the code does

- `main.go` builds a small `WeakCache` that stores `weak.Pointer[Blob]` values in a map.
  - Blob `a` is still strongly referenced by `main`, so `Get("a")` hits after `runtime.GC()`.
  - Blob `b` is only weakly referenced, so after `runtime.GC()` its `Value()` is `nil` and `Get("b")` misses.
  - After `a = nil` and another GC, `Get("a")` misses too.
- `runtime.AddCleanup` registers a callback per blob that prints when the blob has been collected.
- It also shows that two `weak.Pointer`s made from the same object compare equal with `==`.

Expected output (observed in a Linux run with go1.26.3):

```
weak.Make(a) == weak.Make(a): true
  [cleanup] blob 已被回收: b
Get(a): 命中, a 1048576 bytes
Get(b): 已被 GC 回收 → 需要重新產生
  [cleanup] blob 已被回收: a
Get(a): 強參照消失後，也被回收了
```

The `[cleanup]` lines run on a separate goroutine, so their timing and order are not guaranteed.

## Concepts you can learn

- **Weak vs strong references**: the GC does not follow weak pointers when marking, so an object reachable only through them is collected, and `Value()` then returns `nil`.
- **Comparable weak pointers**: `weak.Make(p) == weak.Make(p)` holds even after the object is gone, so they work as map keys.
- **`runtime.AddCleanup` vs `SetFinalizer`**: a cleanup must not reference the object itself, or the object never becomes unreachable.
- **Not an LRU**: when entries disappear depends on the GC (`GOGC`, `GOMEMLIMIT`), so hit rate is unpredictable.

## How to run

```bash
cd weak_pointer
go run .
```

## References

- https://pkg.go.dev/weak
- https://pkg.go.dev/runtime#AddCleanup
- https://go.dev/blog/cleanups-and-weak
