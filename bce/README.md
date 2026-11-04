# bce

Bounds Check Elimination: when the Go compiler can drop the slice bounds check, and when it has to keep it.

## What the code does

`main.go` defines two summing functions and prints both results:

- `sum(a)` loops `for i := 0; i < len(a); i++` and adds `a[i]`. The loop condition proves `i` is in range, so the compiler removes the bounds check.
- `sum2(a, index)` adds `a[i]` for each `i` taken from a second slice. The compiler cannot prove these indexes are valid, so it checks every access.

`main` calls both on `[]int{1, 2, 3, 4, 5}`. The program prints `15` twice; `println` writes to stderr.

## Concepts you can learn

- **Bounds checks**: every `a[i]` normally compiles to a compare-and-panic. In hot loops this adds a branch and can block other optimizations.
- **Provable ranges**: a loop bounded by `len(a)` gives the compiler the proof it needs, so idiomatic loops already avoid the check.
- **Data-dependent indexes**: an index that comes from another slice, a map or input always keeps its check.
- **Checking the compiler's decisions**: the `check_bce` debug flag reports every remaining check, so you can confirm what was removed without reading assembly.

## How to run

```bash
cd bce
go run .
go build -gcflags="-d=ssa/check_bce/debug=1" .   # report remaining bounds checks
```

On Go 1.27.1, the `check_bce` build reports `Found IsInBounds` at `main.go:21` (inside `sum2`) and at `main.go:30` (where `sum2` is inlined into `main`). It reports nothing for `sum`.

## Notes

- `bce.exe` is a build artifact committed to the repo.
