# time.Time Values: Layouts, Locations, the Monotonic Clock, and Calendar Math

> Go version used: `go1.26.3`. Everything here works on Go 1.20+ (`time.DateTime` / `time.DateOnly` were added in 1.20; `time/tzdata` in 1.15; monotonic readings in 1.9).

## Why learn this

Almost every service touches time, and time bugs are quiet: they compile, pass tests written in UTC, and break at midnight, at month end, or twice a year at DST changes.

- Parsing timestamps from logs, CSV files, HTTP headers, and JSON
- Building daily/monthly reports ("orders for Nov 28, Taipei time")
- Measuring latency and timeouts correctly even if the system clock jumps
- Deduplicating or caching by timestamp (`map[time.Time]...`)
- Scheduling ("same time tomorrow", "end of next month")

`time_timer` in this repo covers timers and tickers; this topic covers the **value** type `time.Time` itself.

Docs: https://pkg.go.dev/time

## Core concepts

### 1. What a `time.Time` contains

```go
type Time struct {
    wall uint64    // wall-clock seconds + nanoseconds, plus a hasMonotonic flag bit
    ext  int64     // wall seconds, or the monotonic reading when hasMonotonic is set
    loc  *Location // time zone used for display; nil means UTC
}
```

Source: https://cs.opensource.google/go/go/+/refs/tags/go1.26.3:src/time/time.go

So one value can carry three things: an **instant**, an optional **monotonic clock reading**, and a **Location**.
The Location only affects how the instant is presented (`Format`, `Hour()`, `Day()`...). Two values in different
locations can be the same instant. Pass `time.Time` by value, not by pointer (the docs recommend this).

### 2. Layouts: the reference time

Go has no `YYYY`/`MM`/`DD` tokens. You write the layout by showing how one specific moment would look:

```
Mon Jan 2 15:04:05 MST 2006     (= 01/02 03:04:05PM '06 -0700)
```

`2006` year, `01` month, `02` day, `15` hour (24h), `03` hour (12h), `04` minute, `05` second, `-0700` offset.
Predefined constants: `time.RFC3339`, `time.DateTime` (`"2006-01-02 15:04:05"`), `time.DateOnly`, `time.TimeOnly`.
Any text that is not a recognized element is copied literally, so a wrong layout does not error; it just prints garbage.

### 3. Parse vs ParseInLocation

- `time.Parse(layout, s)`: if `s` has no zone/offset, the result is **UTC**.
- `time.ParseInLocation(layout, s, loc)`: if `s` has no zone, interpret it as wall-clock time in `loc`.

`"2026-11-28 09:00:00"` read from a Taiwanese log file is 01:00 UTC, not 09:00 UTC. Parse also validates ranges
(`2026-02-30` gives `day out of range`), unlike `time.Date`, which normalizes.

### 4. The monotonic clock

Since Go 1.9, `time.Now()` also records a reading from the OS monotonic clock (a counter that never jumps backwards,
unlike the wall clock, which NTP or a user can change). Rules from https://pkg.go.dev/time#hdr-Monotonic_Clocks:

- `t.Sub(u)`, `time.Since`, `time.Until`, `Before`, `After`, `Equal`, `Compare` use the monotonic readings **if both values have one**.
- `t.Round(0)`, `t.In(loc)`, `t.UTC()`, `t.Local()`, `t.AddDate(...)`, `t.Truncate(d)` / `t.Round(d)` strip it; `t.Add(d)` keeps it.
- `time.Date`, `time.Parse`, `time.Unix` and decoding (JSON, gob...) never have one.
- `t.String()` shows it as `m=+0.000012345`.

So measuring elapsed time with `time.Since(start)` is correct even if the wall clock is moved during the measurement.

### 5. Comparing: `Equal`, not `==`

`==` compares the struct fields: wall bits, ext, and the `*Location` pointer. Two values for the same instant compare
unequal with `==` if they differ in location or in whether they carry a monotonic reading. The same applies to map
keys. Use `t.Equal(u)` (or `t.Compare(u) == 0`, Go 1.20+); for map keys normalize first, e.g. `t.UTC().Round(0)` or use `t.UnixNano()`.

### 6. Calendar arithmetic and DST

- `t.Add(d)` adds an exact **duration** (elapsed time).
- `t.AddDate(y, m, d)` adds **calendar** units in `t`'s location, then normalizes like `time.Date` does:
  October 32 becomes November 1, so Jan 31 + 1 month = Mar 3 (2026 is not a leap year).
- On the day DST starts in New York, `Add(24*time.Hour)` lands at 13:00, while `AddDate(0, 0, 1)` keeps 12:00.
  That day is 23 hours long.
- A wall-clock time inside the DST gap (02:30 on 2026-03-08 in New York) does not exist. `time.Date` does not
  fail: the docs say it returns a time that is correct in one of the two zones involved, but **does not guarantee which**.
  In this run on go1.26.3 it returned 01:30 EST; do not depend on that.

### 7. Where zone data comes from (important on Windows)

`time.LoadLocation("Asia/Taipei")` needs the IANA database. Lookup order (https://pkg.go.dev/time#LoadLocation):
the `ZONEINFO` env var, OS locations (on Unix `/usr/share/zoneinfo`), `$GOROOT/lib/time/zoneinfo.zip`, and finally
the embedded copy from `time/tzdata`. Windows has no IANA files, so a binary copied to a machine without Go installed
can fail with `unknown time zone`. `import _ "time/tzdata"` (or `go build -tags timetzdata`) embeds ~450 KB of zone data.
Docs: https://pkg.go.dev/time/tzdata

Note: Go prints Taipei's abbreviation as `CST` (the tz database's name for UTC+8 there). Abbreviations are ambiguous;
use numeric offsets (`-0700`, RFC 3339) in anything a machine will read.

## Full example

See [main.go](main.go). Run:

```bash
go run .
```

Expected output (verified with go1.26.3 on linux/amd64; also cross-compiled for windows/amd64):

```
== 1. layouts ==
2026-11-08 09:05:07
2026-11-08 09:05:07
2026-11-08T09:05:07Z
Nov 8, 2026 at 9:05am (UTC)
YYYY-MM-DD
2026-08-11
2026-11-08 09:05

== 2. parsing ==
Parse:           2026-11-28 09:00:00 +0000 UTC
ParseInLocation: 2026-11-28 09:00:00 +0800 CST
same instant?    false difference: 8h0m0s
invalid date:    parsing time "2026-02-30": day out of range

== 3. monotonic clock ==
Now has monotonic reading:      true
Round(0) has monotonic reading: false
start == stripped:              false
start.Equal(stripped):          true
elapsed >= 20ms:                true

== 4. comparing ==
2026-01-01 00:00:00 +0000 UTC | 2026-01-01 08:00:00 +0800 UTC+8
==    : false
Equal : true
map lookup with other location: false
map lookup after .UTC():        true
zero value: 0001-01-01 00:00:00 +0000 UTC IsZero: true

== 5. calendar math ==
Jan 31 + 1 month: 2026-03-03
end of next month: 2026-02-28
Date(2026, 13, 0): 2026-12-31
duration: 1h30m30.25s | hours: 1.5084027777777778 | truncated: 1h30m0s

== 6. daylight saving time ==
start:            2026-03-07 12:00:00 -0500 EST
Add(24h):         2026-03-08 13:00:00 -0400 EDT
AddDate(0, 0, 1): 2026-03-08 12:00:00 -0400 EDT
02:30 in the gap: 2026-03-08 01:30:00 -0500 EST
```

The monotonic section prints booleans, so the output is stable from run to run.

## Common pitfalls

1. **Wrong layout, no error.** `"YYYY-MM-DD"` prints literally; `"2006-02-01"` swaps month and day; `"03"` is the
   12-hour clock. `go vet` (the `timeformat` analyzer) only catches the `2006-02-01` case for literal layouts.
2. **`time.Parse` on zone-less strings gives UTC.** Local log times shift by your UTC offset. Use `ParseInLocation`.
3. **`==` and map keys.** Same instant, different location or monotonic flag means `==` is false and the map lookup misses.
   Use `Equal`, and normalize keys.
4. (Bonus) **"One day" as `24*time.Hour`.** Around DST it is not the same as `AddDate(0, 0, 1)`; also `AddDate` overflows
   at month end instead of clamping.

## Exercise

Write `func monthlyBillingDates(start time.Time, n int) []time.Time` for a subscription that started on `start`
(for example 2026-01-31 10:00 in `Asia/Taipei`). Requirements:

- Return the next `n` billing times. If the month has no such day, use the **last day** of that month
  (Jan 31 -> Feb 28 -> Mar 31 -> Apr 30 ...); do not drift (after Feb 28, March is still the 31st).
- Keep 10:00 local wall-clock time, also for a start in `America/New_York` across the March DST change.
- Parse the start from the string `"2026-01-31 10:00"` with the correct location.
- Print each date with an RFC 3339 layout and show that `Equal` holds between your result and a value built with `time.Date`.
- Bonus: explain in a comment why `start.AddDate(0, i, 0)` alone is wrong.
