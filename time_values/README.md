# time_values

How `time.Time` values really behave in Go.

- Layouts use the reference time `Mon Jan 2 15:04:05 MST 2006` (`01/02 03:04:05PM '06 -0700`); `YYYY-MM-DD` is printed literally.
- `time.Parse` assumes UTC when the input has no zone; use `time.ParseInLocation` for local wall-clock strings.
- `time.Now()` carries a monotonic clock reading used by `Sub`/`Since`/`Before`/`After`; `Round(0)` strips it.
- Compare instants with `Equal`, not `==`; normalize (`.UTC()`) before using `time.Time` as a map key.
- `AddDate` normalizes overflow (Jan 31 + 1 month = Mar 3), it does not clamp.
- Across DST, `Add(24*time.Hour)` and `AddDate(0, 0, 1)` give different wall-clock results.
- `import _ "time/tzdata"` embeds the zone database so `LoadLocation` works on machines without one (e.g. Windows without Go installed).

```bash
go run .
```

See [LESSON.md](LESSON.md) for the full lesson.
