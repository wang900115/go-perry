# goroutine_confinement

A ticket-selling simulation that avoids locks by confining the shared counter to a single goroutine; all other goroutines talk to it through a channel. Compare with `goroutine_mutex/`, which solves the same problem with a `sync.Mutex`.

## What the code does

`main.go` starts with 200 tickets and 2000 users:

- `managerTickets` is the only goroutine that reads or writes the ticket count. It loops forever on a `select`: for each user ID received on `ticketChannel` it either decrements and prints `User N purchased the ticket, remaining M tickets`, or prints `User N found no ticket`.
- `main` starts 2000 `buyTicket` goroutines; each just sends its user ID on the unbuffered `ticketChannel`.
- After `wg.Wait()`, `main` sends on `doneChannel`, and the manager prints `Ticket remaing 0`.

Expected result: exactly 200 purchases and 1800 "found no ticket" lines. Which user IDs get tickets is random.

## Concepts you can learn

- **Confinement**: if only one goroutine ever touches a piece of state, there is no data race and no lock is needed. "Share memory by communicating."
- **Serializing requests through a channel**: the unbuffered channel turns 2000 concurrent buyers into a sequential stream of requests for the manager.
- **`sync.WaitGroup`** to know when all buyers have finished sending.

## How to run

Own `go.mod` (`module confinement`):

```bash
cd goroutine_confinement
go run .
go run -race .     # race detector; requires cgo on Windows
```

## Notes

- `managerTickets` never returns; on `doneChannel` it prints and keeps looping. The program ends only because `main` returns. The final `Ticket remaing ...` line (typo in source, no trailing newline) may not print, since `main` can exit right after the send.
- The count is passed in as a pointer to `main`'s variable; confinement holds only because `main` does not read it after starting the manager.
