# runtime

A set of small examples of `runtime` package functions: `Stack`, `Gosched`, `GOMAXPROCS`, `Goexit`, `NumGoroutine`, `SetFinalizer` and `GC`.

## What the code does

`main.go` is `package main` but has **no `func main`**, so none of these functions are called and the module does not build as is. The functions:

- `riskFunc()`: dereferences a nil `*int`. A deferred `recover` catches the panic and prints it with the built-in `println` (which shows only the interface's type/data pointers, not the message), along with the current goroutine's stack trace from `runtime.Stack(buf, false)` (1 KB buffer).
- `goschedExample()`: starts `riskFunc` in a goroutine, calls `runtime.Gosched()` to yield, then prints `Program finished.`.
- `goMaxProcesExample()`: sets `GOMAXPROCS(1)`, prints the previous value, starts 5 goroutines that print their id, then calls `Gosched()`.
- `exitExample()`: a goroutine calls `runtime.Goexit()`. Its deferred print still runs, and the line after `Goexit` is never reached.
- `seeGoroutine()`: starts 5 sleeping goroutines and prints `runtime.NumGoroutine()`.
- `fininalizerExample()`: attaches a finalizer to a `*Resource`, drops the reference, calls `runtime.GC()` and then `Gosched()`.

## Concepts you can learn

- **recover + `runtime.Stack`**: capture where a panic happened without crashing the process. A nil dereference becomes a recoverable runtime panic.
- **`Gosched`**: voluntarily yields the processor. It is a hint, not a synchronization primitive.
- **`GOMAXPROCS`**: limits how many OS threads run Go code at the same time, which affects parallelism but not concurrency.
- **`Goexit`**: ends only the calling goroutine and still runs its deferred calls, unlike `os.Exit`.
- **`NumGoroutine`**: a cheap way to spot goroutine leaks.
- **Finalizers**: run after the GC finds an object unreachable, on a separate goroutine, at an unspecified time.

## How to run

The module has its own `go.mod`, but `cd runtime && go build .` fails with `function main is undeclared in the main package`. To try an example, add a `func main()` that calls one of the functions, then `go run .`.

## Notes

- Most examples rely on `Gosched()` to let other goroutines run before returning. That is not guaranteed, so output may be missing or appear out of order. Use `sync.WaitGroup` for deterministic results.
- The finalizer in `fininalizerExample` may not run before the function returns, or at all if the program exits first.
- `goMaxProcesExample` sets GOMAXPROCS to 1 and never restores it.
- The module is named `run`, not `runtime`.
- Function names contain typos (`goMaxProcesExample`, `fininalizerExample`). Comments are in Chinese.
