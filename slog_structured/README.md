# slog_structured

Structured logging with the standard library `log/slog` (Go 1.21+): `TextHandler`
vs `JSONHandler`, key-value pairs vs typed `Attr`s, `With`/`WithGroup`, dynamic
levels with `LevelVar`, lazy and redacted values with `LogValuer`, a
context-aware custom `Handler`, and Go 1.26 `slog.NewMultiHandler`.

```bash
go run .
```

Standard library only. This module declares `go 1.26` (`NewMultiHandler` needs Go 1.26).
Verified on go1.26 linux/amd64 (3 runs, identical output); `go vet .` and `GOOS=windows go vet .` pass.
The time attribute is removed with `ReplaceAttr` so the output is deterministic.

## Expected output

```text
== Part 1: key-value pairs vs typed Attrs ==
level=INFO msg="user login" user=perry attempt=3
level=INFO msg="user login" user=perry attempt=3
level=INFO msg="odd args" !BADKEY=user
== Part 2: With and WithGroup (JSON) ==
{"level":"INFO","msg":"handled","service":"api","req":{"method":"GET","path":"/users","status":200}}
{"level":"INFO","msg":"nested","db":{"op":"select","ms":7}}
== Part 3: dynamic level and lazy values ==
LogValue calls while level=INFO:  0
level=DEBUG msg=query detail.rows=42 detail.table=users
LogValue calls while level=DEBUG: 1
Enabled(Debug) = true
== Part 4: LogValuer for redaction ==
level=INFO msg=auth user=perry token=REDACTED
== Part 5: custom Handler reading the context ==
level=INFO msg="charge created" component=billing amount=1200 request_id=req-7f3a
level=INFO msg="no context, no request_id" component=billing
== Part 6: Go 1.26 NewMultiHandler ==
level=INFO msg="only text handler"
level=WARN msg="both handlers" disk_free_pct=4
{"level":"WARN","msg":"both handlers","disk_free_pct":4}
```

See [LESSON.md](LESSON.md) for the full lesson.
