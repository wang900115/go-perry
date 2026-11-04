# directio

Linux `O_DIRECT` file I/O, which bypasses the kernel page cache, compared against normal buffered writes.

## What the code does

- `main_linux.go` (helpers only, no `func main`):
  - `AigenSzie = 4096`: the assumed block size.
  - `MakeAlignedBuffer(size)`: returns a buffer from an anonymous `syscall.Mmap`, which is page-aligned. It rejects sizes that are not multiples of 4096.
  - `FreeAlignedBuffer`: `Munmap`.
  - `OpenDirectFile`: `os.OpenFile` with `syscall.O_DIRECT` added.
- `main_test.go`: two benchmarks that each write a 4096-byte block `b.N` times to a file in `b.TempDir()`:
  - `BenchmarkBufferedWrite`: a normal `os.OpenFile` and a plain `make([]byte)` buffer.
  - `BenchmarkDirectIOWrite`: `OpenDirectFile` and the mmap-aligned buffer.

## Concepts you can learn

- **Page cache vs direct I/O**: buffered writes return once data is copied into kernel memory. `O_DIRECT` sends data straight to the device, so it is slower per call but avoids double caching. Databases use it to manage their own cache.
- **Alignment requirements**: `O_DIRECT` requires the buffer address, file offset and length to be aligned to the device's logical block size. A Go `make([]byte)` slice is not guaranteed to be aligned, which is why the buffer comes from `mmap`.
- **Measuring I/O**: the gap between the two benchmarks shows how much the page cache hides.

## How to run

**Linux only** (`_linux.go` and `syscall.O_DIRECT`; this does not build on Windows/macOS). On Linux:

```bash
cd directio
go test -bench . -benchmem ./...
```

There is nothing to `go run`. On Windows, `GOOS=linux go vet ./...` passes; the benchmarks were not executed.

## Notes

- `b.TempDir()` lives under `$TMPDIR` (usually `/tmp`). Some filesystems, such as some tmpfs setups, reject `O_DIRECT` with `EINVAL`. If the direct benchmark fails to open the file, point `TMPDIR` at a disk-backed filesystem.
- The buffered benchmark never calls `fsync`, so it measures writes into the page cache, not to disk.
- 4096 is hard-coded. Devices with a different logical block size may need another alignment.
- `AigenSzie` is a typo for "AlignSize".
