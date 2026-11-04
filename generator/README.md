# generator

The generator pattern with channels: a goroutine produces an endless stream of prime numbers, and the consumer stops it with a `done` channel.

## What the code does

`main.go`:

- `isPrime` / `nextPrime` find the next prime after `n` by trial division up to `sqrt(n)` (`nextPrime(1)` returns 2).
- `primeGenerator(done, operation)` starts a goroutine that repeatedly applies `operation` (here `nextPrime`) and sends each result on an unbuffered channel. Each send is in a `select` with `<-done`, so the generator exits (and closes its channel) when signalled.
- `main` ranges over the generator, prints `Prime # 1 :  2` through `Prime # 10 :  29`, breaks after 10, then sends `done <- true`.

## Concepts you can learn

- **Generator as a goroutine + channel**: values are produced lazily, one per receive, so an infinite sequence is safe to consume.
- **Stopping a producer with a done channel**: breaking out of `range` alone would leave the goroutine blocked forever on its send (a goroutine leak); the `select` on `done` lets it exit.
- **Passing behavior as a function** (`operation func(int) int`): the generator does not know it produces primes.
- **Returning a receive-only channel** (`<-chan int`) so callers cannot send into it.

## How to run

This directory has no `go.mod`; it belongs to the root module. From the repo root:

```bash
go run generator/main.go
```

## Notes

- `Received done signal, stopping generator.` may not print: `main` returns right after its `done <- true` send completes, which can end the process before the goroutine reaches `fmt.Println`.
- `done` is sent a value rather than closed; that works for one generator, but `close(done)` is the usual way to broadcast to many.
