# channel_buffered

A coffee-shop example of a **buffered** channel (capacity 3): the producer can get ahead of a slow consumer until the buffer is full. `channel_unbuffered/` is the same program with an unbuffered channel, for comparison.

## What the code does

`main.go` creates `orders := make(chan string, 3)`. A goroutine sends `Coffee order #1` .. `#5`, prints `Place: ...` after each send succeeds, then closes the channel. `main` ranges over the channel, printing `Preparing: ...`, sleeping 2 s, then printing `Served: ...`.

Because the buffer holds 3 orders, several `Place:` lines print right away, before the first order is served; the producer only blocks once the buffer is full. Total run time is about 10 s (5 orders x 2 s). The exact interleaving of lines is not deterministic.

## Concepts you can learn

- **Buffered channels decouple producer and consumer**: a send blocks only when the buffer is full, so a fast producer is not tied to the consumer's pace (up to the capacity).
- **Backpressure**: once the buffer is full the producer blocks, which limits how far it can get ahead.
- **`close` + `for range`**: closing the channel ends the consumer loop after the buffered values are drained.

## How to run

This directory has no `go.mod`; it belongs to the root module. From the repo root:

```bash
go run channel_buffered/main.go
```

Then run `go run channel_unbuffered/main.go` and compare when the `Place:` lines appear.
