# 2026-11-21 · Generics: type parameters and constraints

## Why learn this
- The standard library is now built on generics: `slices`, `maps`, `cmp`, `iter`, `sync/atomic.Pointer[T]`, `weak.Pointer[T]`, `unique.Handle[T]`. Earlier topics in this repo (`iter`, `sync_atomic`, `weak_pointer`, `unique`, `slice_internals`) already used them; this lesson explains the rules behind them.
- Real uses: type-safe containers (stacks, sets, caches, LRU), helpers such as `Map`/`Filter`/`MaxOf`, typed wrappers around `sync.Pool` or `atomic.Value`, and repository/service layers that share code across entity types without `interface{}` and type assertions.
- Knowing what constraints mean (type sets, `~`, `comparable`) is what lets you read compile errors like `possibly missing ~ for float64` and avoid run-time panics that generics do *not* prevent.

## Core concepts
1. **Type parameters** are declared in square brackets: `func Map[In, Out any](s []In, f func(In) Out) []Out`. Each parameter has a *constraint*. Generic types work the same way: `type Stack[T any] struct{...}` (https://go.dev/ref/spec#Type_parameter_declarations).
2. **Instantiation** substitutes type arguments: `Map[string, int]`. It happens at compile time. Usually you do not write the arguments; the compiler **infers** them from the function arguments (https://go.dev/ref/spec#Type_inference).
3. **Constraints are interfaces.** An interface now defines a *type set*. `any` = all types. `comparable` = types that support `==`/`!=`. `cmp.Ordered` = types that support `<`. A union `int | float64` = exactly those types. **`~float64`** = every type whose *underlying* type is `float64`, including named types like `type Celsius float64` (https://go.dev/ref/spec#General_interfaces).
4. Inside a generic function you can only use operations that **every type in the type set** supports. `+` works in `Sum[T Number]` because all of `~int | ~int64 | ~float64` support `+`.
5. **Constraints can mix methods and types**: `interface{ ~string; Len() int }`. Interfaces with type elements (unions, `~T`) can only be used as constraints, not as ordinary variable types.
6. **`S ~[]E` pattern**: `func Filter[S ~[]E, E any](s S, ...) S` returns the caller's own named slice type (`main.Words`), not a plain `[]string`. This is how `slices.Clone`, `slices.Sort` etc. are declared.
7. **Generic types and methods**: methods of `Stack[T]` use `T`, but a method **cannot declare its own type parameters** (`func (s *Stack[T]) Map[U any]` is a syntax error). Write a top-level function instead.
8. **Zero value**: `var zero T` is the standard way to return "nothing" of type `T`.
9. **Inference limits**: `Zero[T any]() T` cannot be inferred (T appears only in the result), and a generic function cannot be used as a value without instantiation (`f := Sum[float64]`).
10. **`T` is not a dynamic type.** To branch on the actual type, convert to an interface first: `switch any(v).(type)`.
11. **`comparable` and interfaces**: since Go 1.20, interface types such as `any` satisfy `comparable` (https://go.dev/blog/comparable). So `Index[any]` compiles, but comparing two `any` values holding slices **panics at run time**: `comparing uncomparable type []int`.
12. **Go 1.26: self-referential constraints.** A generic type may now refer to itself in its own type parameter list: `type Lesser[T Lesser[T]] interface{ Less(T) bool }`. Before Go 1.26 this was rejected (verified with go1.25.0: `invalid recursive type: Lesser refers to itself`). See the Go 1.26 release notes (https://go.dev/doc/go1.26).

### How it works (compiler/runtime)
- Go does **not** create one copy of the code per type argument (full monomorphization like C++), and it does **not** box everything (like Java erasure). The gc compiler uses **GC shape stenciling with dictionaries** (https://github.com/golang/proposal/blob/master/design/generics-implementation-dictionaries-go1.18.md):
  - Type arguments with the same *GC shape* (same underlying memory layout from the garbage collector's point of view) share one compiled copy. **All pointer types share one shape**, so `Stack[*User]` and `Stack[*Order]` share code.
  - Each instantiation passes a hidden **dictionary** argument with type information (runtime type descriptors, method addresses for the constraint's methods, etc.).
  - For example `Sum[int]` and `Sum[float64]` have different shapes → separate code; `Sum[float64]` and `Sum[Celsius]` have the same underlying type, so they most likely share one copy.
- Consequence: calling a method through a constraint on a pointer-shaped `T` goes through the dictionary (an indirect call), similar in cost to an interface call. Generic code is usually fast, but it is not guaranteed to be as fast as hand-written code for one type. This is an implementation detail of the gc compiler and may change between versions; measure with `benchmark` before optimizing.
- You can see the stenciled functions with `go build -gcflags=-S` (names contain `[go.shape.int]` etc.). I am not certain the exact symbol format is stable across versions.
- Type checking of constraints happens fully at compile time; the only run-time failures come from operations the type set allows but a dynamic value does not (Part 5).

## Full example
See `main.go`. Verified with `go run .` (go1.26.3 linux/amd64), `go vet .` and `GOOS=windows go vet .`. Each commented-out line was uncommented on its own to confirm the compile error. Expected output is in [README.md](README.md).

Parts:
1. `Map` with two type parameters, inference vs explicit instantiation, `Filter[S ~[]E, E any]` keeping `main.Words`.
2. `Number` with `~`, `Sum` on `Celsius`, `StrictNumber` without `~` rejecting `Celsius`, `cmp.Ordered`, `comparable`, a methods-plus-types constraint.
3. `Stack[T]` with `var zero T`, `Pair[K, V]` with a `String` method and `slices.SortFunc`.
4. `Zero[T]()` needs explicit type arguments, `f := Sum[float64]`, type switch via `any(v)`.
5. `Index[any]` works for comparable dynamic values and panics for `[]int`.
6. `Lesser[T Lesser[T]]` and `MinOf` (Go 1.26).

## Common pitfalls
1. **Forgetting `~`.** `int | float64` does not include `type Celsius float64` or `time.Duration` (`~int64`). The compiler tells you: `possibly missing ~ for float64`. Use `~` in reusable constraints unless you really mean "exactly this type".
2. **Believing `comparable` makes `==` safe.** When `T` is `any` or another interface, `==` can panic at run time if the dynamic values are slices, maps or functions. Generic sets/caches keyed by `any` have this risk.
3. **Over-using generics / fighting the language.** Methods cannot have type parameters, you cannot switch on `T` directly, and `T` cannot be inferred from a return type. If a function only needs behaviour, a plain interface parameter (`io.Reader`) is often simpler. The Go team's guideline: write generic code when you would otherwise write the same code for several types (https://go.dev/blog/when-generics).

## Exercise (no answer given)
Build a generic set in a new file, standard library only:
- `type Set[T comparable] struct { m map[T]struct{} }` with `NewSet[T comparable](items ...T) *Set[T]`, `Add`, `Has`, `Len`.
- Add `All() iter.Seq[T]` so that `for v := range s.All()` works (see the `iter` topic).
- Write a **top-level** function `Union[T comparable](a, b *Set[T]) *Set[T]`. Then try to write it as a method that maps to a different element type, `func (s *Set[T]) Map[U comparable](f func(T) U) *Set[U]`, read the error, and write it as a function instead.
- Write `Sorted[T cmp.Ordered](s *Set[T]) []T` using `slices.Sorted` and `maps.Keys`.
- Create a `Set[any]` and add `[]int{1}`. What happens and when (compile time or run time)? Explain using the concept of `comparable` in this lesson.
- Bonus: add a benchmark comparing `Set[int]` with a hand-written `IntSet` (see the `benchmark` topic) and explain the result using GC shapes and dictionaries.

## Suggested directory
`generics_constraints`
