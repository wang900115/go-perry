# io_reader_writer

The `io.Reader` / `io.Writer` contracts and composition: reading data that arrives
together with `io.EOF`, writing custom readers and writers, `io.TeeReader`,
`io.MultiWriter`, `io.LimitReader`, the `io.WriterTo` fast path inside `io.Copy`,
and `io.Pipe` between goroutines with `CloseWithError`.

```bash
go run .
```

Standard library only. This module declares `go 1.26`.
Verified on go1.26 linux/amd64 (3 runs, identical output); `GOOS=windows go vet .` passes.

## Expected output

```text
== Part 1: the Read contract ==
wrong loop: "hello, gophe"
right loop: "hello, gopher" err=<nil>
io.ReadAll: "hello, gopher" err=<nil> (ReadAll never returns io.EOF)
== Part 2: custom Reader and Writer ==
SMALL INTERFACES COMPOSE
copied 25 bytes, err=<nil>, cw.bytes=25
== Part 3: TeeReader + MultiWriter + LimitReader ==
copied 29 bytes, err=<nil>
sha256 = bce2aeea9e6fc31f...
first 8 bytes via LimitReader: "line one"
== Part 4: io.Copy fast paths ==
with WriterTo:    100000 bytes in 1 Write call(s)
without WriterTo: 100000 bytes in 4 Write call(s)
== Part 5: io.Pipe connects a writer goroutine to a reader ==
got: event 1
got: event 2
got: event 3
scanner err: producer finished early
```

See [LESSON.md](LESSON.md) for the full lesson.
