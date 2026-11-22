# Closures and Loop Variables

> Go version used: `go1.26.3`. The module's `go` line is `go 1.26.3`, which matters (see Pitfall 1).

## Why learn this

Closures are everywhere in real Go code:

- HTTP middleware and handlers (`func(next http.Handler) http.Handler { return http.HandlerFunc(func(w, r) {...}) }`)
- Goroutines started in a loop (`for _, job := range jobs { go func() { process(job) }() }`)
- `sort.Slice`, `slices.SortFunc`, `sync.OnceFunc`, `errgroup.Go`, `defer func() { ... }()`
- Functional options, generators, memoization, iterators (`iter.Seq`)

The classic bug "all goroutines see the last loop value" was so common that Go 1.22 changed the language.
You need to know both the new rule and the old one, because old code and old `go.mod` files still exist.

## Core concepts

### 1. A closure captures variables, not values

A function literal can refer to variables of the surrounding function. The spec says these variables are
*shared* between the surrounding function and the closure, and they survive as long as they are accessible.
([spec: Function literals](https://go.dev/ref/spec#Function_literals))

```go
x := 1
show := func() { fmt.Println(x) }
x = 2
show() // 2
```

### 2. How the runtime implements it

- A closure value is a pointer to a small struct: the function's code pointer followed by the captured variables.
  When a closure is called, the compiler passes this struct in a dedicated "closure context" register
  (`DX` on amd64), so the function body can reach its captured variables.
- If a captured variable may outlive the stack frame (the closure escapes, e.g. is returned), the compiler
  **moves the variable to the heap** (escape analysis). Check with `go build -gcflags=-m`:

  ```
  ./main.go:16:2: moved to heap: n
  ./main.go:17:9: func literal escapes to heap
  ```
- Optimization: if a captured variable is never reassigned after capture and is small, the compiler
  **captures it by value** (copies it into the closure struct) because nobody can observe the difference.
  `go build -gcflags=-m=2` shows `capturing by value: i (addr=false assign=false width=8)`.
  The exact threshold (I believe 128 bytes) is an implementation detail and may change.

### 3. Go 1.22: per-iteration loop variables

Before Go 1.22, `for i := 0; i < n; i++` declared **one** `i` for the whole loop. Since Go 1.22, each iteration
has **its own** `i`. The same applies to `for k, v := range ...`.
([Go 1.22 release notes](https://go.dev/doc/go1.22#language),
[blog: Fixing For Loops in Go 1.22](https://go.dev/blog/loopvar-preview),
[spec: For statements with for clause](https://go.dev/ref/spec#For_clause))

For a 3-clause loop the spec defines the copy precisely: each iteration has its own separate declared variable;
the variable used by the first iteration is declared by the init statement, and the variable used by each
subsequent iteration is declared implicitly before executing the post statement and initialized to the value of
the previous iteration's variable at that moment. So `i++` inside the body carries over to the next iteration
(example section 4).

The rule is applied **per module/file based on the `go` line in `go.mod`** (or a `//go:build go1.xx` line),
not based on the toolchain you install.

### 4. Method values and `defer` evaluate early

- `f := acc.Value` evaluates the receiver **now**. A value receiver copies `acc`; a pointer receiver takes `&acc`.
  ([spec: Method values](https://go.dev/ref/spec#Method_values))
- `defer fmt.Println(x)` evaluates `x` **now**; `defer func() { fmt.Println(x) }()` reads `x` **later**.
  ([spec: Defer statements](https://go.dev/ref/spec#Defer_statements))

## Full example

`main.go` in this directory. Run with `go run .`

Expected output:

```
== 1. capture by reference ==
  x = 2
  c1: 1 2 3
  c2: 1
== 2. per-iteration loop variables (Go 1.22+) ==
  0  1  2
  distinct addresses: true
== 3. one shared variable (what Go <= 1.21 did) ==
  3  3  3
== 4. body changes carry over to the next iteration ==
  1  3  5
== 5. goroutines in a loop ==
  seen: [a b c]
== 6. method values bind the receiver early ==
  value receiver  : 100
  pointer receiver: 999
== 7. defer argument vs closure ==
  deferred closure: 2
  deferred arg    : 1
```

Section 5 uses `sync.WaitGroup.Go`, added in Go 1.25 ([pkg.go.dev/sync#WaitGroup.Go](https://pkg.go.dev/sync#WaitGroup.Go)).

## Common pitfalls

1. **The `go` line in `go.mod` decides the semantics.** A module with `go 1.21` compiled by Go 1.26 still has the
   old shared-variable behavior. Verified: the same 3-closure loop prints `3 3 3` with `go 1.21` and `0 1 2` with
   `go 1.22+`. Bumping the `go` line can change behavior (usually fixing bugs, rarely breaking tests that relied on
   the old behavior). `go vet`'s `loopclosure` check only catches some cases (goroutine/defer as the last statement).
2. **Per-iteration variables do not fix data races.** Each goroutine now gets its own `name`, but appending to a
   shared slice still needs a mutex (or a channel). Run `go test -race` / `go run -race` to check.
3. **Closures keep captured data alive.** A long-lived closure (callback, cache, handler) that captures a large
   slice or struct keeps it from being garbage collected, and captured variables move to the heap, adding
   allocations in hot paths. Capture only what you need (copy the small field into a local first).

## Exercise

Write `memoize(f func(int) int) func(int) int` that returns a closure caching results in a `map[int]int`.

1. Use it to compute `fib(40)` quickly (the memoized function must call itself through the cache, think about how
   a closure can refer to itself).
2. Make it safe for concurrent use from 8 goroutines started in a `for range 8` loop, and verify with `go run -race`.
3. Bonus: use `go build -gcflags=-m` to find which variables move to the heap and explain why.

## References

- https://go.dev/ref/spec#Function_literals
- https://go.dev/ref/spec#For_clause
- https://go.dev/doc/go1.22#language
- https://go.dev/blog/loopvar-preview
- https://go.dev/wiki/LoopvarExperiment
- https://go.dev/ref/spec#Method_values
- https://go.dev/ref/spec#Defer_statements
