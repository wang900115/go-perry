# coroutine

Call the Go runtime's internal coroutine primitives (`runtime.newcoro` / `runtime.coroswitch`) through `//go:linkname`, and hand control back and forth between two coroutines by hand.

## What the code does

- `main.go` (package `coroutine`, a library with no `main`):
  - Declares `newcoro` and `coroswitch` without bodies and binds them to the runtime functions with `//go:linkname`.
  - `Create(f)` wraps `newcoro`. The callback receives the `*Coro` and a raw `*coro`, which is the same object stored in `Coro.handle` (the runtime calls `c.f(c)`). The test names it `parent`, but it is not a separate handle.
  - `SwitchTo()` switches to the coroutine. `SwitchRaw(c)` switches to a raw handle and ignores `nil`.
- `main_test.go`: `TestRawCoroutineSwitch` creates a task coroutine and calls `SwitchTo()` three times. The task stores the parent handle and yields back with `SwitchRaw` after each phase. The prints alternate in a fixed order: `[Main Thread] A, B` → `[Coroutine] 1, 2` → `[Main Thread] C, D` → `[Coroutine] 4, 5` → `[Main Thread] E, F` → `[Coroutine] 7` → `[Main Thread] G`. Steps 3 and 6 are never printed.

## Concepts you can learn

- **Runtime coroutines**: `coroswitch` transfers control directly to another goroutine without going through the scheduler's run queue. This is the mechanism behind `iter.Pull`.
- **`//go:linkname` pull**: a body-less function declaration plus `import _ "unsafe"` binds a local name to an unexported runtime symbol.
- **Linkname restrictions**: since Go 1.23 the linker rejects references to runtime symbols that are not explicitly allowed. That is why this test needs `-checklinkname=0`.
- **Symmetric hand-off**: there is one `*coro` per pairing. `coroswitch(c)` parks the caller on `c` and resumes whichever goroutine was parked there, so both sides use the same call on the same handle. There is no separate `yield`/`resume` pair.

## How to run

```bash
cd coroutine
go test -v -ldflags=-checklinkname=0 ./...
```

Verified on Go 1.27.1: the test passes with the flag. Without it, linking fails with `link: coroutine: invalid reference to runtime.newcoro`.

## Notes

- This depends on unexported runtime internals that can change or disappear in any Go release. Do not use it in production; `iter.Pull` is the supported API.
- The test only checks that nothing deadlocks or panics. It has no assertions.
