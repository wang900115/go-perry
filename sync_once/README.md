# sync_once

`sync.Once`, `sync.OnceFunc`, `sync.OnceValue`, `sync.OnceValues`: run an initialization exactly once, safely, even when many goroutines race to trigger it.

## Run

```
go run .
go run -race .
```

## What the code does

- **Part 1 – `sync.Once`**: 100 goroutines call `GetConfig()` at the same time; the loader runs only once (`loadCount == 1`).
- **Part 2 – `OnceValue` / `OnceValues` (Go 1.21+)**: wrap "compute once, return the value" (and `(value, error)`) into a plain function variable. No package-level `once` + result pair needed.
- **Part 3 – panic behavior**: if `f` panics, `Once.Do` treats it as *done* — the 2nd call silently returns and the value stays uninitialized. `OnceFunc`/`OnceValue(s)` instead re-panic with the same value on every call.

Expected output (observed on Linux with go1.26.3, also under `-race`; `GOOS=windows go build` also succeeds):

```
== 1. sync.Once ==
100 個 goroutine 呼叫 GetConfig，實際載入次數 = 1
config = {DSN:file:app.db MaxConn:10}
== 2. OnceValue / OnceValues ==
  (計算平方表中…只會印一次)
table[3] = 9
table[9] = 81
  (讀取 secret 中…只會印一次)
第 1 次 loadSecret: "", err=secret file not found
第 2 次 loadSecret: "", err=secret file not found
== 3. panic 時的行為 ==
  Once.Do 第 1 次 -> panic: init failed
  Once.Do 第 2 次 -> 正常返回（f 沒有再被執行）
  OnceFunc 第 1 次 -> panic: init failed
  OnceFunc 第 2 次 -> panic: init failed
```

## Concepts you can learn

- **How `Once` works**: an `atomic.Bool done` fast path (inlined at the call site) plus a `Mutex` slow path. `done` is set only *after* `f` returns, so every caller of `Do` is guaranteed that `f` has finished — a plain `CompareAndSwap` would not give that guarantee.
- **Errors are cached too**: `OnceValues` memoizes `(v, err)`; a failed init is not retried. If you need retry-on-error, use a `Mutex` + flag yourself.
- **Deadlock**: calling `once.Do` again from inside `f` on the same `Once` deadlocks (non-reentrant mutex).
- **Do not copy a `Once`**: it embeds `noCopy`, `go vet` (copylocks) reports copies.

References: <https://pkg.go.dev/sync#Once>, <https://pkg.go.dev/sync#OnceValue>, <https://go.dev/doc/go1.21#sync>
