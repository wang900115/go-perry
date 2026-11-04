# logger

Structured logging with [Uber zap](https://github.com/uber-go/zap). One logger writes JSON to a file and colored console output to stdout at the same time.

## What the code does

`main.go`:

- Builds two encoders with `TimeKey: time`, `LevelKey: level`, `MessageKey: msg`, and ISO8601 timestamps:
  - a **JSON encoder** with lowercase levels
  - a **console encoder** with lowercase *colored* levels
- Opens `log.txt` in append mode (it's created if missing) and builds two cores:
  - file core: JSON encoder, Debug level and above
  - stdout core: console encoder, Debug level and above
- Combines them with `zapcore.NewTee` into one `zap.New(core)` logger. It logs `logger.Info("測試日誌", zap.String("user", "Alice"))` and calls `Sync()` on exit.

Running it appends a line like this to `log.txt`:

```json
{"level":"info","time":"2025-03-13T23:31:29.702+0800","msg":"測試日誌","user":"Alice"}
```

It also prints a tab-separated colored line to the console.

A second, commented-out `main` shows the alternative approach: `zap.Config{...}.Build()` with `OutputPaths: ["stdout", "log.txt"]`. The older entries in the committed `log.txt` came from that version.

## Concepts you can learn

- **Cores, encoders, and sinks**: zap separates *how* to format (encoder), *where* to write (`WriteSyncer`), and *which levels* to keep (`LevelEnabler`). A `Core` combines all three.
- **Tee to multiple outputs**: `NewTee` sends each entry to several cores, so the file can be machine-readable JSON while the terminal stays human-readable.
- **Typed fields instead of string formatting**: `zap.String("user", ...)` avoids reflection and allocations, which is why zap is fast.
- **`Sync()` on exit** flushes buffered output.

## How to run

```bash
cd logger
go run .
```

Running it appends to `logger/log.txt` (a tracked file).

## Notes

- The error from `os.OpenFile` is ignored.
- `EncodeCaller` is set, but `CallerKey` is empty and the logger isn't built with `zap.AddCaller()`, so no caller information is printed.
- The module name is `07`.
