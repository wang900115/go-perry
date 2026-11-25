# interface_values

What a Go interface value really stores.

- An interface value is two words: a type word (`itab` for non-empty interfaces, `*_type` for `any`) and a data word.
- Assigning to an interface copies the value; store a pointer if you want shared state.
- Boxing may allocate: pointers, zero-sized values and small integers (0-255) do not; larger ints and strings usually do.
- Type assertions (`v.(T)`, `v.(I)`) and type switches read the type word; the single-value form panics on mismatch.
- Optional interfaces (`io.WriterTo`, `http.Flusher`) are discovered by assertion, and wrappers hide them.
- `==` on interfaces compares dynamic type and value, and panics if the dynamic type is not comparable.
- A nil pointer stored in an interface makes a non-nil interface.
- `var _ I = (*T)(nil)` checks interface satisfaction at compile time.

```bash
go run .
go build -gcflags=-m . 2>&1 | grep "escapes to heap"
```

See [LESSON.md](LESSON.md) for the full lesson.
