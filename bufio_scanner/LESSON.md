# 2026-11-17 · bufio: Scanner, Reader, Writer

## Why learn this
- Reading input line by line is everywhere: log files, CSV-like data, stdin for CLI tools, line-based protocols (Redis RESP, SMTP, the HTTP/1.1 header block). `bufio.Scanner` is the standard way to do it.
- `bufio.Writer` turns many small writes into a few large ones. Each `Write` on an `*os.File` or `net.Conn` is a system call, so buffering matters for CLI output, file export, and network protocols.
- `bufio` is also the source of classic production bugs: the 64 KiB line limit (`token too long`), data lost because `Flush` was never called, and slices reused after the next `Scan`.
- It builds directly on the `io_reader_writer` lesson: a `Scanner`, `Reader` and `Writer` all wrap an `io.Reader` / `io.Writer`.

## Core concepts
1. **Scanner loop**: `for sc.Scan() { sc.Text() }`, then check `sc.Err()`. `Scan` returns `false` both at EOF and on an error; `Err` returns `nil` for a clean EOF (https://pkg.go.dev/bufio#Scanner).
2. **Split functions** decide what a token is: `ScanLines` (default), `ScanWords`, `ScanRunes`, `ScanBytes`, or your own `SplitFunc`.
3. **`ScanLines` details**: it removes the trailing `\n` and one `\r` before it (`dropCR`), so Windows `\r\n` files work. An empty line is a token (`""`). The last line is returned even without a final newline.
4. **`SplitFunc` contract**: `func(data []byte, atEOF bool) (advance int, token []byte, err error)`. Return `0, nil, nil` to ask for more data. Return `advance > 0` to consume bytes. A non-nil `token` is delivered to the caller. Return `bufio.ErrFinalToken` to stop scanning without an error (https://pkg.go.dev/bufio#SplitFunc).
5. **Token size limit**: the buffer starts at 4 KiB and doubles up to `bufio.MaxScanTokenSize` (64 KiB). A token that does not fit stops scanning with `bufio.ErrTooLong`. Call `sc.Buffer(buf, max)` before the first `Scan` to raise the limit.
6. **`Bytes` aliases the buffer**: `sc.Bytes()` points into the Scanner's internal buffer and can be overwritten by the next `Scan`. `sc.Text()` allocates a new string, so it is safe to keep.
7. **`bufio.Reader`**: `ReadString('\n')` / `ReadBytes` return data up to and including the delimiter. At the end, they return the remaining data together with `io.EOF` (same rule as `io.Reader`). `Peek(n)` looks ahead without consuming. `ReadLine` is low level; prefer `ReadString` or `Scanner` (https://pkg.go.dev/bufio#Reader.ReadLine).
8. **`bufio.Writer`**: writes go into a buffer (4096 bytes by default) and reach the underlying writer when the buffer is full or on `Flush`. You must call `Flush` at the end, and check its error.

### How it works (standard library)
Based on the go1.26 source (`src/bufio/scan.go`, `src/bufio/bufio.go`); internal details may change.
- `Scanner.Scan` loops: call the split function on `buf[start:end]`; if it returns a token, return `true`. Otherwise, move the unread bytes to the front of the buffer (when the buffer is full or more than half is consumed), double the buffer if it is still full (up to `maxTokenSize`, else `ErrTooLong`), then `Read` more data. That is why Part 2 reassembles words from a source that returns only 3 bytes per `Read`, and why Part 4 shows `Bytes()` slices being overwritten: the data was moved.
- When the reader returns an error (including `io.EOF`), the split function is called with `atEOF=true` so it can return the last partial token.
- If a split function keeps returning empty tokens at EOF without advancing, `Scan` panics after 100 tries (`maxConsecutiveEmptyReads`) to avoid an infinite loop.
- `Scanner.Buffer` panics if called after `Scan`.
- `Writer.Write`: if `p` fits in the free space, it is copied into the buffer. If not and the buffer is empty, `p` is written directly to the underlying writer (no copy). Otherwise it fills the buffer, flushes, and repeats. Part 6 shows both paths.
- `Writer` remembers the first write error; every later `Write` and `Flush` returns it.

## Full example
See `main.go`. Verified with `go run .` (2 runs, identical output) and `go vet .` on go1.26 linux/amd64; `GOOS=windows go vet .` also passes. Expected output is in [README.md](README.md).

Parts:
1. `ScanLines` on `\r\n` input with an empty line and no final newline.
2. `ScanWords` over a reader that returns at most 3 bytes per `Read`.
3. A comma-splitting `SplitFunc` with `ErrFinalToken`.
4. `ErrTooLong` with a 70 KiB line, `Scanner.Buffer` to fix it, and `Bytes()` aliasing vs `bytes.Clone`.
5. `Reader.Peek` and `ReadString` returning the last line with `io.EOF`.
6. `Writer` buffering, `Flush`, the large-write bypass, and a forgotten `Flush`.

## Common pitfalls
1. **Not checking `sc.Err()`.** When a line is longer than 64 KiB, the loop just ends early. Without `if err := sc.Err(); err != nil`, the program silently processes half a file. Use `sc.Buffer` for long lines (minified JSON, base64 blobs), or use `bufio.Reader.ReadString` when lines have no bound.
2. **Forgetting `Flush` (or ignoring its error).** `defer w.Flush()` drops the error; disk-full or a closed connection goes unnoticed. Call `if err := w.Flush(); err != nil { ... }` explicitly before closing the file.
3. **Keeping `sc.Bytes()` across iterations.** Appending `sc.Bytes()` to a slice, or using it as a map key via an unsafe conversion, keeps a view of a buffer that the next `Scan` rewrites. Use `sc.Text()` or `bytes.Clone(sc.Bytes())`.

## Exercise (no answer given)
Write a small `wc`-like tool that reads stdin (`bufio.NewScanner(os.Stdin)`) and prints line, word and byte counts, plus the longest line:
- Make it correct for Windows files (`\r\n`) and for a file whose last line has no newline. Compare your byte count with the real file size and explain any difference.
- Write a custom `SplitFunc` that splits on blank lines (paragraphs), and count paragraphs.
- Feed it a 200 KiB single-line file. First observe the failure, then fix it in two different ways (`Scanner.Buffer` and `bufio.Reader`). Which one handles a 2 GiB line without extra memory planning?
- Write the report with a `bufio.Writer` wrapped around `os.Stdout`, and use a `countingWriter` (like in `main.go`) to show how many real writes happen with and without buffering.

## Suggested directory
`bufio_scanner`
