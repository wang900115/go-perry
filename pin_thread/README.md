# pin_thread

Compare a goroutine locked to an OS thread with `runtime.LockOSThread` against one that is not locked, by printing the OS thread ID each one runs on.

## What the code does

`main.go` starts two goroutines and waits for both with a `sync.WaitGroup`:

- **[1]** calls `runtime.LockOSThread()` (and defers `UnlockOSThread`). It then prints `[1] This goroutine is locked to OS thread: <id>` five times, sleeping 2 s between prints.
- **[2]** does the same without locking and prints `[2] This goroutine is NOT locked to OS thread: <id>`.

The thread ID comes from `windows.GetCurrentThreadId()` (`golang.org/x/sys/windows`). Goroutine [1] should print the same ID every time. Goroutine [2] may print different IDs after waking from `Sleep`, because the scheduler can resume it on any thread. The run takes about 10 s.

## Concepts you can learn

- **M:N scheduling**: goroutines are not tied to OS threads. After blocking or sleeping, a goroutine can resume on a different M (thread).
- **`runtime.LockOSThread`**: pins the calling goroutine to its current thread, and no other goroutine runs on that thread until it is unlocked. This is needed for thread-local OS state, such as GUI/OpenGL contexts, some C libraries, Linux namespaces or `setns`.
- **Pairing lock and unlock**: deferring `UnlockOSThread` releases the thread when the goroutine exits. If a goroutine exits while still locked, the runtime terminates that thread.

## How to run

```bash
cd pin_thread
go run .
```

**Windows only**: the code imports `golang.org/x/sys/windows`. `GOOS=linux go vet` fails with `build constraints exclude all Go files`, and it builds on Windows.

## Notes

- The repo's `.claude/CLAUDE.md` lists `pin_thread` as Linux-only. The code is actually Windows-only.
- An unlocked goroutine may happen to stay on one thread, especially on an idle machine, so identical IDs for [2] do not mean the demo failed.
