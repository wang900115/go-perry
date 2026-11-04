# escape

Escape analysis: how the Go compiler decides whether a value lives on the stack or the heap, and how inlining changes that decision.

## What the code does

`main.go` defines `NewUser(name string) *User`, which returns `&User{Name: name}`. `main` calls `NewUser("Perry")` and prints `u.Name` (`Perry`). The comments list the three common escape triggers (returning a pointer, storing a value in an interface, capturing it in a closure) and some general allocation-reduction tips.

Output of `go build -gcflags=-m .` on Go 1.27.1:

```
./main.go:14:6: can inline NewUser
./main.go:20:14: inlining call to NewUser
./main.go:21:13: inlining call to fmt.Println
./main.go:14:14: leaking param: name
./main.go:15:7: &User{...} escapes to heap
./main.go:20:14: &User{...} does not escape
./main.go:21:13: ... argument does not escape
./main.go:21:15: u.Name escapes to heap
```

## Concepts you can learn

- **Returning a pointer makes the value escape**: inside `NewUser` on its own, `&User{}` escapes to the heap (line 15) because the pointer outlives the function.
- **Inlining can undo an escape**: after `NewUser` is inlined into `main`, the same allocation `does not escape` (line 20), so it stays on the stack. Small constructors are often free for this reason.
- **Interfaces cause escapes**: `fmt.Println` takes `...any`, so `u.Name` escapes when it is converted to an interface (line 21). This is why the comments advise against `fmt` in hot paths.
- **`leaking param`**: the compiler tracks that the `name` string flows into a returned value.
- **Why it matters**: heap allocations cost GC work and cache locality, while stack allocations are close to free.

## How to run

```bash
cd escape
go run .
go build -gcflags=-m .        # escape and inlining decisions
go build -gcflags="-m -m" .   # more detailed reasons
```

## Notes

- The "HPC core" tips at the bottom of `main.go` are general guidance. This code only exercises the return-pointer and interface cases; closures, `sync.Pool` and preallocation are not demonstrated.
- `escape.exe` is a committed build artifact.
