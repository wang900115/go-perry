# nosplit-demo

Try the `//go:nosplit` directive, which tells the compiler not to insert a stack-growth check in a function's prologue.

## What the code does

`main.go` defines two `//go:nosplit` functions:

- `hello()` prints `Hello`.
- `ba()` declares a local `[10000]int` (80 KB) array and prints `buf[0]`, which is `0`.

`main` calls both. Running it on Go 1.27.1 / Windows prints:

```
Hello
0
```

## Concepts you can learn

- **Goroutine stack growth**: goroutines start with a small stack. Normally every function prologue checks whether there is enough room and grows the stack if not.
- **`//go:nosplit`**: removes that check. The runtime uses it for code that must not be preempted or must not grow the stack (signal handlers, scheduler internals). The function must then fit in a small fixed reserve.
- **Linker nosplit check**: the linker rejects call chains of nosplit functions whose frames exceed that reserve, with a `nosplit stack over N byte limit` error.

## How to run

```bash
cd nosplit-demo
go run .
go build -gcflags=-S . 2>&1 | grep -A5 'main.ba STEXT'   # inspect the prologue/frame
```

## Notes

- The program builds and runs without a nosplit error because the compiler optimizes the array away: only the zero value `buf[0]` is read, and `go build -gcflags=-S` shows `main.ba` with a 16-byte frame (`$16-0`). To actually trigger the check, index with a variable and prevent inlining (`//go:noinline`, `buf[idx] = 1`); the linker then fails with `main.ba: nosplit stack over 792 byte limit` (verified on Go 1.27.1).
- `main.o` is a committed compiler output (ar archive). The code does not show how it was produced.
