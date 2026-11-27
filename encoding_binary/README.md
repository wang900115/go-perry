# encoding_binary

Turning numbers and fixed-size structs into bytes (and back) with the standard library `encoding/binary` package.

- `binary.BigEndian` / `binary.LittleEndian` decide byte order; both sides of a protocol must agree. Network protocols are usually big-endian.
- `PutUint32` / `AppendUint32` / `Uint32` are fast, reflection-free helpers. `PutUint32` panics if the buffer is too short.
- `binary.Write` / `Read` / `Size` (and Go 1.23 `Encode` / `Decode` / `Append`) use reflection and accept only fixed-size types: `int`, `uint`, strings, slices inside structs, and maps are rejected.
- `binary.Size` is the wire size with no padding; `unsafe.Sizeof` is the in-memory size with padding.
- Varints: `Uvarint` uses 7 bits per byte; `Varint` adds zigzag so small negative numbers stay small. Casting a negative `int64` to `uint64` costs 10 bytes.
- `Uvarint` returns `n == 0` for a short buffer and `n < 0` for overflow; always check `n`.
- Length-prefixed framing: uvarint length + payload, read with `ReadUvarint` + `io.ReadFull`, and cap the length.

```bash
go run .
```

See [LESSON.md](LESSON.md) for the full lesson.
