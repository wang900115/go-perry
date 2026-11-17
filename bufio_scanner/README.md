# bufio_scanner

Buffered I/O with the `bufio` package: `bufio.Scanner` for lines and words
(CRLF handling, a final line without a newline, sources that return a few bytes
per `Read`), a custom `SplitFunc` with `bufio.ErrFinalToken`, the 64 KiB token
limit (`bufio.ErrTooLong`) and `Scanner.Buffer`, `Scanner.Bytes` aliasing,
`bufio.Reader` `Peek`/`ReadString` (data returned together with `io.EOF`), and
`bufio.Writer` buffering, the large-write bypass, and the forgotten `Flush`.

```bash
go run .
```

Standard library only, no files or network. Verified on go1.26 linux/amd64
(2 runs, identical output); `go vet .` and `GOOS=windows go vet .` pass.
See [LESSON.md](LESSON.md) for the full lesson.

## Expected output

```text
== Part 1: Scanner lines (CRLF and missing final newline) ==
line 1: "alpha"
line 2: "beta"
line 3: ""
line 4: "gamma"
lines: 4 (\r is dropped, empty line kept, last line returned)

== Part 2: Scanner words and a 3-byte-per-Read source ==
words: ["the" "quick" "brown" "fox"]
tokens are reassembled even though each Read returned at most 3 bytes

== Part 3: custom SplitFunc + ErrFinalToken ==
field: "red"
field: "green"
field: "blue"
err after ErrFinalToken: <nil>
no END: ["a" "b" "last"] err=<nil>

== Part 4: token too long and Scanner.Buffer ==
default: got line of 5 bytes
default: err = bufio.Scanner: token too long | is ErrTooLong: true
default MaxScanTokenSize = 65536
1MiB max: got line of 5 bytes
1MiB max: got line of 71680 bytes
1MiB max: got line of 5 bytes
1MiB max: err = <nil>
kept without copy: ["six" "two" "six"]
kept with copy:    ["one" "two" "six"]

== Part 5: bufio.Reader ReadString and Peek ==
peek: "HTTP" buffered=36
ReadString: "HTTP/1.1 200 OK\n" err=<nil>
ReadString: "last line no newline" err=EOF

== Part 6: bufio.Writer buffering and Flush ==
before Flush: underlying calls=1, buffered=9, written="row0\nrow1\nrow2\nr"
after Flush:  underlying calls=2, buffered=0, total bytes=25
large write on empty buffer: calls=1 buffered=0
no Flush: underlying got 0 bytes (data still in buffer: 27)
```
