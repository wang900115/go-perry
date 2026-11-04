# capProto

Building, serializing, and reading a Cap'n Proto message in Go with `capnproto.org/go/capnp/v3`. Cap'n Proto is a schema-based binary format whose in-memory layout is also its wire format.

## What the code does

- `schema/user.capnp` - schema: `struct User { id @0 :UInt64; name @1 :Text; age @2 :UInt8; active @3 :Bool; }`, with Go package name `generated`.
- `schema/user.capnp.go` - **generated** by `capnpc-go` (the Cap'n Proto Go plugin) from `schema/user.capnp`. Do not edit it by hand. It is imported as `capproto/schema` under the alias `generated`.
- `main.go`:
  1. Creates a single-segment message with `capnp.NewMessage(capnp.SingleSegment(nil))` and a root `User` with `NewRootUser`.
  2. Sets the four fields, calls `msg.Marshal()`, and prints the size and hex bytes.
  3. Calls `capnp.Unmarshal(buf)`, then `ReadRootUser`, and prints the fields.

Actual output:

```
Binary size: 48 bytes
Binary data: 00000000050000000000000002000100e903000000000000190100000000000001000000320000005065727279000000

User:
ID: 1001
Name: Perry
Age: 25
Active: true
```

## Concepts you can learn

- **No encode/decode step**: setters write directly into the segment, and getters read from it. `Marshal` mostly adds the segment table header (the first 8 bytes) in front of the segment.
- **Fixed data section + pointer section**: `id`, `age`, and `active` are packed into one data section (the 2 words seen in the struct pointer `02000100`). `name` is a pointer to a separate text blob. This is why the message is word-aligned (multiples of 8 bytes) and padded.
- **Segments and arenas**: `SingleSegment` puts everything in one contiguous buffer, which is the simplest and fastest case.
- **Errors on pointer fields**: `Name()` returns `(string, error)` because following a pointer can fail on malformed input. The code discards that error.

## How to run

Own `go.mod` (module `capproto`):

```bash
cd capProto
go run .
```

To regenerate `user.capnp.go` you need the `capnp` compiler and the `capnpc-go` plugin, plus the `go.capnp` schema from go-capnp on the include path (the schema does `import "/go.capnp"`). The exact include flag depends on where go-capnp is installed. Not verified here.

## Notes

- The schema's annotation `$Go.import("go-perry/capProto/generated")` does not match the real import path, which is `capproto/schema`. It only affects code generation metadata, and the example builds and runs as is.
- The error from `user.SetName("Perry")` and from `user2.Name()` is ignored.
