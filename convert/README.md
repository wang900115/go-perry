# convert

Type conversion in Go: (1) normal numeric conversions, (2) number to string, (3) string to number, and (4) converting between struct types.

## What the code does

- `01.go` - **Normal type conversion**. `a := 0x1234; uint8(a)` keeps only the low byte (prints `34` in hex), and `b := 1234.56; int(b)` truncates to `1234`. These work because they convert variables; `uint8(0x1234)` or `int(1234.56)` on constants would not compile, since constant conversions must be exact. `float64(c)` with `c := 256` prints `256.000000`.
- `02.go` - **Number to string**. Converts `int 10` with `strconv.Itoa` and with `fmt.Sprintf("%d", ...)`, printing type and value each time.
- `03.go` - **String to number**. `strconv.Atoi("5")` gives `int 5`. `strconv.Atoi("abc")` returns an error, which is printed.
- `04.go` - **Struct conversion**. `type bar foo` defines a new type with the same underlying struct, and `foo(b)` converts a `bar` to a `foo`.

Actual output:

```
# 01.go
34
1234
256.000000
# 02.go (no newlines in the format strings)
int 10string 10string 10
# 03.go
int 5
strconv.Atoi: parsing "abc": invalid syntax
# 04.go
main.bar {1}main.foo {1}
```

## Concepts you can learn

- **Conversions are explicit and can lose data**: narrowing an integer wraps (keeps the low bits), and float to int truncates toward zero. Go never does this implicitly, so every lossy step is visible in the code.
- **`string(n)` is not number-to-text**: use `strconv.Itoa` (fast, int only) or `fmt.Sprintf` (flexible, but slower: it parses the format string and boxes the argument in an interface).
- **Parsing returns an error**: `strconv.Atoi` reports bad input as an error instead of returning 0 silently, and you must check it.
- **Conversion between named types**: two types with identical underlying types can be converted with `T(x)`. This is a compile-time conversion, not a type assertion or type switch.

## How to run

Part of the root module (no `go.mod` here). Each numbered file has its own `func main`, so run them one at a time from the repo root:

```bash
go run convert/01.go
go run convert/02.go
go run convert/03.go
go run convert/04.go
```

`go build ./convert` fails by design.

## Notes

- The previous README labeled example 4 "Struct convert (type switch)". The code does not use a type switch. It is a plain conversion between types that share an underlying type.
- In `03.go`, the error from the first `Atoi` call is assigned but never checked.
