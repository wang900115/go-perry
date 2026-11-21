# generics_constraints

Go generics from the ground up: type parameters on functions and types, how
type arguments are inferred, and constraints as interfaces that describe *type
sets*. Covers `any`, `comparable`, `cmp.Ordered`, unions with `~` (underlying
types), constraints that mix methods and types, the `S ~[]E` pattern that keeps
a caller's named slice type, generic types and their methods (which cannot add
type parameters), the limits of inference (`Zero[T]()`, function values), type
switches via `any(v)`, the run-time panic when `comparable` is satisfied by an
interface type, and Go 1.26 self-referential constraints (`Lesser[T Lesser[T]]`).

```bash
go run .
```

Standard library only. Requires Go 1.26 (Part 6). Verified with go1.26.3
linux/amd64; `go vet .` and `GOOS=windows go vet .` pass. The commented-out lines
in `main.go` are real compile errors (checked by uncommenting them one by one).
See [LESSON.md](LESSON.md) for the full lesson.

## Expected output

```text
== Part 1: generic functions and inference ==
lengths     = [2 8 3 3]
upper       = [GO GENERICS ARE FUN]
long        = [generics are fun] (type main.Words)

== Part 2: constraints ==
Sum ints    = 6
Sum floats  = 3.75
Sum Celsius = 42 (type main.Celsius)
SumStrict   = 9
MaxOf       = 9 zoo
Index       = 2
Longest     = perry

== Part 3: generic types ==
pop         = 2 true
pop         = 1 true
pop empty   = 0 false
sorted pairs= [a=1 b=2 c=3]

== Part 4: inference limits ==
Zero[string]= "", Zero[*int] = <nil>
f is func([]float64) float64, f(...) = 1
int 42 | string "hi" | other float64

== Part 5: comparable and runtime panics ==
[]any ok    = 1 <nil>
[]any slice = -1 recovered: runtime error: comparing uncomparable type []int

== Part 6: self-referential constraint (Go 1.26) ==
MinOf       = {1 9} <nil>
MinOf empty = empty slice
```
