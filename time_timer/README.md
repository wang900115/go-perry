# time_timer

`time.Timer`, `time.Ticker` and `time.AfterFunc`, focusing on the Go 1.23+ semantics:
timer channels are synchronous (unbuffered), and after `Stop`/`Reset` returns no stale
value can be received.

```bash
go run .
```

Requires Go 1.23+ **and** `go 1.23` or later in the main module's `go.mod`
(this module declares `go 1.26`). Verified on go1.26 linux/amd64; `GOOS=windows go vet` passes.

## Expected output

```text
== Part 1: timeout in a loop ==
event 1
event 2
event 3
idle timeout, leaving loop
== Part 2: Stop/Reset leave no stale value ==
after Reset, waited >= 150ms: true
Stop on expired-but-unreceived timer returns: true
t2.C is empty after Stop
== Part 3: Ticker drops ticks for a slow receiver ==
ticks produced ~15, received < 8: true
== Part 4: AfterFunc ==
AfterFunc timer channel is nil: true
Stop before firing returns: true
Stop after firing returns: false
fired total: 10
== Part 5: timer channel is unbuffered ==
len=0 cap=0
but a receive still gets the value
```

## Compare with the pre-1.23 behaviour

```bash
GODEBUG=asynctimerchan=1 go run .        # Linux/macOS
$env:GODEBUG="asynctimerchan=1"; go run . # Windows PowerShell
```

Part 2 then prints `waited >= 150ms: false`, `Stop ... returns: false` and
`received stale value (old behaviour)`, and Part 5 prints `len=1 cap=1`.

See [LESSON.md](LESSON.md) for the full lesson.
