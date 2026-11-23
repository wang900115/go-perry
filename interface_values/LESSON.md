# Interface Values

> Go version used: `go1.26.3`. Runs on Windows, macOS and Linux with the standard toolchain (amd64/arm64; the `unsafe` peek assumes the gc compiler).

## Why learn this

Interfaces are the backbone of Go APIs: `error`, `io.Reader`/`io.Writer`, `http.Handler`, `fmt.Stringer`, `sort.Interface`, `any` in `json.Unmarshal`.
Knowing what an interface value *is* explains a set of real bugs and costs:

- `if err != nil` is true even though you "returned nil" (typed nil).
- `map[any]T` or `==` on interfaces panics at run time.
- Passing small structs as `any` (logging, `fmt.Println`, `context.WithValue`) allocates.
- `io.Copy` gets faster when your type implements `io.WriterTo` (an interface-to-interface assertion).

## Core concepts

### 1. An interface value is two machine words

An interface value holds a **dynamic type** and a **dynamic value** ([spec: Interface types](https://go.dev/ref/spec#Interface_types), [spec: Variables](https://go.dev/ref/spec#Variables)).
In the gc runtime (`runtime/runtime2.go`):

- `any` is an `eface`: `{_type *_type, data unsafe.Pointer}`.
- A non-empty interface (e.g. `Shape`) is an `iface`: `{tab *itab, data unsafe.Pointer}`. The `itab` holds the dynamic type plus a table of function pointers for the interface's methods, so a method call is one indirect call.
- itabs are built lazily at run time and cached in a global hash table (`runtime/iface.go`), so the second conversion of the same (interface, type) pair is cheap. Compile-time-known pairs are emitted statically by the compiler.

The example prints `unsafe.Sizeof(any(nil)) == 16` on 64-bit platforms and peeks at the two words with `unsafe`. This layout is an implementation detail: never depend on it in real code.

### 2. Boxing: storing a value in an interface copies it

Assigning `r` to `var s Shape = r` copies `r`. Later changes to `r` are not visible through `s`.
The data word must be a pointer, so:

- **Pointer-shaped** types (`*T`, `map`, `chan`, `func`, single-pointer structs) are stored *directly* in the data word: no allocation.
- Other values are copied to the heap by `runtime.convT*` helpers: usually **one allocation**.
- Optimizations (implementation details, verified on go1.26.3 by the example): integers `< 256` point into a static table `runtime.staticuint64s`; zero-size values point at `runtime.zerobase`; constants can be boxed into read-only data at compile time.

### 3. Type assertions

- `v, ok := x.(T)` never panics; on failure `v` is the zero value of `T`.
- `v := x.(T)` panics with `*runtime.TypeAssertionError` on failure.
- If `T` is a concrete type, the check is a single comparison of the type word.
- If `T` is an interface, the runtime checks the dynamic type's method set (cached via itabs). This is how `io.Copy` discovers `WriterTo`/`ReaderFrom`.

([spec: Type assertions](https://go.dev/ref/spec#Type_assertions))

### 4. Type switches

```go
switch x := v.(type) {
case int, int64: // more than one type: x keeps the type of v (any)
case string:     // exactly one type: x is a string
case nil:        // v is a nil interface
}
```

Cases are tried in order, so put specific interfaces before broader ones (`Shape` before `error`, etc.).
([spec: Type switches](https://go.dev/ref/spec#Type_switches))

### 5. Comparing interfaces

Two interface values are equal if both are nil, or their dynamic types are identical **and** their dynamic values are equal.
`any(1) != any(int64(1))` because the types differ. If the dynamic type is not comparable (slice, map, func), `==` compiles but **panics at run time**; using it as a map key panics with "hash of unhashable type".
([spec: Comparison operators](https://go.dev/ref/spec#Comparison_operators))

### 6. Typed nil

An interface is nil only when **both** words are nil. Returning a nil `*MyErr` through an `error` result produces `{type=*MyErr, data=nil}`, which is not nil.
([FAQ: Why is my nil error value not equal to nil?](https://go.dev/doc/faq#nil_error))

### 7. Static satisfaction checks

`var _ Shape = (*Circle)(nil)` makes the build fail if `*Circle` stops implementing `Shape`. Because `Circle.Area` has a pointer receiver, only `*Circle` is in the method set, so `var _ Shape = Circle{}` does not compile ([spec: Method sets](https://go.dev/ref/spec#Method_sets)).

## Example

`main.go` demonstrates all of the above. Run it:

```bash
go run .
go build -gcflags=-m . 2>&1 | grep escapes   # see which values escape when boxed
```

Expected output (go1.26.3, amd64):

```
== 1. An interface value is two words ==
sizeof(any)   = 16
sizeof(Shape) = 16
nil any: typ == nil: true data == nil: true
same dynamic type -> same type word: true
different values -> different data words: true
pointer stored directly in data word: true

== 2. Boxing copies the value ==
r.Area() = 300  s.Area() = 6

== 3. Boxing may allocate ==
allocs any(int >= 256): 1
allocs any(int < 256):  0
allocs any(Rect):       1
allocs any(*Circle):    0

== 4. Type assertions ==
x.(string) -> "hello", true
x.(int)    -> 0, false (zero value, no panic)
recovered: interface conversion: interface {} is string, not int
strings.Reader is io.WriterTo: true  io.Writer: false

== 5. Type switch ==
nil interface
integer 42 (x is still any here: int)
integer 5 (x is still any here: int64)
string of len 2
Shape with area 4
Shape with area 3
other type main.Circle
error: boom
other type float64

== 6. Comparing interfaces ==
any(1) == any(1):        true
any(1) == any(int64(1)): false (different dynamic types)
recovered: runtime error: comparing uncomparable type []int
recovered: runtime error: hash of unhashable type []int

== 7. The typed-nil trap ==
err == nil: false  dynamic type: *main.MyErr
type word nil: false  data word nil: true
fixed: err == nil: true
```

The allocation counts in section 3 are implementation details and may differ between Go versions or compilers.

## Common pitfalls

1. **Typed nil errors.** Never return a typed nil pointer through an interface result. Return a literal `nil`. `go vet` does not catch this; `nilness` / staticcheck can catch some cases.
2. **Run-time panics from `==` and map keys.** `any` keys and interface comparisons compile fine even if a caller passes a slice. Use `reflect.TypeOf(v).Comparable()` to check first, or avoid `map[any]`.
3. **Hidden allocations.** `fmt.Println(x)`, `slog.Any`, `ctx.Value` and `[]any` all box their arguments. In hot paths, prefer concrete types, generics, or pointer values (and remember that `x.(T)` on a non-pointer value gives you a copy).

## Exercise

Write a function `func SameValue(a, b any) (equal bool, err error)` that:

- returns `false, nil` when the dynamic types differ,
- returns an error (instead of panicking) when the dynamic type is not comparable,
- treats `[]byte` as a special case and compares the contents with `bytes.Equal`.

Then use `testing.AllocsPerRun` to measure how many allocations `SameValue(Rect{1,2}, Rect{1,2})` performs, and explain the number.

## References

- [The Go Programming Language Specification: Interface types](https://go.dev/ref/spec#Interface_types)
- [Go FAQ: nil error](https://go.dev/doc/faq#nil_error)
- [Russ Cox: Go Data Structures: Interfaces](https://research.swtch.com/interfaces) (2009; the core idea still holds, details changed)
- Runtime source: [`runtime/iface.go`](https://go.dev/src/runtime/iface.go), [`runtime/runtime2.go`](https://go.dev/src/runtime/runtime2.go)
