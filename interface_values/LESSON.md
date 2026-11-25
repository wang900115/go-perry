# Interface Values

> Go version used: `go1.26.3`.

## Why learn this

Interfaces are the main abstraction tool in Go, and you already use them everywhere:

- `error`, `io.Reader`/`io.Writer`, `fmt.Stringer`, `http.Handler`, `sort.Interface`
- `any` parameters in `fmt.Println`, `encoding/json`, `slog` attributes, `context.WithValue`
- Mocks and dependency injection in tests
- Optional capabilities: `io.Copy` checks for `io.WriterTo`/`io.ReaderFrom`, `net/http` checks for `http.Flusher`

Knowing what an interface value *is* at run time explains several real bugs and costs:
why `err != nil` is true for a nil pointer, why `==` can panic, why `fmt.Println(x)` may allocate,
and why wrapping a `ResponseWriter` can silently break streaming.

## Core concepts

### 1. An interface value is two words

The spec says an interface variable holds a *dynamic type* and a *dynamic value*
([spec: Interface types](https://go.dev/ref/spec#Interface_types),
[spec: Variables](https://go.dev/ref/spec#Variables)). The gc runtime stores this as two machine words
(see `src/runtime/runtime2.go`):

```go
type iface struct { // non-empty interface, e.g. io.Reader
    tab  *itab          // (interface type, concrete type, method table)
    data unsafe.Pointer // pointer to the value (or the value itself if it is pointer-shaped)
}
type eface struct { // empty interface: any
    _type *_type
    data  unsafe.Pointer
}
```

- `itab` caches the concrete type and the function pointers for the interface's methods.
  Calling `r.Read(p)` is: load `tab`, load the method pointer, call it with `data` as the receiver.
  This is "dynamic dispatch". The runtime builds an `itab` the first time a (interface, type) pair is needed
  and caches it in a global hash table.
- `any` has no methods, so it only needs the type descriptor.
- The layout is an implementation detail of the gc compiler, not part of the language. The example reads it with
  `unsafe` only to show it; never depend on it.

### 2. Assignment copies the value

`var s Shape = sq` copies `sq` into the interface. Changing `sq` afterwards does not change `s`.
If you store a pointer, the pointer is copied, so both refer to the same object.
You cannot modify the value inside an interface in place: `s.(Square).Side = 1` does not compile.

### 3. Boxing and allocation

The data word is one pointer. If the value is not pointer-shaped, the compiler calls a runtime helper
(`convT64`, `convTstring`, `convT`, ... in `src/runtime/iface.go`) that usually copies the value to the heap.
Exceptions in the current runtime:

- pointers, maps, channels, funcs, and single-pointer structs: stored directly, no allocation
- zero-sized values: point to a shared `zerobase`
- integers 0–255 (single byte-sized values): point into a static table `staticuint64s`
- constants and read-only values: the compiler may point to static data

These are optimizations, not guarantees, and they may change between versions. Use
`testing.AllocsPerRun` or `go build -gcflags=-m` to check real code.

### 4. Method sets decide satisfaction

A type `T` satisfies an interface if its **method set** contains all the methods
([spec: Method sets](https://go.dev/ref/spec#Method_sets)).
The method set of `T` contains value-receiver methods; the method set of `*T` contains both.
So if `Area` has a pointer receiver, `*Circle` is a `Shape` but `Circle` is not.

Use a compile-time assertion so the build breaks when a type stops satisfying an interface:

```go
var _ Shape = (*Circle)(nil)
```

### 5. Type assertions and type switches

- `x.(T)` with a concrete `T`: compares the type word with `T`'s descriptor (cheap pointer compare).
- `x.(I)` with an interface `I`: checks at run time whether the dynamic type has all of `I`'s methods
  (cached `itab` lookup). This is how **optional interfaces** work.
- `v, ok := x.(T)` never panics; `v := x.(T)` panics with `*runtime.TypeAssertionError`.
- In a type switch, a `case` with one type gives `x` that type; a `case` listing several types (`case int, int64`)
  leaves `x` with the switch expression's type.
  ([spec: Type switches](https://go.dev/ref/spec#Type_switches))

### 6. Comparing interfaces

Two interface values are equal if both are nil, or they have identical dynamic types **and** equal dynamic
values ([spec: Comparison operators](https://go.dev/ref/spec#Comparison_operators)).

- `any(1) == any(int64(1))` is `false`: different dynamic types.
- If the dynamic type is not comparable (slice, map, func), `==` compiles but **panics at run time**.
  The same happens when such a value is used as a key of `map[any]V`.

### 7. nil interface vs typed nil

An interface is `nil` only when **both** words are empty. Returning a nil `*Circle` as `Shape`
sets the type word, so the result is not `nil`. (This was also covered in `error_wrapping`.)
Return a literal `nil` on the "nothing" path instead.

## Full example

See [main.go](main.go). Run with `go run .`.

Expected output (gc compiler, Go 1.26.3, amd64):

```
== 1. an interface is two words ==
  sizeof(Shape) = 16  sizeof(any) = 16
  any(42): type word set = true , data word set = true
  *(*int)(data) = 42
== 2. assignment copies the value ==
  original area = 100  interface area = 4
  pointer stored: interface sees the change, area = 12
== 3. boxing and allocation ==
  small int (<256): 0 allocs/op
  big int:          1 allocs/op
  string:           1 allocs/op
  pointer:          0 allocs/op
  struct{}{}:       0 allocs/op
== 4. type assertions and type switches ==
  s.(Square) -> {3} true
  s.(*Circle) ok = false
  s.(fmt.Stringer) ok = false
  describe: nil interface
  describe: integer 5 (x is still any here: int)
  describe: Stringer: 21.5°C
  describe: error: boom
  describe: other []int
  recovered: interface conversion: main.Shape is main.Square, not *main.Circle
== 5. optional interfaces ==
  *strings.Reader implements io.WriterTo: true
  wrapper implements io.WriterTo:       false (wrapping hides methods)
== 6. comparing interface values ==
  any(1) == any(1):         true
  any(1) == any(int64(1)):  false
  map[any] with 1 and "1": len = 2
  recovered: runtime error: comparing uncomparable type []int
== 7. nil interface vs typed nil ==
  zero Shape == nil: true
  find(false) == nil: false (dynamic type *main.Circle)
```

Notes:

- Section 3 numbers come from runtime optimizations and may differ on other versions or compilers.
- The `*(*int)(data)` read in section 1 depends on the gc layout; it is for illustration only.

## Common pitfalls

### Pitfall 1: Typed nil makes a non-nil interface

```go
func find() Shape {
    var c *Circle
    return c // != nil to the caller
}
```

Fix: `return nil` explicitly. Be careful with `error` in particular.

### Pitfall 2: `==` on interfaces can panic

`var x, y any = []int{1}, []int{1}; x == y` compiles and panics. The same applies to `map[any]T` keys
and to structs that contain an interface field. If values can be slices or maps, use `reflect.DeepEqual`
or compare in a type-specific way.

### Pitfall 3: Wrappers hide optional interfaces

If you wrap an `http.ResponseWriter` (for logging, metrics) with a struct that only embeds the
`http.ResponseWriter` interface, the wrapper no longer has `Flush`, so `w.(http.Flusher)` fails and
streaming silently stops working. Go 1.20 added `http.ResponseController`, which looks for an
`Unwrap() http.ResponseWriter` method on wrappers
([net/http.ResponseController](https://pkg.go.dev/net/http#ResponseController)).

## Exercise

Write a small `Logger` interface with one method `Log(msg string)`, and two implementations:
`ConsoleLogger` (value receiver) and `*BufferLogger` (pointer receiver, stores messages in a slice).

1. Add compile-time checks for both. Try `var _ Logger = BufferLogger{}` and read the compiler error.
2. Add an optional interface `Flusher { Flush() error }` that only `*BufferLogger` implements.
   Write `func Close(l Logger) error` that calls `Flush` if the logger supports it.
3. Write a `PrefixLogger` that wraps another `Logger`. Make sure `Close(PrefixLogger{...})` still flushes
   the inner logger (hint: think about how `http.ResponseController` uses `Unwrap`).
4. Use `testing.AllocsPerRun` to measure how many allocations `l.Log("hi")` costs for each implementation,
   and explain the result with what you learned in section 3.

## References

- Spec, Interface types: https://go.dev/ref/spec#Interface_types
- Spec, Method sets: https://go.dev/ref/spec#Method_sets
- Spec, Type assertions: https://go.dev/ref/spec#Type_assertions
- Spec, Type switches: https://go.dev/ref/spec#Type_switches
- Spec, Comparison operators: https://go.dev/ref/spec#Comparison_operators
- FAQ, Why is my nil error value not equal to nil?: https://go.dev/doc/faq#nil_error
- Effective Go, Interfaces: https://go.dev/doc/effective_go#interfaces
- Russ Cox, Go Data Structures: Interfaces (older, but the model still applies): https://research.swtch.com/interfaces
