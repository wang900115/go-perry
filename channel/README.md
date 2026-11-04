# channel

Eight small programs that introduce goroutines, `sync.WaitGroup`, channels, `select`, and the `context` package, one idea per file.

## What the code does

- `01.go`: goroutine with an anonymous func. It sends a string on an unbuffered channel; `main` receives it and prints `Hello from goroutine!`.
- `02.go`: goroutine with sync (`Add`, `Done`, `Wait`). `main` waits on a `sync.WaitGroup` until the goroutine prints `goroutine 執行中`, then prints `主執行緒結束`.
- `03.go`: goroutine with channel. Shows that a send on an unbuffered channel blocks until it is received. It prints `準備發送` and `準備接收` (order not fixed), then `收到: 42`. `發送完畢` may or may not appear, because `main` can return before the goroutine prints it.
- `04.go`: channel writer and reader, `chan<-` / `<-chan`. `sendOnly(ch chan<- int)` sends 42 and `receiveOnly(ch <-chan int)` prints it.
- `05.go`: channel operation (`close`). A goroutine writes 1, 2, 3 into a buffered channel (capacity 3) and closes it; `for range` prints `1`, `2`, `3` and exits when the channel is closed.
- `06.go`: producer and consumer. The producer sends 1..5 with a 1 s pause after each send, then closes the channel; the consumer ranges over it. `生產: n` and `消費: n` lines interleave, then `完成`. Takes about 5 s.
- `07.go`: `select` over channels. Two goroutines each send a message after 1 s; a 2-iteration loop `select`s on `ch1`, `ch2`, and `time.After(2 * time.Second)`. Both messages print (order not fixed), then `完成`. The timeout branch (`超時了`) does not normally fire.
- `08.go`: the `context` package. A worker prints `Worker 1 工作中` every 500 ms until `cancel()` is called after 2 s, then prints `Worker 1 退出: context canceled`; `main` prints `主程序結束` 1 s later.

> [!TIP]
> First In First Out: a channel delivers values in the order they were sent (see `05.go`).

## Concepts you can learn

- **Unbuffered channels synchronize**: a send waits for a receiver (01, 03, 04), so a channel is both a data pipe and a sync point.
- **`sync.WaitGroup`**: waits for goroutines that return no value (02, 06). Without it, `main` can exit before they run.
- **Directional channel types** (`chan<-`, `<-chan`): the compiler rejects misuse, which documents who owns which end (04, 06).
- **`close` + `for range`**: the sender closes to signal "no more values", which ends the receiver's loop (05, 06).
- **Buffered channel capacity**: sends do not block until the buffer is full (05).
- **`select` with `time.After`**: wait on several channels at once and add a timeout (07).
- **Cancellation with `context.WithCancel`**: a goroutine checks `ctx.Done()` in a `select` so it can stop cooperatively (08).

## How to run

`channel/` has no `go.mod`; it is part of the root module, and each file declares its own `func main`. Run one file at a time from the repo root:

```bash
go run channel/01.go
go run channel/02.go
# ...
go run channel/08.go
```

`go build ./channel` fails by design (several `main` functions in one package).

## Notes

- Program output is in Chinese.
- In `06.go`, `defer wg.Done()` is the last statement of each function. It works, but the deferred call runs only at return; in `producer` that is after `close(ch)`. Putting the `defer` first is the usual style.
- In `08.go`, the `select` has a `default` branch, so the worker only notices cancellation between its 500 ms sleeps.
