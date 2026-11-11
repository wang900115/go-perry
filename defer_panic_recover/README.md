# defer_panic_recover

`defer`, `panic` and `recover`: LIFO order, when deferred arguments are evaluated,
turning a panic into an `error` with a named result, the rule that `recover` only
works when called directly by a deferred function, `panic(nil)` since Go 1.21,
`runtime.Goexit`, and recovering panics inside goroutines.

```bash
go run .
```

Requires Go 1.25+ (`sync.WaitGroup.Go`, `range` over int). This module declares `go 1.26`.
Verified on go1.26 linux/amd64 (3 runs, identical output); `GOOS=windows go vet .` passes.

## Expected output

```text
== Part 1: defer order and evaluation ==
function body done, x = 42
loop defer 2
loop defer 1
loop defer 0
deferred closure sees x = 42
deferred with arg x = 1
double() = 42
== Part 2: panic -> error ==
string panic: recovered panic: boom
  errors.Is(err, ErrPanic): true
index: recovered panic: runtime error: index out of range [5] with length 3
  is runtime.Error: true
nil map: recovered panic: assignment to entry in nil map
== Part 3: recover must be called directly by a deferred func ==
  helper: recover() = <nil>
helper: recovered panic: not caught by helper
  Goexit: recover() = <nil>
== Part 4: panic(nil) since Go 1.21 ==
nil: recovered panic: panic called with nil argument
  is *runtime.PanicNilError: true
== Part 5: recovering inside goroutines ==
worker 1: recovered panic: worker 1 failed
main still running: ok
```

See [LESSON.md](LESSON.md) for the full lesson.
