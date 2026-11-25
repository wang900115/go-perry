# interface_values

What an interface value really is in Go.

- An interface value is two words: a type word (`*_type` for `any`, `*itab` for non-empty interfaces) and a data word.
- Storing a value in an interface copies it; non-pointer values usually allocate (ints `< 256` and zero-size values do not).
- `x.(T)` panics on mismatch, `v, ok := x.(T)` does not; interface-to-interface assertions check the method set at run time.
- In a type switch case with several types, the variable keeps the interface type.
- Comparing interfaces whose dynamic type is a slice/map/func panics at run time (also as map keys).
- A nil pointer stored in an interface makes a non-nil interface (typed nil).
- `var _ I = (*T)(nil)` checks interface satisfaction at compile time.

```bash
go run .
go build -gcflags=-m . 2>&1 | grep escapes
```

See [LESSON.md](LESSON.md) for the full lesson.
