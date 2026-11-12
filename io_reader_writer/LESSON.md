# 2026-11-12 · io.Reader / io.Writer and composition

## Why learn this
- **Almost every I/O API in Go speaks these two interfaces**: `os.File`, `net.Conn`, `http.Request.Body`, `http.ResponseWriter`, `bytes.Buffer`, `gzip.Reader`, `sha256` hashes, `bufio`, `json.Decoder`/`Encoder`. Learn them once and every library plugs together.
- **Streaming instead of loading everything into memory**: hash an upload while saving it to disk, gzip an HTTP response on the fly, cap a request body with `io.LimitReader`.
- **Testable code**: a function that takes an `io.Reader` can be tested with `strings.NewReader` instead of a real file or socket.
- It builds on earlier lessons: `defer` (closing), `error_wrapping` (`io.EOF` is a sentinel error checked with `errors.Is`), and the concurrency series (`io.Pipe` is a synchronous hand-off between goroutines, like an unbuffered channel).

## Core concepts
1. `io.Reader` has one method, `Read(p []byte) (n int, err error)`. It reads **up to** `len(p)` bytes. It may return `n > 0` **together with** a non-nil error (including `io.EOF`). Callers must process `p[:n]` **before** looking at `err` (https://pkg.go.dev/io#Reader).
2. `io.EOF` means "no more input, normal end". It is a sentinel value, not a failure. Functions like `io.ReadAll` and `io.Copy` treat it as success and return `err == nil` (https://pkg.go.dev/io#ReadAll, https://pkg.go.dev/io#Copy).
3. `io.Writer` has one method, `Write(p []byte) (n int, err error)`. It must return a non-nil error if `n < len(p)`, and it must not keep `p` after returning (https://pkg.go.dev/io#Writer).
4. **Composition**: because the interfaces are tiny, the standard library ships small adapters that wrap them:
   - `io.TeeReader(r, w)`: everything read from `r` is also written to `w` (https://pkg.go.dev/io#TeeReader).
   - `io.MultiWriter(ws...)`: one write fans out to several writers (https://pkg.go.dev/io#MultiWriter).
   - `io.LimitReader(r, n)`: returns EOF after `n` bytes (https://pkg.go.dev/io#LimitReader).
   - `io.Pipe()`: connects code that wants a `Writer` to code that wants a `Reader`, across goroutines (https://pkg.go.dev/io#Pipe).
5. **Optional interfaces**: `io.WriterTo` and `io.ReaderFrom`. `io.Copy` checks for them with a type assertion and uses them instead of its own buffer loop (https://pkg.go.dev/io#WriterTo, https://pkg.go.dev/io#ReaderFrom).

### How it works (standard library)
Based on `src/io/io.go` and `src/io/pipe.go` in go1.26; internal details may change between versions.
- **`io.Copy` / `copyBuffer`**: first `if wt, ok := src.(WriterTo); ok { return wt.WriteTo(dst) }`, then `if rf, ok := dst.(ReaderFrom); ok { return rf.ReadFrom(src) }`. Only if neither exists does it allocate a buffer (default `32 * 1024` bytes, smaller if `src` is a `*LimitedReader` with a smaller limit) and loop Read -> Write. That is why `strings.Reader` (which has `WriteTo`) copies 100,000 bytes in **1** Write call, but the same reader hidden behind a wrapper takes **4** calls (32 KiB x 3 + remainder). On some platforms `*os.File.ReadFrom` and `*net.TCPConn.ReadFrom` use `sendfile`/`splice`, which is how the `zero-copy` topic gets used from plain `io.Copy` (I have not verified this on Windows specifically).
- **Wrapping hides methods**: a struct that embeds only `io.Reader` exposes only `Read`. Any optional interface of the wrapped value is lost. This matters for performance, and also for things like `http.Flusher` on a wrapped `ResponseWriter`.
- **`io.Pipe`** has no internal buffer. `Write` sends the slice on an unbuffered channel and waits for the reader to report how many bytes it copied; a large write is delivered in several reads. `CloseWithError(err)` on the write side makes the reader get `err` instead of `io.EOF` (with `Close()` / `CloseWithError(nil)`, the reader gets `io.EOF`).

## Full example
See `main.go`. Verified with `go run .` (3 runs, identical output) and `go vet` on go1.26 linux/amd64; `GOOS=windows go vet .` also passes. Expected output is in [README.md](README.md).

Parts:
1. A reader that returns its last bytes together with `io.EOF`: the wrong loop drops data, the right loop keeps it, `io.ReadAll` returns `err == nil`.
2. A custom `upperReader` (wraps a Reader) and `countingWriter` (wraps a Writer).
3. `TeeReader` into a SHA-256 hash while copying through `MultiWriter`; `LimitReader`.
4. `io.Copy` with and without the `WriterTo` fast path, counting Write calls.
5. `io.Pipe` with a producer goroutine, `bufio.Writer.Flush`, and `CloseWithError` showing up as `Scanner.Err()`.

## Common pitfalls
1. **Checking `err` before using `n`**: `if err != nil { break }` before `buf[:n]` loses the final chunk when a reader returns data with `io.EOF` (Part 1). Always consume `p[:n]` first. Also, a `Read` that returns `0, nil` is allowed and does not mean EOF.
2. **Treating `io.EOF` as a failure, or comparing wrapped errors with `==`**: `io.ReadAll`/`io.Copy` already return `nil` at EOF. When you do get it yourself, prefer `errors.Is(err, io.EOF)`. Never return a wrapped `io.EOF` from your own `Read` (the docs say callers may test for it with `==`; see https://pkg.go.dev/io#EOF).
3. **Forgetting to flush or close**: `bufio.Writer`, `gzip.Writer` and `io.PipeWriter` hold data or block until you call `Flush`/`Close`. Missing `pw.Close()` makes the reader side of a pipe block forever; a missing `Flush` silently truncates output. Check the error from `Close` on writers.
4. **Wrapping loses optional interfaces**: a logging or counting wrapper around a reader/writer hides `WriterTo`/`ReaderFrom` (Part 4) and also `http.Flusher` etc. Either forward the optional method explicitly or accept the slower path knowingly.

## Exercise (no answer given)
Build a small "upload processor" that reads a fake upload from `strings.NewReader` and, in **one pass** over the data:
- writes it to a file in `os.TempDir()` (clean it up at the end),
- computes its SHA-256,
- counts lines (implement a `lineCounter` that is an `io.Writer`),
- rejects input larger than N bytes: return an error like `upload too large` instead of silently truncating. (Hint: think about why `io.LimitReader(r, N)` alone cannot tell "exactly N bytes" from "more than N bytes".)

Then answer: if you replace the file with a `gzip.Writer` wrapped around the file, what must you call, and in which order, so the file is valid?

## Suggested directory
`io_reader_writer`
