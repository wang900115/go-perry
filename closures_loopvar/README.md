# closures_loopvar

Closures and loop variables in Go.

- A closure captures the **variable**, not a snapshot of its value.
- Closure state (`counter()`) lives on the heap when it outlives the enclosing call.
- Go 1.22+: each iteration of a `for` loop gets its own variable (`for i := ...` and `for k, v := range ...`).
- 3-clause loops copy the previous iteration's value into the new variable before the post statement, so body changes carry over.
- Goroutines that capture loop variables are now safe (data races still need synchronization).
- Method values (`acc.Value`) bind the receiver when evaluated.
- `defer f(x)` evaluates `x` immediately; `defer func(){ ... x ... }()` reads `x` when it runs.

```bash
go run .
go build -gcflags=-m . 2>&1 | grep -E "moved to heap|func literal"
go build -gcflags=-m=2 . 2>&1 | grep capturing
```

See [LESSON.md](LESSON.md) for the full lesson.
