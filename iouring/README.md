# iouring

A TCP echo server driven by Linux `io_uring` through liburing. The ring logic is written in C and the event loop is in Go, connected with cgo.

## What the code does

- `uring.h` / `uring.c`:
  - Wrap liburing in a few functions: `server_socket`, `uring_init`/`uring_cleanup`, `submit_accept`/`submit_recv`/`submit_send`, and `wait_cqe`.
  - Keep a per-connection struct with a 4096-byte buffer plus `len`/`offset`.
  - Allocate a small `request {op, conn}` for each submission and attach it as the SQE user data. `wait_cqe` turns the completion into a `completion_t {op, conn, result}` and frees the request.
  - `submit_send` writes debug lines to stderr.
- `main.go`:
  1. Creates a listening socket on **127.0.0.1:8080** (loopback only, `SO_REUSEADDR`) and a ring with 256 entries, then submits the first accept.
  2. Loops on `wait_cqe`:
     - **ACCEPT**: creates a connection, submits a recv, and re-arms the accept.
     - **RECV**: if the result is 0 (closed) or negative (error), frees the connection. Otherwise it prints `Received N bytes: ...` and submits a send.
     - **SEND**: advances the offset. If bytes remain, it submits another send. Otherwise it submits the next recv.

## Concepts you can learn

- **Submission and completion queues**: operations are queued as SQEs and their results arrive later as CQEs. The model is completion-based, unlike readiness-based epoll.
- **Keeping per-operation context in user data**: `io_uring_sqe_set_data` and `io_uring_cqe_get_data` let you route completions without a lookup table.
- **Handling partial sends**: `connection_consume` and `connection_remaining` resubmit the unsent tail.
- **Re-arming accepts**: each accept SQE completes once, so the loop must submit a new one after every connection.
- **cgo across languages**: calling C functions, using C enums and structs from Go, linking with `#cgo LDFLAGS: -luring`, and copying C memory with `C.GoBytes`.

## How to run

**Linux only.** You need a kernel with io_uring support, cgo enabled (a C compiler), and the liburing development package (for example, `liburing-dev` on Debian/Ubuntu).

```bash
cd iouring
go run .
# in another terminal:
nc 127.0.0.1 8080
```

I couldn't build or run this topic on the Windows dev machine.

## Notes

- The SIGINT/SIGTERM check only runs between completions. `signal.Notify` captures Ctrl+C, so the process won't exit until the next CQE arrives (for example, a new connection or data on an existing one).
- If `server_socket` fails, the code panics. It also panics on a ring init failure or a failed first accept.
- `submit_accept` and `submit_recv` don't free their `request` when `io_uring_submit` fails (only `submit_send` does). If `submit_recv` fails right after an accept, `main.go` frees the connection without logging anything.
- `submit_send` prints debug lines to stderr on every send.
- The listener binds to loopback only, so it can't be reached from other hosts.
