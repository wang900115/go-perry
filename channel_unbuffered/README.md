# channel_unbuffered

A coffee-shop example of an **unbuffered** channel: every send waits until the consumer takes the order. `channel_buffered/` is the same program with a buffer of 3, for comparison.

## What the code does

`main.go` creates `orders := make(chan string)`. A goroutine sends `Coffee order #1` .. `#5`, prints `Place: ...` after each send succeeds, then closes the channel. `main` ranges over the channel, printing `Preparing: ...`, sleeping 2 s, then printing `Served: ...`.

With no buffer, the producer can hand over the next order only when `main` comes back to receive it, so `Place:` lines appear one at a time, in step with the consumer. Total run time is about 10 s. Whether a `Place:` line prints before or after the matching `Preparing:` line is not deterministic.

## Concepts you can learn

- **Unbuffered channels are a rendezvous**: send and receive complete together, so the producer runs at the consumer's pace.
- **Synchronization through communication**: no mutex or WaitGroup is needed; the channel alone coordinates the two goroutines.
- **`close` + `for range`**: closing the channel ends the consumer loop.

## How to run

This directory has no `go.mod`; it belongs to the root module. From the repo root:

```bash
go run channel_unbuffered/main.go
```
