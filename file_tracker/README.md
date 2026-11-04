# file_tracker

A small wrapper around `github.com/fsnotify/fsnotify` that watches a directory and reports file create, write, and remove events through a callback.

## What the code does

`main.go` defines a `Watcher` type:

- `New(send)` creates an `fsnotify.Watcher` and starts a `loop` goroutine. The loop forwards only `Write`, `Create`, and `Remove` events to `send`, and logs errors from the `Errors` channel. It exits when either channel closes.
- `Watch(path)`, `Unwatch(path)`, and `Close()` wrap `Add`, `Remove`, and `Close`.

`main` watches the current directory (`.`) and logs every forwarded event (for example, `event: WRITE "./example.txt"`). It then blocks forever with `select {}`. `example.txt` is a scratch file you can edit to trigger events.

## Concepts you can learn

- **OS file-change notification**: fsnotify uses the platform's native API (inotify on Linux), so you don't have to poll file timestamps.
- **Callback + goroutine event loop**: wrapping the two channels (`Events`, `Errors`) in a `select` loop keeps the caller's API to a single function.
- **Filtering events with `event.Has(op)`**: one event can carry several ops at once, so the code checks each flag instead of comparing with `==`.
- **Channel close as a shutdown signal**: the loop exits when `Close()` closes the channels (`ok == false`).

## How to run

```bash
cd file_tracker
go run .
# in another terminal, edit/create/delete files in file_tracker/
```

Stop it with Ctrl+C.

The repository docs list this topic as Linux-only. However, the code has no build constraint and fsnotify supports other platforms. The module compiled on Windows (`go build`), but I only tested that it compiles, not how it behaves at runtime on any OS.

## Notes

- The watch isn't recursive. Changes inside subdirectories of `.` aren't reported.
- `Rename` and `Chmod` events are dropped by the filter.
- Editors often save by writing a temp file and renaming it, so one save can show up as `Create`/`Remove` instead of `Write`.
- `main` blocks forever, so the deferred `Close()` never runs.
