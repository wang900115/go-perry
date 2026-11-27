# encoding/binary: Byte Order, Fixed-Size Structs, and Varints

> Go version used: `go1.26.3`. `Encode`/`Decode`/`Append` need Go 1.23+, `NativeEndian` needs Go 1.21+, `AppendUint32` and friends need Go 1.19+. Everything else has existed since Go 1.0.
> Docs: https://pkg.go.dev/encoding/binary

## Why learn this

Any time bytes leave your process in a format you control, not JSON:

- Binary protocols and file formats: TCP message headers, WAL / SSTable records (see `lsm/` in this repo), PNG/ZIP/ELF headers
- Length-prefixed framing over a TCP stream (a stream has no message boundaries)
- Compact IDs and keys: big-endian `uint64` keys sort correctly as raw bytes in KV stores (Badger, Pebble, BoltDB)
- Hashing / checksums over a stable byte layout (Bloom filters, see `boomfilter/`)
- Understanding Protocol Buffers: protobuf uses exactly the varint and zigzag encodings in this package

## Core concepts

### 1. Byte order

A `uint32` is 4 bytes, and something has to decide their order on the wire.

```go
binary.BigEndian.AppendUint32(nil, 0x01020304)    // 01 02 03 04  (most significant byte first, "network order")
binary.LittleEndian.AppendUint32(nil, 0x01020304) // 04 03 02 01  (x86/amd64 and arm64 memory order)
```

`ByteOrder` has `Uint16/32/64` and `PutUint16/32/64`; `AppendByteOrder` adds `AppendUint16/32/64`.
They are plain shifts and are not reflection based. The compiler recognises the pattern and turns it into a single load/store (plus a byte swap when needed) on common architectures.
`binary.NativeEndian` (Go 1.21) is the host order: use it only for data that never leaves the machine.

### 2. Fixed-size values and reflection

`binary.Write`, `binary.Read`, `binary.Size`, and the Go 1.23 `Encode`, `Decode`, `Append` accept "fixed-size values":
fixed-size numeric types (`int8`...`uint64`, `float32/64`, `complex64/128`), `bool`, and arrays or structs made only of those.

- `int` and `uint` are **not** fixed-size (32 or 64 bits depending on the platform), so they are rejected with an error, and `Size` returns `-1`.
- Strings, slices inside structs, maps, and pointers are rejected too. A top-level slice of fixed-size values such as `[]uint32` is accepted.
- Fields are written in declaration order with **no padding**: `binary.Size(Header{})` is 18 while `unsafe.Sizeof` is 24 on amd64 (the compiler pads `Length` to a 4-byte boundary and `Seq` to an 8-byte boundary).
- `_` (blank) fields are written as zeros and skipped on read, which is handy for reserved bytes.
- For a pointer to a basic type or a slice of basic types (`*uint32`, `[]uint64`), `Read`/`Write` take a fast path without reflection; structs go through `reflect`.

Source: https://cs.opensource.google/go/go/+/refs/tags/go1.26.3:src/encoding/binary/binary.go

### 3. Encode / Decode / Append (Go 1.23)

`Write`/`Read` need an `io.Writer`/`io.Reader`. When you already have a `[]byte`:

- `Encode(buf, order, v) (int, error)` writes into `buf` and returns an error if it is too small
- `Decode(buf, order, &v) (int, error)` reads from `buf` and returns how many bytes it consumed
- `Append(buf, order, v) ([]byte, error)` appends, growing the slice as needed

Release notes: https://go.dev/doc/go1.23 (section "encoding/binary")

### 4. Varints

`PutUvarint`/`AppendUvarint` store 7 bits per byte, low bits first; the top bit of each byte means "more bytes follow".
`300 = 0b1_0010_1100` becomes `ac 02` (`0x2c | 0x80`, then `0x02`). A `uint64` takes 1 to 10 bytes (`MaxVarintLen64 = 10`).

Signed `PutVarint`/`AppendVarint` first apply **zigzag**: `ux := uint64(x) << 1; if x < 0 { ux = ^ux }`.
That maps `0, -1, 1, -2, 2 ...` to `0, 1, 2, 3, 4 ...`, so `-1` takes 1 byte. Without zigzag, `uint64(int64(-1))` is all ones and takes 10 bytes.

Decoding returns `(value, n)`:

| `n` | meaning |
| --- | --- |
| `n > 0` | ok, `n` bytes consumed |
| `n == 0` | buffer too small (truncated input) |
| `n < 0` | value overflows 64 bits; `-n` bytes were read |

`ReadUvarint(io.ByteReader)` returns `io.EOF` only if no byte was read, and `io.ErrUnexpectedEOF` if the input ends in the middle of a varint.

## Complete example

`main.go` in this directory (run with `go run .`):

1. The same `uint32` in big- and little-endian, and what you get when you read with the wrong order
2. A `Header` struct written with `binary.Write` and read back with `binary.Read`; `int` field rejected
3. Go 1.23 `Encode` into a too-small buffer, `Append`, `Decode`
4. Uvarint sizes for 1, 127, 128, 300, 2^32
5. Zigzag `Varint` vs casting a negative number to `uint64`
6. `Uvarint` return values for ok / truncated / empty / overflow input
7. Length-prefixed framing with `AppendUvarint`, `ReadUvarint`, `io.ReadFull`, and a size cap

Expected output (amd64; `unsafe.Sizeof` would print 20 on 32-bit x86 because `uint64` is only 4-byte aligned there):

```text

== 1. byte order ==
BigEndian    01 02 03 04
LittleEndian 04 03 02 01
read LE bytes as BE -> 0x04030201 (wrong!)

== 2. fixed-size struct ==
unsafe.Sizeof=24 binary.Size=18
wire bytes   50 52 52 59 01 80 00 00 00 05 00 00 00 00 00 00 00 2a
decoded      {Magic:[80 82 82 89] Version:1 Flags:128 Length:5 Seq:42} equal=true
int field    -> binary.Write: some values are not fixed-sized in type main.BadHeader
Size(int)    -> -1

== 3. Encode/Decode/Append (Go 1.23+) ==
Encode into 10 bytes -> buffer too small
Append len=19 first=aa 50 52 52 59
Decode n=18 err=<nil> seq=42

== 4. uvarint sizes ==
1          -> 1 byte(s) 01
127        -> 1 byte(s) 7f
128        -> 2 byte(s) 80 01
300        -> 2 byte(s) ac 02
4294967296 -> 5 byte(s) 80 80 80 80 10

== 5. zigzag (Varint) vs casting to uint64 ==
   0 -> Varint 00     (1 B) | Uvarint(uint64)  1 B
  -1 -> Varint 01     (1 B) | Uvarint(uint64) 10 B
   1 -> Varint 02     (1 B) | Uvarint(uint64)  1 B
  -2 -> Varint 03     (1 B) | Uvarint(uint64) 10 B
   2 -> Varint 04     (1 B) | Uvarint(uint64)  1 B
 -64 -> Varint 7f     (1 B) | Uvarint(uint64) 10 B
  64 -> Varint 80 01  (2 B) | Uvarint(uint64)  1 B

== 6. Uvarint error signalling ==
in=ac 02                             -> x=300 n=2 (ok)
in=ac                                -> x=0 n=0 (buffer too small)
in=                                  -> x=0 n=0 (buffer too small)
in=ff ff ff ff ff ff ff ff ff ff ff  -> x=0 n=-11 (overflow)

== 7. length-prefixed frames ==
stream 02 68 69 00 06 67 6f 70 68 65 72
frame "hi"
frame ""
frame "gopher"
done (clean EOF)
truncated frame -> read payload: unexpected EOF
```

## Common pitfalls

1. **Using `int` in a wire struct.** `binary.Write` returns `some values are not fixed-sized` and `binary.Size` returns `-1`. Use `int32`/`int64`/`uint32` explicitly. If you ignore the error you silently write nothing.
2. **Mismatched byte order, or assuming `unsafe.Sizeof` is the wire size.** Reading little-endian bytes as big-endian gives `0x04030201` instead of `0x01020304`, with no error. Buffers sized with `unsafe.Sizeof` include padding; use `binary.Size`.
3. **Ignoring `n` from `Uvarint`/`Varint`.** On truncated input you get `x = 0, n = 0`, which looks like a valid zero if you only look at `x`. Also, a length prefix read from the network must be capped before `make([]byte, size)`, or a hostile peer can make you allocate gigabytes.
4. **Using `binary.Write` on a hot path.** Struct encoding goes through `reflect`. For a few fields, `AppendUint32`/`AppendUint64` by hand is much faster and allocation-free.

## Exercise

Write a tiny record format for a key-value log, like one entry of a WAL:

```
[crc32 uint32 BE][keyLen uvarint][valLen uvarint][key bytes][val bytes]
```

- `func AppendRecord(dst []byte, key, val []byte) []byte` appends one record; the CRC (`hash/crc32.ChecksumIEEE`) covers everything after the CRC field.
- `func ReadRecord(r *bufio.Reader) (key, val []byte, err error)` reads one record and returns an error on a CRC mismatch, a truncated record, or a length over 1 MiB.
- In `main`, write 3 records into a `bytes.Buffer`, flip one byte in the middle of the second record, and show that the first record decodes, the second reports a CRC error, and a clean `io.EOF` is returned at the end of the stream.
- Bonus: benchmark `binary.Write` of a struct against hand-written `AppendUint32`/`AppendUint64` with `testing.B` and `-benchmem`.
