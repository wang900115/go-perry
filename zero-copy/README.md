# zero-copy

Linux zero-copy data paths: `mmap` shared memory, `sendfile`, `splice`, `vmsplice` and `MSG_ZEROCOPY`. Each one avoids copying data between user space and the kernel.

## What the code does

One module (`go.mod`, depends on `golang.org/x/sys`). Each subdirectory holds standalone `package main` programs:

- `mmap/shared_memory_writer_linux.go`: creates `/tmp/my_shared_mem`, sizes it to a `Packet{ID, Score uint32; Ready bool}` and `mmap`s it with `MAP_SHARED`. It writes `ID=1, Score=100, Ready=true` through a pointer cast, sleeps 10 s, then removes the file.
- `mmap/shared_memory_reader_linux.go`: maps the same file read-only, polls `Ready` every 100 ms and prints `ID: 1, Score: 100`.
- `sendfile/file_writer_linux.go`: opens `data.txt`, listens on `:8080`, accepts one connection and calls `syscall.Sendfile` once (up to 1 MiB) from the file to the socket. Prints `file sent`.
- `splice/pipe_k2k_linux.go`: opens `data.txt`, listens on `:8080` and, for one client, loops file → pipe → socket with `unix.Splice` in 64 KiB chunks. It prints `spliced N bytes` per chunk and `transfer completed` at the end.
- `vmsplice/pipe_u2k_linux.go`: `unix.Vmsplice` maps a user buffer (`"hello from vmsplice\n"`) into a pipe, reads it back from the other end, and prints `vmsplice: 20 bytes` and the string.
- `msg_zerocopy/network_linux.go`: connects to `127.0.0.1:9090`, enables `SO_ZEROCOPY` and sends 64 KiB of `'A'` with `SendmsgN(..., MSG_ZEROCOPY)`. Prints `sent 65536 bytes with MSG_ZEROCOPY`.
- `reader_linux`, `writer_linux`: committed Linux ELF binaries, apparently built from the mmap reader and writer.

## Concepts you can learn

- **mmap shared memory**: two processes map the same file pages, so a write by one is visible to the other with no `read`/`write` syscalls.
- **`sendfile`**: copies file → socket entirely inside the kernel. This is the classic static-file-server fast path.
- **`splice`**: moves data between a file descriptor and a pipe without a user-space copy. Two splices (file→pipe, pipe→socket) connect any two fds.
- **`vmsplice`**: hands user-space pages to a pipe (user → kernel), the counterpart of `splice`.
- **`MSG_ZEROCOPY`**: the kernel transmits directly from the user buffer, after the socket opts in with `SO_ZEROCOPY`.
- **`SyscallConn().Control`**: how to get the raw fd of a `net.Conn` safely for low-level syscalls.

## How to run

**Linux only.** Run each file individually, because `mmap/` has two `main` functions in one directory:

```bash
cd zero-copy

# mmap: start the writer, then the reader within 10 s (second terminal)
go run mmap/shared_memory_writer_linux.go
go run mmap/shared_memory_reader_linux.go

# sendfile / splice: need ./data.txt, then connect a client to :8080
echo hello > data.txt
go run ./sendfile        # or: go run ./splice
nc localhost 8080

go run ./vmsplice

# msg_zerocopy: needs a TCP server on 127.0.0.1:9090 first
nc -l 9090 > /dev/null &
go run ./msg_zerocopy
```

On Windows, `GOOS=linux go vet` passes for each program; none were executed.

## Notes

- `go build ./...` fails in `mmap/`: both files declare `main` and `Packet` in the same package.
- The mmap reader polls `Ready` with no atomic operation or memory barrier. That is fine for a demo but not a correct synchronization protocol. The writer sets `Ready` last in program order, but the code does not guarantee that the reader sees the writes in that order.
- `sendfile` makes a single call capped at 1 MiB. Larger files, or a short write, are not handled.
- Over loopback (`127.0.0.1`) the kernel performs a deferred copy for MSG_ZEROCOPY packets (per the kernel's `msg_zerocopy` documentation; not verified here), so this setup does not actually avoid the copy; use a remote peer to observe the effect. The program also prints `len(data)` rather than the byte count returned by `SendmsgN`.
- `msg_zerocopy` never reads the socket error queue for completion notifications, so it cannot know when the buffer is safe to reuse. It sleeps 1 s instead.
- `sendfile` panics inside the `Control` callback on error instead of returning the error.
